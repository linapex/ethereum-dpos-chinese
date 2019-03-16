
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:46</date>
//</624342668872257536>

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
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/swarm/storage/encryption"
)

type RefEncryption struct {
	spanEncryption encryption.Encryption
	dataEncryption encryption.Encryption
	span           []byte
}

func NewRefEncryption(refSize int) *RefEncryption {
	span := make([]byte, 8)
	binary.LittleEndian.PutUint64(span, uint64(refSize))
	return &RefEncryption{
		spanEncryption: encryption.New(0, uint32(refSize/32), sha3.NewKeccak256),
		dataEncryption: encryption.New(refSize, 0, sha3.NewKeccak256),
		span:           span,
	}
}

func (re *RefEncryption) Encrypt(ref []byte, key []byte) ([]byte, error) {
	encryptedSpan, err := re.spanEncryption.Encrypt(re.span, key)
	if err != nil {
		return nil, err
	}
	encryptedData, err := re.dataEncryption.Encrypt(ref, key)
	if err != nil {
		return nil, err
	}
	encryptedRef := make([]byte, len(ref)+8)
	copy(encryptedRef[:8], encryptedSpan)
	copy(encryptedRef[8:], encryptedData)

	return encryptedRef, nil
}

func (re *RefEncryption) Decrypt(ref []byte, key []byte) ([]byte, error) {
	decryptedSpan, err := re.spanEncryption.Decrypt(ref[:8], key)
	if err != nil {
		return nil, err
	}

	size := binary.LittleEndian.Uint64(decryptedSpan)
	if size != uint64(len(ref)-8) {
		return nil, errors.New("invalid span in encrypted reference")
	}

	decryptedRef, err := re.dataEncryption.Decrypt(ref[8:], key)
	if err != nil {
		return nil, err
	}

	return decryptedRef, nil
}

