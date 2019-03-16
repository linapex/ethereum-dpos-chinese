
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:49</date>
//</624342680083632128>

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

package state

import (
	"encoding"
	"encoding/json"
	"sync"
)

//
//
type InmemoryStore struct {
	db map[string][]byte
	mu sync.RWMutex
}

//
func NewInmemoryStore() *InmemoryStore {
	return &InmemoryStore{
		db: make(map[string][]byte),
	}
}

//
//
func (s *InmemoryStore) Get(key string, i interface{}) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bytes, ok := s.db[key]
	if !ok {
		return ErrNotFound
	}

	unmarshaler, ok := i.(encoding.BinaryUnmarshaler)
	if !ok {
		return json.Unmarshal(bytes, i)
	}

	return unmarshaler.UnmarshalBinary(bytes)
}

//
func (s *InmemoryStore) Put(key string, i interface{}) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	bytes := []byte{}

	marshaler, ok := i.(encoding.BinaryMarshaler)
	if !ok {
		if bytes, err = json.Marshal(i); err != nil {
			return err
		}
	} else {
		if bytes, err = marshaler.MarshalBinary(); err != nil {
			return err
		}
	}

	s.db[key] = bytes
	return nil
}

//
func (s *InmemoryStore) Delete(key string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.db[key]; !ok {
		return ErrNotFound
	}
	delete(s.db, key)
	return nil
}

//
func (s *InmemoryStore) Close() error {
	return nil
}

