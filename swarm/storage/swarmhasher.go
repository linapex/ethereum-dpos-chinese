
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342684449902592>

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

package storage

import (
	"hash"
)

const (
	BMTHash     = "BMT"
SHA3Hash    = "SHA3" //
	DefaultHash = BMTHash
)

type SwarmHash interface {
	hash.Hash
	ResetWithLength([]byte)
}

type HashWithLength struct {
	hash.Hash
}

func (h *HashWithLength) ResetWithLength(length []byte) {
	h.Reset()
	h.Write(length)
}

