
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:31</date>
//</624342604443553792>


package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/log"
)

//
func makeWizard(network string) *wizard {
	return &wizard{
		network: network,
		conf: config{
			Servers: make(map[string][]byte),
		},
		servers:  make(map[string]*sshClient),
		services: make(map[string][]string),
		in:       bufio.NewReader(os.Stdin),
	}
}

//
//
func (w *wizard) run() {
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("| Welcome to puppeth, your Ethereum private network manager |")
	fmt.Println("|                                                           |")
	fmt.Println("| This tool lets you create a new Ethereum network down to  |")
	fmt.Println("| the genesis block, bootnodes, miners and ethstats servers |")
	fmt.Println("| without the hassle that it would normally entail.         |")
	fmt.Println("|                                                           |")
	fmt.Println("| Puppeth uses SSH to dial in to remote servers, and builds |")
	fmt.Println("| its network components out of Docker containers using the |")
	fmt.Println("| docker-compose toolset.                                   |")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println()

//
//
	if w.network == "" {
		fmt.Println("Please specify a network name to administer (no spaces or hyphens, please)")
		for {
			w.network = w.readString()
			if !strings.Contains(w.network, " ") && !strings.Contains(w.network, "-") {
				fmt.Printf("\nSweet, you can set this via --network=%s next time!\n\n", w.network)
				break
			}
			log.Error("I also like to live dangerously, still no spaces or hyphens")
		}
	}
	log.Info("Administering Ethereum network", "name", w.network)

//
	w.conf.path = filepath.Join(os.Getenv("HOME"), ".puppeth", w.network)

	blob, err := ioutil.ReadFile(w.conf.path)
	if err != nil {
		log.Warn("No previous configurations found", "path", w.conf.path)
	} else if err := json.Unmarshal(blob, &w.conf); err != nil {
		log.Crit("Previous configuration corrupted", "path", w.conf.path, "err", err)
	} else {
//
		var pend sync.WaitGroup
		for server, pubkey := range w.conf.Servers {
			pend.Add(1)

			go func(server string, pubkey []byte) {
				defer pend.Done()

				log.Info("Dialing previously configured server", "server", server)
				client, err := dial(server, pubkey)
				if err != nil {
					log.Error("Previous server unreachable", "server", server, "err", err)
				}
				w.lock.Lock()
				w.servers[server] = client
				w.lock.Unlock()
			}(server, pubkey)
		}
		pend.Wait()
		w.networkStats()
	}
//
	for {
		fmt.Println()
		fmt.Println("What would you like to do? (default = stats)")
		fmt.Println(" 1. Show network stats")
		if w.conf.Genesis == nil {
			fmt.Println(" 2. Configure new genesis")
		} else {
			fmt.Println(" 2. Manage existing genesis")
		}
		if len(w.servers) == 0 {
			fmt.Println(" 3. Track new remote server")
		} else {
			fmt.Println(" 3. Manage tracked machines")
		}
		if len(w.services) == 0 {
			fmt.Println(" 4. Deploy network components")
		} else {
			fmt.Println(" 4. Manage network components")
		}

		choice := w.read()
		switch {
		case choice == "" || choice == "1":
			w.networkStats()

		case choice == "2":
			if w.conf.Genesis == nil {
				w.makeGenesis()
			} else {
				w.manageGenesis()
			}
		case choice == "3":
			if len(w.servers) == 0 {
				if w.makeServer() != "" {
					w.networkStats()
				}
			} else {
				w.manageServers()
			}
		case choice == "4":
			if len(w.services) == 0 {
				w.deployComponent()
			} else {
				w.manageComponents()
			}

		default:
			log.Error("That's not something I can do")
		}
	}
}

