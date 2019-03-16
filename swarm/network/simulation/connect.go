
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342673305636864>


package simulation

import (
	"strings"

	"github.com/ethereum/go-ethereum/p2p/discover"
)

//ConnectToPivotNode将节点与提供的节点ID连接起来
//到透视节点，已由simulation.setPivotNode方法设置。
//它在构建星型网络拓扑时很有用
//当模拟动态添加和删除节点时。
func (s *Simulation) ConnectToPivotNode(id discover.NodeID) (err error) {
	pid := s.PivotNodeID()
	if pid == nil {
		return ErrNoPivotNode
	}
	return s.connect(*pid, id)
}

//ConnectToLastNode将节点与提供的节点ID连接起来
//到上一个节点，并避免连接到自身。
//它在构建链网络拓扑结构时很有用
//当模拟动态添加和删除节点时。
func (s *Simulation) ConnectToLastNode(id discover.NodeID) (err error) {
	ids := s.UpNodeIDs()
	l := len(ids)
	if l < 2 {
		return nil
	}
	lid := ids[l-1]
	if lid == id {
		lid = ids[l-2]
	}
	return s.connect(lid, id)
}

//connecttorandomnode将节点与provided nodeid连接起来
//向上的随机节点发送。
func (s *Simulation) ConnectToRandomNode(id discover.NodeID) (err error) {
	n := s.RandomUpNode(id)
	if n == nil {
		return ErrNodeNotFound
	}
	return s.connect(n.ID, id)
}

//ConnectNodesFull将所有节点连接到另一个。
//它在网络中提供了完整的连接
//这应该是很少需要的。
func (s *Simulation) ConnectNodesFull(ids []discover.NodeID) (err error) {
	if ids == nil {
		ids = s.UpNodeIDs()
	}
	l := len(ids)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			err = s.connect(ids[i], ids[j])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//connectnodeschain连接链拓扑中的所有节点。
//如果ids参数为nil，则所有打开的节点都将被连接。
func (s *Simulation) ConnectNodesChain(ids []discover.NodeID) (err error) {
	if ids == nil {
		ids = s.UpNodeIDs()
	}
	l := len(ids)
	for i := 0; i < l-1; i++ {
		err = s.connect(ids[i], ids[i+1])
		if err != nil {
			return err
		}
	}
	return nil
}

//ConnectNodesRing连接环拓扑中的所有节点。
//如果ids参数为nil，则所有打开的节点都将被连接。
func (s *Simulation) ConnectNodesRing(ids []discover.NodeID) (err error) {
	if ids == nil {
		ids = s.UpNodeIDs()
	}
	l := len(ids)
	if l < 2 {
		return nil
	}
	for i := 0; i < l-1; i++ {
		err = s.connect(ids[i], ids[i+1])
		if err != nil {
			return err
		}
	}
	return s.connect(ids[l-1], ids[0])
}

//connectnodestar连接星形拓扑中的所有节点
//中心位于提供的节点ID。
//如果ids参数为nil，则所有打开的节点都将被连接。
func (s *Simulation) ConnectNodesStar(id discover.NodeID, ids []discover.NodeID) (err error) {
	if ids == nil {
		ids = s.UpNodeIDs()
	}
	l := len(ids)
	for i := 0; i < l; i++ {
		if id == ids[i] {
			continue
		}
		err = s.connect(id, ids[i])
		if err != nil {
			return err
		}
	}
	return nil
}

//ConnectNodessTarPivot连接星形拓扑中的所有节点
//中心位于已设置的轴节点。
//如果ids参数为nil，则所有打开的节点都将被连接。
func (s *Simulation) ConnectNodesStarPivot(ids []discover.NodeID) (err error) {
	id := s.PivotNodeID()
	if id == nil {
		return ErrNoPivotNode
	}
	return s.ConnectNodesStar(*id, ids)
}

//连接连接两个节点，但忽略已连接的错误。
func (s *Simulation) connect(oneID, otherID discover.NodeID) error {
	return ignoreAlreadyConnectedErr(s.Net.Connect(oneID, otherID))
}

func ignoreAlreadyConnectedErr(err error) error {
	if err == nil || strings.Contains(err.Error(), "already connected") {
		return nil
	}
	return err
}

