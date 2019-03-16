
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:31</date>
//</624342604724572160>


package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

//
func (w *wizard) deployNode(boot bool) {
//
	if w.conf.Genesis == nil {
		log.Error("No genesis block configured")
		return
	}
	if w.conf.ethstats == "" {
		log.Error("No ethstats server configured")
		return
	}
//
	server := w.selectServer()
	if server == "" {
		return
	}
	client := w.servers[server]

//
	infos, err := checkNode(client, w.network, boot)
	if err != nil {
		if boot {
			infos = &nodeInfos{port: 30303, peersTotal: 512, peersLight: 256}
		} else {
			infos = &nodeInfos{port: 30303, peersTotal: 50, peersLight: 0, gasTarget: 4.7, gasPrice: 18}
		}
	}
	existed := err == nil

	infos.genesis, _ = json.MarshalIndent(w.conf.Genesis, "", "  ")
	infos.network = w.conf.Genesis.Config.ChainID.Int64()

//找出用户希望存储持久数据的位置
	fmt.Println()
	if infos.datadir == "" {
		fmt.Printf("Where should data be stored on the remote machine?\n")
		infos.datadir = w.readString()
	} else {
		fmt.Printf("Where should data be stored on the remote machine? (default = %s)\n", infos.datadir)
		infos.datadir = w.readDefaultString(infos.datadir)
	}
	if w.conf.Genesis.Config.Ethash != nil && !boot {
		fmt.Println()
		if infos.ethashdir == "" {
			fmt.Printf("Where should the ethash mining DAGs be stored on the remote machine?\n")
			infos.ethashdir = w.readString()
		} else {
			fmt.Printf("Where should the ethash mining DAGs be stored on the remote machine? (default = %s)\n", infos.ethashdir)
			infos.ethashdir = w.readDefaultString(infos.ethashdir)
		}
	}
//找出要监听的端口
	fmt.Println()
	fmt.Printf("Which TCP/UDP port to listen on? (default = %d)\n", infos.port)
	infos.port = w.readDefaultInt(infos.port)

//计算允许多少对等点（根据节点类型不同）
	fmt.Println()
	fmt.Printf("How many peers to allow connecting? (default = %d)\n", infos.peersTotal)
	infos.peersTotal = w.readDefaultInt(infos.peersTotal)

//计算允许多少光对等（根据节点类型不同）
	fmt.Println()
	fmt.Printf("How many light peers to allow connecting? (default = %d)\n", infos.peersLight)
	infos.peersLight = w.readDefaultInt(infos.peersLight)

//
	fmt.Println()
	if infos.ethstats == "" {
		fmt.Printf("What should the node be called on the stats page?\n")
		infos.ethstats = w.readString() + ":" + w.conf.ethstats
	} else {
		fmt.Printf("What should the node be called on the stats page? (default = %s)\n", infos.ethstats)
		infos.ethstats = w.readDefaultString(infos.ethstats) + ":" + w.conf.ethstats
	}
//
	if !boot {
		if w.conf.Genesis.Config.Ethash != nil {
//Ethash based miners only need an etherbase to mine against
			fmt.Println()
			if infos.etherbase == "" {
				fmt.Printf("What address should the miner use?\n")
				for {
					if address := w.readAddress(); address != nil {
						infos.etherbase = address.Hex()
						break
					}
				}
			} else {
				fmt.Printf("What address should the miner use? (default = %s)\n", infos.etherbase)
				infos.etherbase = w.readDefaultAddress(common.HexToAddress(infos.etherbase)).Hex()
			}
		} else if w.conf.Genesis.Config.Clique != nil {
//
			if infos.keyJSON != "" {
				if key, err := keystore.DecryptKey([]byte(infos.keyJSON), infos.keyPass); err != nil {
					infos.keyJSON, infos.keyPass = "", ""
				} else {
					fmt.Println()
					fmt.Printf("Reuse previous (%s) signing account (y/n)? (default = yes)\n", key.Address.Hex())
					if w.readDefaultString("y") != "y" {
						infos.keyJSON, infos.keyPass = "", ""
					}
				}
			}
//基于集团的签名者需要一个密钥文件和解锁密码，询问是否不可用
			if infos.keyJSON == "" {
				fmt.Println()
				fmt.Println("Please paste the signer's key JSON:")
				infos.keyJSON = w.readJSON()

				fmt.Println()
				fmt.Println("What's the unlock password for the account? (won't be echoed)")
				infos.keyPass = w.readPassword()

				if _, err := keystore.DecryptKey([]byte(infos.keyJSON), infos.keyPass); err != nil {
					log.Error("Failed to decrypt key with given passphrase")
					return
				}
			}
		}
//
		fmt.Println()
		fmt.Printf("What gas limit should empty blocks target (MGas)? (default = %0.3f)\n", infos.gasTarget)
		infos.gasTarget = w.readDefaultFloat(infos.gasTarget)

		fmt.Println()
		fmt.Printf("What gas price should the signer require (GWei)? (default = %0.3f)\n", infos.gasPrice)
		infos.gasPrice = w.readDefaultFloat(infos.gasPrice)
	}
//
	nocache := false
	if existed {
		fmt.Println()
		fmt.Printf("Should the node be built from scratch (y/n)? (default = no)\n")
		nocache = w.readDefaultString("n") != "n"
	}
	if out, err := deployNode(client, w.network, w.conf.bootnodes, infos, nocache); err != nil {
		log.Error("Failed to deploy Ethereum node container", "err", err)
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

