
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342683082559488>

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
	"encoding/binary"
	"hash"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

//
//
//
type LookupParams struct {
	UpdateLookup
	Limit uint32
}

//
func (r *LookupParams) RootAddr() storage.Address {
	return r.rootAddr
}

func NewLookupParams(rootAddr storage.Address, period, version uint32, limit uint32) *LookupParams {
	return &LookupParams{
		UpdateLookup: UpdateLookup{
			period:   period,
			version:  version,
			rootAddr: rootAddr,
		},
		Limit: limit,
	}
}

//
func LookupLatest(rootAddr storage.Address) *LookupParams {
	return NewLookupParams(rootAddr, 0, 0, 0)
}

//
func LookupLatestVersionInPeriod(rootAddr storage.Address, period uint32) *LookupParams {
	return NewLookupParams(rootAddr, period, 0, 0)
}

//
func LookupVersion(rootAddr storage.Address, period, version uint32) *LookupParams {
	return NewLookupParams(rootAddr, period, version, 0)
}

//
type UpdateLookup struct {
	period   uint32
	version  uint32
	rootAddr storage.Address
}

//
//
//
const updateLookupLength = 4 + 4 + storage.KeyLength

//
func (u *UpdateLookup) UpdateAddr() (updateAddr storage.Address) {
	serializedData := make([]byte, updateLookupLength)
	u.binaryPut(serializedData)
	hasher := hashPool.Get().(hash.Hash)
	defer hashPool.Put(hasher)
	hasher.Reset()
	hasher.Write(serializedData)
	return hasher.Sum(nil)
}

//
func (u *UpdateLookup) binaryPut(serializedData []byte) error {
	if len(serializedData) != updateLookupLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to serialize UpdateLookup. Expected %d, got %d", updateLookupLength, len(serializedData))
	}
	if len(u.rootAddr) != storage.KeyLength {
		return NewError(ErrInvalidValue, "UpdateLookup.binaryPut called without rootAddr set")
	}
	binary.LittleEndian.PutUint32(serializedData[:4], u.period)
	binary.LittleEndian.PutUint32(serializedData[4:8], u.version)
	copy(serializedData[8:], u.rootAddr[:])
	return nil
}

//
func (u *UpdateLookup) binaryLength() int {
	return updateLookupLength
}

//
func (u *UpdateLookup) binaryGet(serializedData []byte) error {
	if len(serializedData) != updateLookupLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to read UpdateLookup. Expected %d, got %d", updateLookupLength, len(serializedData))
	}
	u.period = binary.LittleEndian.Uint32(serializedData[:4])
	u.version = binary.LittleEndian.Uint32(serializedData[4:8])
	u.rootAddr = storage.Address(make([]byte, storage.KeyLength))
	copy(u.rootAddr[:], serializedData[8:])
	return nil
}

