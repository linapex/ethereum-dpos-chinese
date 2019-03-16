
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:48</date>
//</624342674933026816>

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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/swarm/log"
)

var (
	nodeCount = 16
)

//
//
//
//
//
//
//
func TestOverlaySim(t *testing.T) {
t.Skip("Test is flaky, see: https://
//
	log.Info("Start simulation backend")
//
	net := newSimulationNetwork()
//
	sim := newOverlaySim(net)
//
	srv := httptest.NewServer(sim)
	defer srv.Close()

	log.Debug("Http simulation server started. Start simulation network")
//
	resp, err := http.Post(srv.URL+"/start", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected Status Code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	log.Debug("Start mocker")
//
	resp, err = http.PostForm(srv.URL+"/mocker/start",
		url.Values{
			"node-count":  {fmt.Sprintf("%d", nodeCount)},
			"mocker-type": {simulations.GetMockerList()[0]},
		})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		reason, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("Expected Status Code %d, got %d, response body %s", http.StatusOK, resp.StatusCode, string(reason))
	}

//
	var upCount int
	trigger := make(chan discover.NodeID)

//
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

//
	go watchSimEvents(net, ctx, trigger)

//
LOOP:
	for {
		select {
		case <-trigger:
//
			upCount++
//
			if upCount == nodeCount {
				break LOOP
			}
		case <-ctx.Done():
			t.Fatalf("Timed out waiting for up events")
		}

	}

//
	log.Info("Get number of nodes")
//
	resp, err = http.Get(srv.URL + "/nodes")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("err %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

//
	var nodesArr []simulations.Node
	err = json.Unmarshal(b, &nodesArr)
	if err != nil {
		t.Fatal(err)
	}

//
	if len(nodesArr) != nodeCount {
		t.Fatal(fmt.Errorf("Expected %d number of nodes, got %d", nodeCount, len(nodesArr)))
	}

//
//
	time.Sleep(1 * time.Second)

	log.Info("Stop the network")
//
	resp, err = http.Post(srv.URL+"/stop", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("err %s", resp.Status)
	}

	log.Info("Reset the network")
//
	resp, err = http.Post(srv.URL+"/reset", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("err %s", resp.Status)
	}
}

//
func watchSimEvents(net *simulations.Network, ctx context.Context, trigger chan discover.NodeID) {
	events := make(chan *simulations.Event)
	sub := net.Events().Subscribe(events)
	defer sub.Unsubscribe()

	for {
		select {
		case ev := <-events:
//
			if ev.Type == simulations.EventTypeNode {
				if ev.Node.Up {
					log.Debug("got node up event", "event", ev, "node", ev.Node.Config.ID)
					select {
					case trigger <- ev.Node.Config.ID:
					case <-ctx.Done():
						return
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

