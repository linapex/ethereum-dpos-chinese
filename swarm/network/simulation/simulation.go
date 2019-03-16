
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:48</date>
//</624342674442293248>

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

package simulation

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

//
var (
	ErrNodeNotFound = errors.New("node not found")
	ErrNoPivotNode  = errors.New("no pivot node set")
)

//
//
type Simulation struct {
//
//
	Net *simulations.Network

	serviceNames []string
	cleanupFuncs []func()
	buckets      map[discover.NodeID]*sync.Map
	pivotNodeID  *discover.NodeID
	shutdownWG   sync.WaitGroup
	done         chan struct{}
	mu           sync.RWMutex

httpSrv *http.Server        //
handler *simulations.Server //
runC    chan struct{}       //
}

//
//
//
//
//
//
//
//
type ServiceFunc func(ctx *adapters.ServiceContext, bucket *sync.Map) (s node.Service, cleanup func(), err error)

//
//
func New(services map[string]ServiceFunc) (s *Simulation) {
	s = &Simulation{
		buckets: make(map[discover.NodeID]*sync.Map),
		done:    make(chan struct{}),
	}

	adapterServices := make(map[string]adapters.ServiceFunc, len(services))
	for name, serviceFunc := range services {
		s.serviceNames = append(s.serviceNames, name)
		adapterServices[name] = func(ctx *adapters.ServiceContext) (node.Service, error) {
			b := new(sync.Map)
			service, cleanup, err := serviceFunc(ctx, b)
			if err != nil {
				return nil, err
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			if cleanup != nil {
				s.cleanupFuncs = append(s.cleanupFuncs, cleanup)
			}
			s.buckets[ctx.Config.ID] = b
			return service, nil
		}
	}

	s.Net = simulations.NewNetwork(
		adapters.NewSimAdapter(adapterServices),
		&simulations.NetworkConfig{ID: "0"},
	)

	return s
}

//
//
type RunFunc func(context.Context, *Simulation) error

//
type Result struct {
	Duration time.Duration
	Error    error
}

//
//
func (s *Simulation) Run(ctx context.Context, f RunFunc) (r Result) {
//
//
	start := time.Now()
	if s.httpSrv != nil {
		log.Info("Waiting for frontend to be ready...(send POST /runsim to HTTP server)")
//
		select {
		case <-s.runC:
		case <-ctx.Done():
			return Result{
				Duration: time.Since(start),
				Error:    ctx.Err(),
			}
		}
		log.Info("Received signal from frontend - starting simulation run.")
	}
	errc := make(chan error)
	quit := make(chan struct{})
	defer close(quit)
	go func() {
		select {
		case errc <- f(ctx, s):
		case <-quit:
		}
	}()
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errc:
	}
	return Result{
		Duration: time.Since(start),
		Error:    err,
	}
}

//
//
var maxParallelCleanups = 10

//
//
//
//
//
//
func (s *Simulation) Close() {
	close(s.done)

//
//
//
	for _, c := range s.Net.Conns {
		if c.Up {
			s.Net.Disconnect(c.One, c.Other)
		}
	}
	s.shutdownWG.Wait()
	s.Net.Shutdown()

	sem := make(chan struct{}, maxParallelCleanups)
	s.mu.RLock()
	cleanupFuncs := make([]func(), len(s.cleanupFuncs))
	for i, f := range s.cleanupFuncs {
		if f != nil {
			cleanupFuncs[i] = f
		}
	}
	s.mu.RUnlock()
	var cleanupWG sync.WaitGroup
	for _, cleanup := range cleanupFuncs {
		cleanupWG.Add(1)
		sem <- struct{}{}
		go func(cleanup func()) {
			defer cleanupWG.Done()
			defer func() { <-sem }()

			cleanup()
		}(cleanup)
	}
	cleanupWG.Wait()

	if s.httpSrv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := s.httpSrv.Shutdown(ctx)
		if err != nil {
			log.Error("Error shutting down HTTP server!", "err", err)
		}
		close(s.runC)
	}
}

//
//
//
func (s *Simulation) Done() <-chan struct{} {
	return s.done
}

