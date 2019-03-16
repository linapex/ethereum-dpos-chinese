
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342683220971520>

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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

//
//
type ResourceMetadata struct {
StartTime Timestamp      //
Frequency uint64         //
Name      string         //
Owner     common.Address //
}

const frequencyLength = 8 //
const nameLengthLength = 1

//
//
//
//
//
//
//
const minimumMetadataLength = chunkPrefixLength + timestampLength + frequencyLength + nameLengthLength + common.AddressLength

//
func (r *ResourceMetadata) binaryGet(serializedData []byte) error {
	if len(serializedData) < minimumMetadataLength {
		return NewErrorf(ErrInvalidValue, "Metadata chunk to deserialize is too short. Expected at least %d. Got %d.", minimumMetadataLength, len(serializedData))
	}

//
	if serializedData[0] != 0 || serializedData[1] != 0 {
		return NewError(ErrCorruptData, "Invalid metadata chunk")
	}

	cursor := 2
metadataLength := int(binary.LittleEndian.Uint16(serializedData[cursor : cursor+2])) //
	if metadataLength+chunkPrefixLength != len(serializedData) {
		return NewErrorf(ErrCorruptData, "Incorrect declared metadata length. Expected %d, got %d.", metadataLength+chunkPrefixLength, len(serializedData))
	}

	cursor += 2

	if err := r.StartTime.binaryGet(serializedData[cursor : cursor+timestampLength]); err != nil {
		return err
	}
	cursor += timestampLength

	r.Frequency = binary.LittleEndian.Uint64(serializedData[cursor : cursor+frequencyLength])
	cursor += frequencyLength

	nameLength := int(serializedData[cursor])
	if nameLength+minimumMetadataLength > len(serializedData) {
		return NewErrorf(ErrInvalidValue, "Metadata chunk to deserialize is too short when decoding resource name. Expected at least %d. Got %d.", nameLength+minimumMetadataLength, len(serializedData))
	}
	cursor++
	r.Name = string(serializedData[cursor : cursor+nameLength])
	cursor += nameLength

	copy(r.Owner[:], serializedData[cursor:])
	cursor += common.AddressLength
	if cursor != len(serializedData) {
		return NewErrorf(ErrInvalidValue, "Metadata chunk has leftover data after deserialization. %d left to read", len(serializedData)-cursor)
	}
	return nil
}

//
func (r *ResourceMetadata) binaryPut(serializedData []byte) error {
	metadataChunkLength := r.binaryLength()
	if len(serializedData) != metadataChunkLength {
		return NewErrorf(ErrInvalidValue, "Need a slice of exactly %d bytes to serialize this metadata, but got a slice of size %d.", metadataChunkLength, len(serializedData))
	}

//
//
	cursor := 2
binary.LittleEndian.PutUint16(serializedData[cursor:cursor+2], uint16(metadataChunkLength-chunkPrefixLength)) //
	cursor += 2

	r.StartTime.binaryPut(serializedData[cursor : cursor+timestampLength])
	cursor += timestampLength

	binary.LittleEndian.PutUint64(serializedData[cursor:cursor+frequencyLength], r.Frequency)
	cursor += frequencyLength

//
//
	nameLength := len(r.Name)
	if nameLength > 255 {
		nameLength = 255
	}
	serializedData[cursor] = uint8(nameLength)
	cursor++
	copy(serializedData[cursor:cursor+nameLength], []byte(r.Name[:nameLength]))
	cursor += nameLength

	copy(serializedData[cursor:cursor+common.AddressLength], r.Owner[:])
	cursor += common.AddressLength

	return nil
}

func (r *ResourceMetadata) binaryLength() int {
	return minimumMetadataLength + len(r.Name)
}

//
//
func (r *ResourceMetadata) serializeAndHash() (rootAddr, metaHash []byte, chunkData []byte, err error) {

	chunkData = make([]byte, r.binaryLength())
	if err := r.binaryPut(chunkData); err != nil {
		return nil, nil, nil, err
	}
	rootAddr, metaHash = metadataHash(chunkData)
	return rootAddr, metaHash, chunkData, nil

}

//
func (metadata *ResourceMetadata) newChunk() (chunk *storage.Chunk, metaHash []byte, err error) {
//
//
//
//

//
//
//
	rootAddr, metaHash, chunkData, err := metadata.serializeAndHash()
	if err != nil {
		return nil, nil, err
	}

//
	chunk = storage.NewChunk(rootAddr, nil)
	chunk.SData = chunkData
	chunk.Size = int64(len(chunkData))

	return chunk, metaHash, nil
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
func metadataHash(chunkData []byte) (rootAddr, metaHash []byte) {
	hasher := hashPool.Get().(hash.Hash)
	defer hashPool.Put(hasher)
	hasher.Reset()
	hasher.Write(chunkData[:len(chunkData)-common.AddressLength])
	metaHash = hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(metaHash)
	hasher.Write(chunkData[len(chunkData)-common.AddressLength:])
	rootAddr = hasher.Sum(nil)
	return
}

