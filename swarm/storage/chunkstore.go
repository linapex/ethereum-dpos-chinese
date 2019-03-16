
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:49</date>
//</624342680477896704>

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
	"context"
	"sync"
)

/*







*/

type ChunkStore interface {
Put(context.Context, *Chunk) //
	Get(context.Context, Address) (*Chunk, error)
	Close()
}

//
type MapChunkStore struct {
	chunks map[string]*Chunk
	mu     sync.RWMutex
}

func NewMapChunkStore() *MapChunkStore {
	return &MapChunkStore{
		chunks: make(map[string]*Chunk),
	}
}

func (m *MapChunkStore) Put(ctx context.Context, chunk *Chunk) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.chunks[chunk.Addr.Hex()] = chunk
	chunk.markAsStored()
}

func (m *MapChunkStore) Get(ctx context.Context, addr Address) (*Chunk, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	chunk := m.chunks[addr.Hex()]
	if chunk == nil {
		return nil, ErrChunkNotFound
	}
	return chunk, nil
}

func (m *MapChunkStore) Close() {
}

