
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:46</date>
//</624342668939366400>

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

package api

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/log"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

const maxParallelFiles = 5

type FileSystem struct {
	api *API
}

func NewFileSystem(api *API) *FileSystem {
	return &FileSystem{api}
}

//
//
//
//
//
//
func (fs *FileSystem) Upload(lpath, index string, toEncrypt bool) (string, error) {
	var list []*manifestTrieEntry
	localpath, err := filepath.Abs(filepath.Clean(lpath))
	if err != nil {
		return "", err
	}

	f, err := os.Open(localpath)
	if err != nil {
		return "", err
	}
	stat, err := f.Stat()
	if err != nil {
		return "", err
	}

	var start int
	if stat.IsDir() {
		start = len(localpath)
		log.Debug(fmt.Sprintf("uploading '%s'", localpath))
		err = filepath.Walk(localpath, func(path string, info os.FileInfo, err error) error {
			if (err == nil) && !info.IsDir() {
				if len(path) <= start {
					return fmt.Errorf("Path is too short")
				}
				if path[:start] != localpath {
					return fmt.Errorf("Path prefix of '%s' does not match localpath '%s'", path, localpath)
				}
				entry := newManifestTrieEntry(&ManifestEntry{Path: filepath.ToSlash(path)}, nil)
				list = append(list, entry)
			}
			return err
		})
		if err != nil {
			return "", err
		}
	} else {
		dir := filepath.Dir(localpath)
		start = len(dir)
		if len(localpath) <= start {
			return "", fmt.Errorf("Path is too short")
		}
		if localpath[:start] != dir {
			return "", fmt.Errorf("Path prefix of '%s' does not match dir '%s'", localpath, dir)
		}
		entry := newManifestTrieEntry(&ManifestEntry{Path: filepath.ToSlash(localpath)}, nil)
		list = append(list, entry)
	}

	cnt := len(list)
	errors := make([]error, cnt)
	done := make(chan bool, maxParallelFiles)
	dcnt := 0
	awg := &sync.WaitGroup{}

	for i, entry := range list {
		if i >= dcnt+maxParallelFiles {
			<-done
			dcnt++
		}
		awg.Add(1)
		go func(i int, entry *manifestTrieEntry, done chan bool) {
			f, err := os.Open(entry.Path)
			if err == nil {
				stat, _ := f.Stat()
				var hash storage.Address
				var wait func(context.Context) error
				ctx := context.TODO()
				hash, wait, err = fs.api.fileStore.Store(ctx, f, stat.Size(), toEncrypt)
				if hash != nil {
					list[i].Hash = hash.Hex()
				}
				err = wait(ctx)
				awg.Done()
				if err == nil {
					first512 := make([]byte, 512)
					fread, _ := f.ReadAt(first512, 0)
					if fread > 0 {
						mimeType := http.DetectContentType(first512[:fread])
						if filepath.Ext(entry.Path) == ".css" {
							mimeType = "text/css"
						}
						list[i].ContentType = mimeType
					}
				}
				f.Close()
			}
			errors[i] = err
			done <- true
		}(i, entry, done)
	}
	for dcnt < cnt {
		<-done
		dcnt++
	}

	trie := &manifestTrie{
		fileStore: fs.api.fileStore,
	}
	quitC := make(chan bool)
	for i, entry := range list {
		if errors[i] != nil {
			return "", errors[i]
		}
		entry.Path = RegularSlashes(entry.Path[start:])
		if entry.Path == index {
			ientry := newManifestTrieEntry(&ManifestEntry{
				ContentType: entry.ContentType,
			}, nil)
			ientry.Hash = entry.Hash
			trie.addEntry(ientry, quitC)
		}
		trie.addEntry(entry, quitC)
	}

	err2 := trie.recalcAndStore()
	var hs string
	if err2 == nil {
		hs = trie.ref.Hex()
	}
	awg.Wait()
	return hs, err2
}

//
//
//
//
func (fs *FileSystem) Download(bzzpath, localpath string) error {
	lpath, err := filepath.Abs(filepath.Clean(localpath))
	if err != nil {
		return err
	}
	err = os.MkdirAll(lpath, os.ModePerm)
	if err != nil {
		return err
	}

//
	uri, err := Parse(path.Join("bzz:/", bzzpath))
	if err != nil {
		return err
	}
	addr, err := fs.api.Resolve(context.TODO(), uri.Addr)
	if err != nil {
		return err
	}
	path := uri.Path

	if len(path) > 0 {
		path += "/"
	}

	quitC := make(chan bool)
	trie, err := loadManifest(context.TODO(), fs.api.fileStore, addr, quitC, NOOPDecrypt)
	if err != nil {
		log.Warn(fmt.Sprintf("fs.Download: loadManifestTrie error: %v", err))
		return err
	}

	type downloadListEntry struct {
		addr storage.Address
		path string
	}

	var list []*downloadListEntry
	var mde error

	prevPath := lpath
	err = trie.listWithPrefix(path, quitC, func(entry *manifestTrieEntry, suffix string) {
		log.Trace(fmt.Sprintf("fs.Download: %#v", entry))

		addr = common.Hex2Bytes(entry.Hash)
		path := lpath + "/" + suffix
		dir := filepath.Dir(path)
		if dir != prevPath {
			mde = os.MkdirAll(dir, os.ModePerm)
			prevPath = dir
		}
		if (mde == nil) && (path != dir+"/") {
			list = append(list, &downloadListEntry{addr: addr, path: path})
		}
	})
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	errC := make(chan error)
	done := make(chan bool, maxParallelFiles)
	for i, entry := range list {
		select {
		case done <- true:
			wg.Add(1)
		case <-quitC:
			return fmt.Errorf("aborted")
		}
		go func(i int, entry *downloadListEntry) {
			defer wg.Done()
			err := retrieveToFile(quitC, fs.api.fileStore, entry.addr, entry.path)
			if err != nil {
				select {
				case errC <- err:
				case <-quitC:
				}
				return
			}
			<-done
		}(i, entry)
	}
	go func() {
		wg.Wait()
		close(errC)
	}()
	select {
	case err = <-errC:
		return err
	case <-quitC:
		return fmt.Errorf("aborted")
	}
}

func retrieveToFile(quitC chan bool, fileStore *storage.FileStore, addr storage.Address, path string) error {
f, err := os.Create(path) //
	if err != nil {
		return err
	}
	reader, _ := fileStore.Retrieve(context.TODO(), addr)
	writer := bufio.NewWriter(f)
	size, err := reader.Size(context.TODO(), quitC)
	if err != nil {
		return err
	}
	if _, err = io.CopyN(writer, reader, size); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return f.Close()
}

