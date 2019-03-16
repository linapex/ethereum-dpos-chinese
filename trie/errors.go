
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:51</date>
//</624342686828072960>

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

package trie

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

//
//
//
type MissingNodeError struct {
NodeHash common.Hash //
Path     []byte      //
}

func (err *MissingNodeError) Error() string {
	return fmt.Sprintf("missing trie node %x (path %x)", err.NodeHash, err.Path)
}

