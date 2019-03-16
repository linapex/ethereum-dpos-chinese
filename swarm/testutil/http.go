
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342684806418432>

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

package testutil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/swarm/api"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/ethereum/go-ethereum/swarm/storage/mru"
)

type TestServer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

//
type fakeTimeProvider struct {
	currentTime uint64
}

func (f *fakeTimeProvider) Tick() {
	f.currentTime++
}

func (f *fakeTimeProvider) Now() mru.Timestamp {
	return mru.Timestamp{Time: f.currentTime}
}

func NewTestSwarmServer(t *testing.T, serverFunc func(*api.API) TestServer) *TestSwarmServer {
	dir, err := ioutil.TempDir("", "swarm-storage-test")
	if err != nil {
		t.Fatal(err)
	}
	storeparams := storage.NewDefaultLocalStoreParams()
	storeparams.DbCapacity = 5000000
	storeparams.CacheCapacity = 5000
	storeparams.Init(dir)
	localStore, err := storage.NewLocalStore(storeparams, nil)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	fileStore := storage.NewFileStore(localStore, storage.NewFileStoreParams())

//
	resourceDir, err := ioutil.TempDir("", "swarm-resource-test")
	if err != nil {
		t.Fatal(err)
	}

	fakeTimeProvider := &fakeTimeProvider{
		currentTime: 42,
	}
	mru.TimestampProvider = fakeTimeProvider
	rhparams := &mru.HandlerParams{}
	rh, err := mru.NewTestHandler(resourceDir, rhparams)
	if err != nil {
		t.Fatal(err)
	}

	a := api.NewAPI(fileStore, nil, rh.Handler, nil)
	srv := httptest.NewServer(serverFunc(a))
	return &TestSwarmServer{
		Server:            srv,
		FileStore:         fileStore,
		dir:               dir,
		Hasher:            storage.MakeHashFunc(storage.DefaultHash)(),
		timestampProvider: fakeTimeProvider,
		cleanup: func() {
			srv.Close()
			rh.Close()
			os.RemoveAll(dir)
			os.RemoveAll(resourceDir)
		},
	}
}

type TestSwarmServer struct {
	*httptest.Server
	Hasher            storage.SwarmHash
	FileStore         *storage.FileStore
	dir               string
	cleanup           func()
	timestampProvider *fakeTimeProvider
}

func (t *TestSwarmServer) Close() {
	t.cleanup()
}

func (t *TestSwarmServer) GetCurrentTime() mru.Timestamp {
	return t.timestampProvider.Now()
}

