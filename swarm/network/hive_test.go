
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342672324169728>


package network

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	p2ptest "github.com/ethereum/go-ethereum/p2p/testing"
	"github.com/ethereum/go-ethereum/swarm/state"
)

func newHiveTester(t *testing.T, params *HiveParams, n int, store state.Store) (*bzzTester, *Hive) {
//设置
addr := RandomAddr() //测试的对等地址
	to := NewKademlia(addr.OAddr, NewKadParams())
pp := NewHive(params, to, store) //蜂箱

	return newBzzBaseTester(t, n, addr, DiscoverySpec, pp.Run), pp
}

func TestRegisterAndConnect(t *testing.T) {
	params := NewHiveParams()
	s, pp := newHiveTester(t, params, 1, nil)

	id := s.IDs[0]
	raddr := NewAddrFromNodeID(id)
	pp.Register([]OverlayAddr{OverlayAddr(raddr)})

//启动配置单元并等待连接
	err := pp.Start(s.Server)
	if err != nil {
		t.Fatal(err)
	}
	defer pp.Stop()
//检索和广播
	err = s.TestDisconnected(&p2ptest.Disconnect{
		Peer:  s.IDs[0],
		Error: nil,
	})

	if err == nil || err.Error() != "timed out waiting for peers to disconnect" {
		t.Fatalf("expected peer to connect")
	}
}

func TestHiveStatePersistance(t *testing.T) {
	log.SetOutput(os.Stdout)

	dir, err := ioutil.TempDir("", "hive_test_store")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

store, err := state.NewDBStore(dir) //用空的dbstore启动配置单元

	params := NewHiveParams()
	s, pp := newHiveTester(t, params, 5, store)

	peers := make(map[string]bool)
	for _, id := range s.IDs {
		raddr := NewAddrFromNodeID(id)
		pp.Register([]OverlayAddr{OverlayAddr(raddr)})
		peers[raddr.String()] = true
	}

//启动配置单元并等待连接
	err = pp.Start(s.Server)
	if err != nil {
		t.Fatal(err)
	}
	pp.Stop()
	store.Close()

persistedStore, err := state.NewDBStore(dir) //用空的dbstore启动配置单元

	s1, pp := newHiveTester(t, params, 1, persistedStore)

//启动配置单元并等待连接

	pp.Start(s1.Server)
	i := 0
	pp.Overlay.EachAddr(nil, 256, func(addr OverlayAddr, po int, nn bool) bool {
		delete(peers, addr.(*BzzAddr).String())
		i++
		return true
	})
	if len(peers) != 0 || i != 5 {
		t.Fatalf("invalid peers loaded")
	}
}

