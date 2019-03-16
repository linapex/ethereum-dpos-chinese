
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342683581681664>

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

package mru

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const signatureLength = 65

//
type Signature [signatureLength]byte

//
type Signer interface {
	Sign(common.Hash) (Signature, error)
	Address() common.Address
}

//
//
type GenericSigner struct {
	PrivKey *ecdsa.PrivateKey
	address common.Address
}

//
func NewGenericSigner(privKey *ecdsa.PrivateKey) *GenericSigner {
	return &GenericSigner{
		PrivKey: privKey,
		address: crypto.PubkeyToAddress(privKey.PublicKey),
	}
}

//
//
func (s *GenericSigner) Sign(data common.Hash) (signature Signature, err error) {
	signaturebytes, err := crypto.Sign(data.Bytes(), s.PrivKey)
	if err != nil {
		return
	}
	copy(signature[:], signaturebytes)
	return
}

//
func (s *GenericSigner) Address() common.Address {
	return s.address
}

