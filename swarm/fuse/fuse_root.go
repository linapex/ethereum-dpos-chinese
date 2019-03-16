
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342671137181696>

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

package fuse

import (
	"bazil.org/fuse/fs"
)

var (
	_ fs.Node = (*SwarmDir)(nil)
)

type SwarmRoot struct {
	root *SwarmDir
}

func (filesystem *SwarmRoot) Root() (fs.Node, error) {
	return filesystem.root, nil
}

