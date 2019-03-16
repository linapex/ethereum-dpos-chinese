
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:50</date>
//</624342684219215872>

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
	"time"

	"github.com/ethereum/go-ethereum/swarm/log"
	"github.com/ethereum/go-ethereum/swarm/spancontext"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
//
//
//
	netStoreRetryTimeout = 30 * time.Second
//
//
//
	netStoreMinRetryDelay = 3 * time.Second
//
//
//
	searchTimeout = 10 * time.Second
)

//
//
//
//
type NetStore struct {
	localStore *LocalStore
	retrieve   func(ctx context.Context, chunk *Chunk) error
}

func NewNetStore(localStore *LocalStore, retrieve func(ctx context.Context, chunk *Chunk) error) *NetStore {
	return &NetStore{localStore, retrieve}
}

//
//
//
//
//
//
func (ns *NetStore) Get(ctx context.Context, addr Address) (chunk *Chunk, err error) {

	var sp opentracing.Span
	ctx, sp = spancontext.StartSpan(
		ctx,
		"netstore.get.global")
	defer sp.Finish()

	timer := time.NewTimer(netStoreRetryTimeout)
	defer timer.Stop()

//
//
	type result struct {
		chunk *Chunk
		err   error
	}
	resultC := make(chan result)

//
//
	quitC := make(chan struct{})
	defer close(quitC)

//
//
	go func() {
//
//
//
//
		limiter := time.NewTimer(netStoreMinRetryDelay)
		defer limiter.Stop()

		for {
			chunk, err := ns.get(ctx, addr, 0)
			if err != ErrChunkNotFound {
//
//
				select {
				case <-quitC:
//
//
//
				case resultC <- result{chunk: chunk, err: err}:
//
				}
				return

			}
			select {
			case <-quitC:
//
//
//
				return
			case <-limiter.C:
			}
//
			limiter.Reset(netStoreMinRetryDelay)
			log.Debug("NetStore.Get retry chunk", "key", addr)
		}
	}()

	select {
	case r := <-resultC:
		return r.chunk, r.err
	case <-timer.C:
		return nil, ErrChunkNotFound
	}
}

//
func (ns *NetStore) GetWithTimeout(ctx context.Context, addr Address, timeout time.Duration) (chunk *Chunk, err error) {
	return ns.get(ctx, addr, timeout)
}

func (ns *NetStore) get(ctx context.Context, addr Address, timeout time.Duration) (chunk *Chunk, err error) {
	if timeout == 0 {
		timeout = searchTimeout
	}

	var sp opentracing.Span
	ctx, sp = spancontext.StartSpan(
		ctx,
		"netstore.get")
	defer sp.Finish()

	if ns.retrieve == nil {
		chunk, err = ns.localStore.Get(ctx, addr)
		if err == nil {
			return chunk, nil
		}
		if err != ErrFetching {
			return nil, err
		}
	} else {
		var created bool
		chunk, created = ns.localStore.GetOrCreateRequest(ctx, addr)

		if chunk.ReqC == nil {
			return chunk, nil
		}

		if created {
			err := ns.retrieve(ctx, chunk)
			if err != nil {
//
				chunk.SetErrored(ErrChunkUnavailable)
				return nil, err
			}
		}
	}

	t := time.NewTicker(timeout)
	defer t.Stop()

	select {
	case <-t.C:
//
		chunk.SetErrored(ErrChunkNotFound)
		return nil, ErrChunkNotFound
	case <-chunk.ReqC:
	}
	chunk.SetErrored(nil)
	return chunk, nil
}

//
func (ns *NetStore) Put(ctx context.Context, chunk *Chunk) {
	ns.localStore.Put(ctx, chunk)
}

//
func (ns *NetStore) Close() {
	ns.localStore.Close()
}

