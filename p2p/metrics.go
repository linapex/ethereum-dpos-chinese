
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:44</date>
//</624342658625572864>


//包含网络层使用的仪表和计时器。

package p2p

import (
	"net"

	"github.com/ethereum/go-ethereum/metrics"
)

var (
	ingressConnectMeter = metrics.NewRegisteredMeter("p2p/InboundConnects", nil)
	ingressTrafficMeter = metrics.NewRegisteredMeter("p2p/InboundTraffic", nil)
	egressConnectMeter  = metrics.NewRegisteredMeter("p2p/OutboundConnects", nil)
	egressTrafficMeter  = metrics.NewRegisteredMeter("p2p/OutboundTraffic", nil)
)

//meteredconn是一个围绕net.conn的包装器，用于测量
//入站和出站网络流量。
type meteredConn struct {
net.Conn //与计量包装的网络连接
}

//NewMeteredConn创建了一个新的Metered连接，也会影响入口或
//出口连接仪表。如果禁用度量系统，则此函数
//返回原始对象。
func newMeteredConn(conn net.Conn, ingress bool) net.Conn {
//如果禁用度量值，则短路
	if !metrics.Enabled {
		return conn
	}
//否则，请触发连接计数器并包装连接
	if ingress {
		ingressConnectMeter.Mark(1)
	} else {
		egressConnectMeter.Mark(1)
	}
	return &meteredConn{Conn: conn}
}

//读取将网络读取委派给基础连接，从而阻止进入
//一路上的交通表。
func (c *meteredConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	ingressTrafficMeter.Mark(int64(n))
	return
}

//写入将网络写入委托给基础连接，将
//一路上的出口流量表。
func (c *meteredConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	egressTrafficMeter.Mark(int64(n))
	return
}

