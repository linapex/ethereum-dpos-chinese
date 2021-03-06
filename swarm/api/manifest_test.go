
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342670021496832>

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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

func manifest(paths ...string) (manifestReader storage.LazySectionReader) {
	var entries []string
	for _, path := range paths {
		entry := fmt.Sprintf(`{"path":"%s"}`, path)
		entries = append(entries, entry)
	}
	manifest := fmt.Sprintf(`{"entries":[%s]}`, strings.Join(entries, ","))
	return &storage.LazyTestSectionReader{
		SectionReader: io.NewSectionReader(strings.NewReader(manifest), 0, int64(len(manifest))),
	}
}

func testGetEntry(t *testing.T, path, match string, multiple bool, paths ...string) *manifestTrie {
	quitC := make(chan bool)
	fileStore := storage.NewFileStore(nil, storage.NewFileStoreParams())
	ref := make([]byte, fileStore.HashSize())
	trie, err := readManifest(manifest(paths...), ref, fileStore, false, quitC, NOOPDecrypt)
	if err != nil {
		t.Errorf("unexpected error making manifest: %v", err)
	}
	checkEntry(t, path, match, multiple, trie)
	return trie
}

func checkEntry(t *testing.T, path, match string, multiple bool, trie *manifestTrie) {
	entry, fullpath := trie.getEntry(path)
	if match == "-" && entry != nil {
		t.Errorf("expected no match for '%s', got '%s'", path, fullpath)
	} else if entry == nil {
		if match != "-" {
			t.Errorf("expected entry '%s' to match '%s', got no match", match, path)
		}
	} else if fullpath != match {
		t.Errorf("incorrect entry retrieved for '%s'. expected path '%v', got '%s'", path, match, fullpath)
	}

	if multiple && entry.Status != http.StatusMultipleChoices {
		t.Errorf("Expected %d Multiple Choices Status for path %s, match %s, got %d", http.StatusMultipleChoices, path, match, entry.Status)
	} else if !multiple && entry != nil && entry.Status == http.StatusMultipleChoices {
		t.Errorf("Were not expecting %d Multiple Choices Status for path %s, match %s, but got it", http.StatusMultipleChoices, path, match)
	}
}

func TestGetEntry(t *testing.T) {
//
	testGetEntry(t, "a", "a", false, "a")
	testGetEntry(t, "b", "-", false, "a")
testGetEntry(t, "/a//
//
	testGetEntry(t, "/a", "", false, "")
	testGetEntry(t, "/a/b", "a/b", false, "a/b")
//
	testGetEntry(t, "read", "read", true, "readme.md", "readit.md")
	testGetEntry(t, "rf", "-", false, "readme.md", "readit.md")
	testGetEntry(t, "readme", "readme", false, "readme.md")
	testGetEntry(t, "readme", "-", false, "readit.md")
	testGetEntry(t, "readme.md", "readme.md", false, "readme.md")
	testGetEntry(t, "readme.md", "-", false, "readit.md")
	testGetEntry(t, "readmeAmd", "-", false, "readit.md")
	testGetEntry(t, "readme.mdffff", "-", false, "readme.md")
	testGetEntry(t, "ab", "ab", true, "ab/cefg", "ab/cedh", "ab/kkkkkk")
	testGetEntry(t, "ab/ce", "ab/ce", true, "ab/cefg", "ab/cedh", "ab/ceuuuuuuuuuu")
	testGetEntry(t, "abc", "abc", true, "abcd", "abczzzzef", "abc/def", "abc/e/g")
	testGetEntry(t, "a/b", "a/b", true, "a", "a/bc", "a/ba", "a/b/c")
	testGetEntry(t, "a/b", "a/b", false, "a", "a/b", "a/bb", "a/b/c")
testGetEntry(t, "//
}

func TestExactMatch(t *testing.T) {
	quitC := make(chan bool)
	mf := manifest("shouldBeExactMatch.css", "shouldBeExactMatch.css.map")
	fileStore := storage.NewFileStore(nil, storage.NewFileStoreParams())
	ref := make([]byte, fileStore.HashSize())
	trie, err := readManifest(mf, ref, fileStore, false, quitC, nil)
	if err != nil {
		t.Errorf("unexpected error making manifest: %v", err)
	}
	entry, _ := trie.getEntry("shouldBeExactMatch.css")
	if entry.Path != "" {
		t.Errorf("Expected entry to match %s, got: %s", "shouldBeExactMatch.css", entry.Path)
	}
	if entry.Status == http.StatusMultipleChoices {
		t.Errorf("Got status %d, which is unexepcted", http.StatusMultipleChoices)
	}
}

func TestDeleteEntry(t *testing.T) {

}

//
//
//
func TestAddFileWithManifestPath(t *testing.T) {
//
	manifest, _ := json.Marshal(&Manifest{
		Entries: []ManifestEntry{
			{Path: "ab", Hash: "ab"},
			{Path: "ac", Hash: "ac"},
		},
	})
	reader := &storage.LazyTestSectionReader{
		SectionReader: io.NewSectionReader(bytes.NewReader(manifest), 0, int64(len(manifest))),
	}
	fileStore := storage.NewFileStore(nil, storage.NewFileStoreParams())
	ref := make([]byte, fileStore.HashSize())
	trie, err := readManifest(reader, ref, fileStore, false, nil, NOOPDecrypt)
	if err != nil {
		t.Fatal(err)
	}
	checkEntry(t, "ab", "ab", false, trie)
	checkEntry(t, "ac", "ac", false, trie)

//
	entry := &manifestTrieEntry{}
	entry.Path = "a"
	entry.Hash = "a"
	trie.addEntry(entry, nil)
	checkEntry(t, "ab", "ab", false, trie)
	checkEntry(t, "ac", "ac", false, trie)
	checkEntry(t, "a", "a", false, trie)
}

//
//
//
//
//
func TestReadManifestOverSizeLimit(t *testing.T) {
	manifest := make([]byte, manifestSizeLimit+1)
	reader := &storage.LazyTestSectionReader{
		SectionReader: io.NewSectionReader(bytes.NewReader(manifest), 0, int64(len(manifest))),
	}
	_, err := readManifest(reader, storage.Address{}, nil, false, nil, NOOPDecrypt)
	if err == nil {
		t.Fatal("got no error from readManifest")
	}
//
//
	got := err.Error()
	want := fmt.Sprintf("Manifest size of %v bytes exceeds the %v byte limit", len(manifest), manifestSizeLimit)
	if got != want {
		t.Fatalf("got error mesage %q, expected %q", got, want)
	}
}

