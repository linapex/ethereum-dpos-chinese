
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:31</date>
//</624342604259004416>


package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

//
func (w *wizard) deployExplorer() {
//
	if w.conf.Genesis == nil {
		log.Error("No genesis block configured")
		return
	}
	if w.conf.ethstats == "" {
		log.Error("No ethstats server configured")
		return
	}
	if w.conf.Genesis.Config.Ethash == nil {
		log.Error("Only ethash network supported")
		return
	}
//
	server := w.selectServer()
	if server == "" {
		return
	}
	client := w.servers[server]

//
	infos, err := checkExplorer(client, w.network)
	if err != nil {
		infos = &explorerInfos{
			nodePort: 30303, webPort: 80, webHost: client.server,
		}
	}
	existed := err == nil

	chainspec, err := newParityChainSpec(w.network, w.conf.Genesis, w.conf.bootnodes)
	if err != nil {
		log.Error("Failed to create chain spec for explorer", "err", err)
		return
	}
	chain, _ := json.MarshalIndent(chainspec, "", "  ")

//
	fmt.Println()
	fmt.Printf("Which port should the explorer listen on? (default = %d)\n", infos.webPort)
	infos.webPort = w.readDefaultInt(infos.webPort)

//
	if infos.webHost, err = w.ensureVirtualHost(client, infos.webPort, infos.webHost); err != nil {
		log.Error("Failed to decide on explorer host", "err", err)
		return
	}
//找出用户希望存储持久数据的位置
	fmt.Println()
	if infos.datadir == "" {
		fmt.Printf("Where should data be stored on the remote machine?\n")
		infos.datadir = w.readString()
	} else {
		fmt.Printf("Where should data be stored on the remote machine? (default = %s)\n", infos.datadir)
		infos.datadir = w.readDefaultString(infos.datadir)
	}
//
	fmt.Println()
	fmt.Printf("Which TCP/UDP port should the archive node listen on? (default = %d)\n", infos.nodePort)
	infos.nodePort = w.readDefaultInt(infos.nodePort)

//
	fmt.Println()
	if infos.ethstats == "" {
		fmt.Printf("What should the explorer be called on the stats page?\n")
		infos.ethstats = w.readString() + ":" + w.conf.ethstats
	} else {
		fmt.Printf("What should the explorer be called on the stats page? (default = %s)\n", infos.ethstats)
		infos.ethstats = w.readDefaultString(infos.ethstats) + ":" + w.conf.ethstats
	}
//
	nocache := false
	if existed {
		fmt.Println()
		fmt.Printf("Should the explorer be built from scratch (y/n)? (default = no)\n")
		nocache = w.readDefaultString("n") != "n"
	}
	if out, err := deployExplorer(client, w.network, chain, infos, nocache); err != nil {
		log.Error("Failed to deploy explorer container", "err", err)
		if len(out) > 0 {
			fmt.Printf("%s\n", out)
		}
		return
	}
//
	log.Info("Waiting for node to finish booting")
	time.Sleep(3 * time.Second)

	w.networkStats()
}

