
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342673888645120>

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

package simulation

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/p2p/discover"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/swarm/network"
)

//
//
var BucketKeyKademlia BucketKey = "kademlia"

//
//
func (s *Simulation) WaitTillHealthy(ctx context.Context, kadMinProxSize int) (ill map[discover.NodeID]*network.Kademlia, err error) {
//
	var ppmap map[string]*network.PeerPot
	kademlias := s.kademlias()
	addrs := make([][]byte, 0, len(kademlias))
	for _, k := range kademlias {
		addrs = append(addrs, k.BaseAddr())
	}
	ppmap = network.NewPeerPotMap(kadMinProxSize, addrs)

//
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	ill = make(map[discover.NodeID]*network.Kademlia)
	for {
		select {
		case <-ctx.Done():
			return ill, ctx.Err()
		case <-ticker.C:
			for k := range ill {
				delete(ill, k)
			}
			log.Debug("kademlia health check", "addr count", len(addrs))
			for id, k := range kademlias {
//
				addr := common.Bytes2Hex(k.BaseAddr())
				pp := ppmap[addr]
//
				h := k.Healthy(pp)
//
				log.Debug(k.String())
				log.Debug("kademlia", "empty bins", pp.EmptyBins, "gotNN", h.GotNN, "knowNN", h.KnowNN, "full", h.Full)
				log.Debug("kademlia", "health", h.GotNN && h.KnowNN && h.Full, "addr", hex.EncodeToString(k.BaseAddr()), "node", id)
				log.Debug("kademlia", "ill condition", !h.GotNN || !h.Full, "addr", hex.EncodeToString(k.BaseAddr()), "node", id)
				if !h.GotNN || !h.Full {
					ill[id] = k
				}
			}
			if len(ill) == 0 {
				return nil, nil
			}
		}
	}
}

//
//
func (s *Simulation) kademlias() (ks map[discover.NodeID]*network.Kademlia) {
	items := s.UpNodesItems(BucketKeyKademlia)
	ks = make(map[discover.NodeID]*network.Kademlia, len(items))
	for id, v := range items {
		k, ok := v.(*network.Kademlia)
		if !ok {
			continue
		}
		ks[id] = k
	}
	return ks
}

