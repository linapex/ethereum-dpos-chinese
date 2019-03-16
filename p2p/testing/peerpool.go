
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:45</date>
//</624342661926490112>


package testing

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

type TestPeer interface {
	ID() discover.NodeID
	Drop(error)
}

//testpeerpool是演示对等连接注册的示例对等池
type TestPeerPool struct {
	lock  sync.Mutex
	peers map[discover.NodeID]TestPeer
}

func NewTestPeerPool() *TestPeerPool {
	return &TestPeerPool{peers: make(map[discover.NodeID]TestPeer)}
}

func (p *TestPeerPool) Add(peer TestPeer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	log.Trace(fmt.Sprintf("pp add peer  %v", peer.ID()))
	p.peers[peer.ID()] = peer

}

func (p *TestPeerPool) Remove(peer TestPeer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.peers, peer.ID())
}

func (p *TestPeerPool) Has(id discover.NodeID) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, ok := p.peers[id]
	return ok
}

func (p *TestPeerPool) Get(id discover.NodeID) TestPeer {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.peers[id]
}

