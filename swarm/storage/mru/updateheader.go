
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342684038860800>

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
	"github.com/ethereum/go-ethereum/swarm/storage"
)

//
type updateHeader struct {
UpdateLookup        //
multihash    bool   //
metaHash     []byte //
}

const metaHashLength = storage.KeyLength

//
//
//
const updateHeaderLength = updateLookupLength + 1 + metaHashLength

//
func (h *updateHeader) binaryPut(serializedData []byte) error {
	if len(serializedData) != updateHeaderLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to serialize updateHeaderLength. Expected %d, got %d", updateHeaderLength, len(serializedData))
	}
	if len(h.metaHash) != metaHashLength {
		return NewError(ErrInvalidValue, "updateHeader.binaryPut called without metaHash set")
	}
	if err := h.UpdateLookup.binaryPut(serializedData[:updateLookupLength]); err != nil {
		return err
	}
	cursor := updateLookupLength
	copy(serializedData[cursor:], h.metaHash[:metaHashLength])
	cursor += metaHashLength

	var flags byte
	if h.multihash {
		flags |= 0x01
	}

	serializedData[cursor] = flags
	cursor++

	return nil
}

//
func (h *updateHeader) binaryLength() int {
	return updateHeaderLength
}

//
func (h *updateHeader) binaryGet(serializedData []byte) error {
	if len(serializedData) != updateHeaderLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to read updateHeaderLength. Expected %d, got %d", updateHeaderLength, len(serializedData))
	}

	if err := h.UpdateLookup.binaryGet(serializedData[:updateLookupLength]); err != nil {
		return err
	}
	cursor := updateLookupLength
	h.metaHash = make([]byte, metaHashLength)
	copy(h.metaHash[:storage.KeyLength], serializedData[cursor:cursor+storage.KeyLength])
	cursor += metaHashLength

	flags := serializedData[cursor]
	cursor++

	h.multihash = flags&0x01 != 0

	return nil
}

