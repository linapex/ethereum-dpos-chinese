
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:49</date>
//</624342679932637184>

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
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

//
var ErrNotFound = errors.New("ErrorNotFound")

//
var ErrInvalidArgument = errors.New("ErrorInvalidArgument")

//
type DBStore struct {
	db *leveldb.DB
}

//
func NewDBStore(path string) (s *DBStore, err error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DBStore{
		db: db,
	}, nil
}

//
//
//
func (s *DBStore) Get(key string, i interface{}) (err error) {
	has, err := s.db.Has([]byte(key), nil)
	if err != nil || !has {
		return ErrNotFound
	}

	data, err := s.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return ErrNotFound
	}

	unmarshaler, ok := i.(encoding.BinaryUnmarshaler)
	if !ok {
		return json.Unmarshal(data, i)
	}
	return unmarshaler.UnmarshalBinary(data)
}

//
func (s *DBStore) Put(key string, i interface{}) (err error) {
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

	return s.db.Put([]byte(key), bytes, nil)
}

//
func (s *DBStore) Delete(key string) (err error) {
	return s.db.Delete([]byte(key), nil)
}

//
func (s *DBStore) Close() error {
	return s.db.Close()
}

