
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342670575144960>

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
//
//
//
//
//
//
//
//
package bmt

import (
	"hash"
)

//
type RefHasher struct {
maxDataLength int       //
sectionLength int       //
hasher        hash.Hash //
}

//
func NewRefHasher(hasher BaseHasherFunc, count int) *RefHasher {
	h := hasher()
	hashsize := h.Size()
	c := 2
	for ; c < count; c *= 2 {
	}
	return &RefHasher{
		sectionLength: 2 * hashsize,
		maxDataLength: c * hashsize,
		hasher:        h,
	}
}

//
//
func (rh *RefHasher) Hash(data []byte) []byte {
//
	d := make([]byte, rh.maxDataLength)
	length := len(data)
	if length > rh.maxDataLength {
		length = rh.maxDataLength
	}
	copy(d, data[:length])
	return rh.hash(d, rh.maxDataLength)
}

//
//
//
//
func (rh *RefHasher) hash(data []byte, length int) []byte {
	var section []byte
	if length == rh.sectionLength {
//
		section = data
	} else {
//
//
		length /= 2
		section = append(rh.hash(data[:length], length), rh.hash(data[length:], length)...)
	}
	rh.hasher.Reset()
	rh.hasher.Write(section)
	return rh.hasher.Sum(nil)
}

