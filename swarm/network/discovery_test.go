
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342672198340608>

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

package network

import (
	"testing"

	p2ptest "github.com/ethereum/go-ethereum/p2p/testing"
)

/*
 
 
 
 */

func TestDiscovery(t *testing.T) {
	params := NewHiveParams()
	s, pp := newHiveTester(t, params, 1, nil)

	id := s.IDs[0]
	raddr := NewAddrFromNodeID(id)
	pp.Register([]OverlayAddr{OverlayAddr(raddr)})

//
	pp.Start(s.Server)
	defer pp.Stop()

//
	err := s.TestExchanges(p2ptest.Exchange{
		Label: "outgoing subPeersMsg",
		Expects: []p2ptest.Expect{
			{
				Code: 1,
				Msg:  &subPeersMsg{Depth: 0},
				Peer: id,
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}
}

