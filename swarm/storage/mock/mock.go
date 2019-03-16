
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:49</date>
//</624342682243698688>

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
//
package mock

import (
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/common"
)

//
var ErrNotFound = errors.New("not found")

//
//
type NodeStore struct {
	store GlobalStorer
	addr  common.Address
}

//
//
func NewNodeStore(addr common.Address, store GlobalStorer) *NodeStore {
	return &NodeStore{
		store: store,
		addr:  addr,
	}
}

//
//
func (n *NodeStore) Get(key []byte) (data []byte, err error) {
	return n.store.Get(n.addr, key)
}

//
//
func (n *NodeStore) Put(key []byte, data []byte) error {
	return n.store.Put(n.addr, key, data)
}

//
//
//
//
type GlobalStorer interface {
	Get(addr common.Address, key []byte) (data []byte, err error)
	Put(addr common.Address, key []byte, data []byte) error
	HasKey(addr common.Address, key []byte) bool
//
//
//
	NewNodeStore(addr common.Address) *NodeStore
}

//
//
type Importer interface {
	Import(r io.Reader) (n int, err error)
}

//
//
type Exporter interface {
	Export(w io.Writer) (n int, err error)
}

//
//
type ImportExporter interface {
	Importer
	Exporter
}

//
//
type ExportedChunk struct {
	Data  []byte           `json:"d"`
	Addrs []common.Address `json:"a"`
}

