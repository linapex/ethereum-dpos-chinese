
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342670101188608>

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
	"context"
	"path"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

type Response struct {
	MimeType string
	Status   int
	Size     int64
//
	Content string
}

//
//
//
type Storage struct {
	api *API
}

func NewStorage(api *API) *Storage {
	return &Storage{api}
}

//
//
//
//
func (s *Storage) Put(ctx context.Context, content string, contentType string, toEncrypt bool) (storage.Address, func(context.Context) error, error) {
	return s.api.Put(ctx, content, contentType, toEncrypt)
}

//
//
//
//
//
//
//
//
func (s *Storage) Get(ctx context.Context, bzzpath string) (*Response, error) {
	uri, err := Parse(path.Join("bzz:/", bzzpath))
	if err != nil {
		return nil, err
	}
	addr, err := s.api.Resolve(ctx, uri.Addr)
	if err != nil {
		return nil, err
	}
	reader, mimeType, status, _, err := s.api.Get(ctx, nil, addr, uri.Path)
	if err != nil {
		return nil, err
	}
	quitC := make(chan bool)
	expsize, err := reader.Size(ctx, quitC)
	if err != nil {
		return nil, err
	}
	body := make([]byte, expsize)
	size, err := reader.Read(body)
	if int64(size) == expsize {
		err = nil
	}
	return &Response{mimeType, status, expsize, string(body[:size])}, err
}

//
//
//
//
func (s *Storage) Modify(ctx context.Context, rootHash, path, contentHash, contentType string) (newRootHash string, err error) {
	uri, err := Parse("bzz:/" + rootHash)
	if err != nil {
		return "", err
	}
	addr, err := s.api.Resolve(ctx, uri.Addr)
	if err != nil {
		return "", err
	}
	addr, err = s.api.Modify(ctx, addr, path, contentHash, contentType)
	if err != nil {
		return "", err
	}
	return addr.Hex(), nil
}

