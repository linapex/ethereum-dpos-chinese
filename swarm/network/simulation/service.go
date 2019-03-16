
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:48</date>
//</624342674287104000>

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
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

//
//
func (s *Simulation) Service(name string, id discover.NodeID) node.Service {
	simNode, ok := s.Net.GetNode(id).Node.(*adapters.SimNode)
	if !ok {
		return nil
	}
	services := simNode.ServiceMap()
	if len(services) == 0 {
		return nil
	}
	return services[name]
}

//
//
func (s *Simulation) RandomService(name string) node.Service {
	n := s.RandomUpNode()
	if n == nil {
		return nil
	}
	return n.Service(name)
}

//
//
func (s *Simulation) Services(name string) (services map[discover.NodeID]node.Service) {
	nodes := s.Net.GetNodes()
	services = make(map[discover.NodeID]node.Service)
	for _, node := range nodes {
		if !node.Up {
			continue
		}
		simNode, ok := node.Node.(*adapters.SimNode)
		if !ok {
			continue
		}
		services[node.ID()] = simNode.Service(name)
	}
	return services
}

