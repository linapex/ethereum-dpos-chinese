
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342683837534208>

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

package mru

import (
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

const (
	testDbDirName = "mru"
)

type TestHandler struct {
	*Handler
}

func (t *TestHandler) Close() {
	t.chunkStore.Close()
}

//
func NewTestHandler(datadir string, params *HandlerParams) (*TestHandler, error) {
	path := filepath.Join(datadir, testDbDirName)
	rh := NewHandler(params)
	localstoreparams := storage.NewDefaultLocalStoreParams()
	localstoreparams.Init(path)
	localStore, err := storage.NewLocalStore(localstoreparams, nil)
	if err != nil {
		return nil, fmt.Errorf("localstore create fail, path %s: %v", path, err)
	}
	localStore.Validators = append(localStore.Validators, storage.NewContentAddressValidator(storage.MakeHashFunc(resourceHashAlgorithm)))
	localStore.Validators = append(localStore.Validators, rh)
	netStore := storage.NewNetStore(localStore, nil)
	rh.SetStore(netStore)
	return &TestHandler{rh}, nil
}

