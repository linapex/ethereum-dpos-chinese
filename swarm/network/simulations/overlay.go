
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:48</date>
//</624342674811392000>

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

//
//
//
package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
	"github.com/ethereum/go-ethereum/swarm/network"
	"github.com/ethereum/go-ethereum/swarm/state"
	colorable "github.com/mattn/go-colorable"
)

var (
	noDiscovery = flag.Bool("no-discovery", false, "disable discovery (useful if you want to load a snapshot)")
	vmodule     = flag.String("vmodule", "", "log filters for logger via Vmodule")
	verbosity   = flag.Int("verbosity", 0, "log filters for logger via Vmodule")
	httpSimPort = 8888
)

func init() {
	flag.Parse()
//
//
//
//
//
	if *vmodule != "" {
//
		glogger := log.NewGlogHandler(log.StreamHandler(colorable.NewColorableStderr(), log.TerminalFormat(true)))
		if *verbosity > 0 {
			glogger.Verbosity(log.Lvl(*verbosity))
		}
		glogger.Vmodule(*vmodule)
		log.Root().SetHandler(glogger)
	}
}

type Simulation struct {
	mtx    sync.Mutex
	stores map[discover.NodeID]*state.InmemoryStore
}

func NewSimulation() *Simulation {
	return &Simulation{
		stores: make(map[discover.NodeID]*state.InmemoryStore),
	}
}

func (s *Simulation) NewService(ctx *adapters.ServiceContext) (node.Service, error) {
	id := ctx.Config.ID
	s.mtx.Lock()
	store, ok := s.stores[id]
	if !ok {
		store = state.NewInmemoryStore()
		s.stores[id] = store
	}
	s.mtx.Unlock()

	addr := network.NewAddrFromNodeID(id)

	kp := network.NewKadParams()
	kp.MinProxBinSize = 2
	kp.MaxBinSize = 4
	kp.MinBinSize = 1
	kp.MaxRetries = 1000
	kp.RetryExponent = 2
	kp.RetryInterval = 1000000
	kad := network.NewKademlia(addr.Over(), kp)
	hp := network.NewHiveParams()
	hp.Discovery = !*noDiscovery
	hp.KeepAliveInterval = 300 * time.Millisecond

	config := &network.BzzConfig{
		OverlayAddr:  addr.Over(),
		UnderlayAddr: addr.Under(),
		HiveParams:   hp,
	}

	return network.NewBzz(config, kad, store, nil, nil), nil
}

//
func newSimulationNetwork() *simulations.Network {

	s := NewSimulation()
	services := adapters.Services{
		"overlay": s.NewService,
	}
	adapter := adapters.NewSimAdapter(services)
	simNetwork := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: "overlay",
	})
	return simNetwork
}

//
func newOverlaySim(sim *simulations.Network) *simulations.Server {
	return simulations.NewServer(sim)
}

//
func main() {
//
	runtime.GOMAXPROCS(runtime.NumCPU())
//
	runOverlaySim()
}

func runOverlaySim() {
//
	net := newSimulationNetwork()
//
	sim := newOverlaySim(net)
	log.Info(fmt.Sprintf("starting simulation server on 0.0.0.0:%d...", httpSimPort))
//
	http.ListenAndServe(fmt.Sprintf(":%d", httpSimPort), sim)
}

