
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342673557295104>

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
	"sync"
	"testing"
	"time"
)

//
//
//
//
func TestPeerEvents(t *testing.T) {
	sim := New(noopServiceFuncMap)
	defer sim.Close()

	_, err := sim.AddNodes(2)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	events := sim.PeerEvents(ctx, sim.NodeIDs())

//
	expectedEventCount := 2

	var wg sync.WaitGroup
	wg.Add(expectedEventCount)

	go func() {
		for e := range events {
			if e.Error != nil {
				if e.Error == context.Canceled {
					return
				}
				t.Error(e.Error)
				continue
			}
			wg.Done()
		}
	}()

	err = sim.ConnectNodesChain(sim.NodeIDs())
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

func TestPeerEventsTimeout(t *testing.T) {
	sim := New(noopServiceFuncMap)
	defer sim.Close()

	_, err := sim.AddNodes(2)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	events := sim.PeerEvents(ctx, sim.NodeIDs())

	done := make(chan struct{})
	go func() {
		for e := range events {
			if e.Error == context.Canceled {
				return
			}
			if e.Error == context.DeadlineExceeded {
				close(done)
				return
			} else {
				t.Fatal(e.Error)
			}
		}
	}()

	select {
	case <-time.After(time.Second):
		t.Error("no context deadline received")
	case <-done:
//
	}
}

