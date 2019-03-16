
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342670222823424>

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

package api

import (
	"github.com/ethereum/go-ethereum/swarm/network"
)

type Control struct {
	api  *API
	hive *network.Hive
}

func NewControl(api *API, hive *network.Hive) *Control {
	return &Control{api, hive}
}

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
func (c *Control) Hive() string {
	return c.hive.String()
}

