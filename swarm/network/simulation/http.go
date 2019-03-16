
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342673725067264>

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
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/simulations"
)

//
var (
	DefaultHTTPSimAddr = ":8888"
)

//
//
func (s *Simulation) WithServer(addr string) *Simulation {
//
	if addr == "" {
		addr = DefaultHTTPSimAddr
	}
	log.Info(fmt.Sprintf("Initializing simulation server on %s...", addr))
//
	s.handler = simulations.NewServer(s.Net)
	s.runC = make(chan struct{})
//
	s.addSimulationRoutes()
	s.httpSrv = &http.Server{
		Addr:    addr,
		Handler: s.handler,
	}
	go func() {
		err := s.httpSrv.ListenAndServe()
		if err != nil {
			log.Error("Error starting the HTTP server", "error", err)
		}
	}()
	return s
}

//
func (s *Simulation) addSimulationRoutes() {
	s.handler.POST("/runsim", s.RunSimulation)
}

//
func (s *Simulation) RunSimulation(w http.ResponseWriter, req *http.Request) {
	log.Debug("RunSimulation endpoint running")
	s.runC <- struct{}{}
	w.WriteHeader(http.StatusOK)
}

