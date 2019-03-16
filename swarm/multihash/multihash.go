
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342671808270336>

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

package multihash

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	defaultMultihashLength   = 32
	defaultMultihashTypeCode = 0x1b
)

var (
	multihashTypeCode uint8
	MultihashLength   = defaultMultihashLength
)

func init() {
	multihashTypeCode = defaultMultihashTypeCode
	MultihashLength = defaultMultihashLength
}

//
func isSwarmMultihashType(code uint8) bool {
	return code == multihashTypeCode
}

//
//
func GetMultihashLength(data []byte) (int, int, error) {
	cursor := 0
	typ, c := binary.Uvarint(data)
	if c <= 0 {
		return 0, 0, errors.New("unreadable hashtype field")
	}
	if !isSwarmMultihashType(uint8(typ)) {
		return 0, 0, fmt.Errorf("hash code %x is not a swarm hashtype", typ)
	}
	cursor += c
	hashlength, c := binary.Uvarint(data[cursor:])
	if c <= 0 {
		return 0, 0, errors.New("unreadable length field")
	}
	cursor += c

//
	inthashlength := int(hashlength)
	if len(data[c:]) < inthashlength {
		return 0, 0, errors.New("length mismatch")
	}
	return inthashlength, cursor, nil
}

//
//
func FromMultihash(data []byte) ([]byte, error) {
	hashLength, _, err := GetMultihashLength(data)
	if err != nil {
		return nil, err
	}
	return data[len(data)-hashLength:], nil
}

//
func ToMultihash(hashData []byte) []byte {
	buf := bytes.NewBuffer(nil)
	b := make([]byte, 8)
	c := binary.PutUvarint(b, uint64(multihashTypeCode))
	buf.Write(b[:c])
	c = binary.PutUvarint(b, uint64(len(hashData)))
	buf.Write(b[:c])
	buf.Write(hashData)
	return buf.Bytes()
}

