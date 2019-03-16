
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:31</date>
//</624342604904927232>

package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/swarm/api"
	"github.com/ethereum/go-ethereum/swarm/api/client"
	"gopkg.in/urfave/cli.v1"
)

var salt = make([]byte, 32)

func init() {
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
}

func accessNewPass(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) != 1 {
		utils.Fatalf("Expected 1 argument - the ref")
	}

	var (
		ae        *api.AccessEntry
		accessKey []byte
		err       error
		ref       = args[0]
		password  = getPassPhrase("", 0, makePasswordList(ctx))
		dryRun    = ctx.Bool(SwarmDryRunFlag.Name)
	)
	accessKey, ae, err = api.DoPasswordNew(ctx, password, salt)
	if err != nil {
		utils.Fatalf("error getting session key: %v", err)
	}
	m, err := api.GenerateAccessControlManifest(ctx, ref, accessKey, ae)
	if dryRun {
		err = printManifests(m, nil)
		if err != nil {
			utils.Fatalf("had an error printing the manifests: %v", err)
		}
	} else {
		utils.Fatalf("uploading manifests")
		err = uploadManifests(ctx, m, nil)
		if err != nil {
			utils.Fatalf("had an error uploading the manifests: %v", err)
		}
	}
}

func accessNewPK(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) != 1 {
		utils.Fatalf("Expected 1 argument - the ref")
	}

	var (
		ae               *api.AccessEntry
		sessionKey       []byte
		err              error
		ref              = args[0]
		privateKey       = getPrivKey(ctx)
		granteePublicKey = ctx.String(SwarmAccessGrantKeyFlag.Name)
		dryRun           = ctx.Bool(SwarmDryRunFlag.Name)
	)
	sessionKey, ae, err = api.DoPKNew(ctx, privateKey, granteePublicKey, salt)
	if err != nil {
		utils.Fatalf("error getting session key: %v", err)
	}
	m, err := api.GenerateAccessControlManifest(ctx, ref, sessionKey, ae)
	if dryRun {
		err = printManifests(m, nil)
		if err != nil {
			utils.Fatalf("had an error printing the manifests: %v", err)
		}
	} else {
		err = uploadManifests(ctx, m, nil)
		if err != nil {
			utils.Fatalf("had an error uploading the manifests: %v", err)
		}
	}
}

func accessNewACT(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) != 1 {
		utils.Fatalf("Expected 1 argument - the ref")
	}

	var (
		ae          *api.AccessEntry
		actManifest *api.Manifest
		accessKey   []byte
		err         error
		ref         = args[0]
		grantees    = []string{}
		actFilename = ctx.String(SwarmAccessGrantKeysFlag.Name)
		privateKey  = getPrivKey(ctx)
		dryRun      = ctx.Bool(SwarmDryRunFlag.Name)
	)

	bytes, err := ioutil.ReadFile(actFilename)
	if err != nil {
		utils.Fatalf("had an error reading the grantee public key list")
	}
	grantees = strings.Split(string(bytes), "\n")
	accessKey, ae, actManifest, err = api.DoACTNew(ctx, privateKey, salt, grantees)
	if err != nil {
		utils.Fatalf("error generating ACT manifest: %v", err)
	}

	if err != nil {
		utils.Fatalf("error getting session key: %v", err)
	}
	m, err := api.GenerateAccessControlManifest(ctx, ref, accessKey, ae)
	if err != nil {
		utils.Fatalf("error generating root access manifest: %v", err)
	}

	if dryRun {
		err = printManifests(m, actManifest)
		if err != nil {
			utils.Fatalf("had an error printing the manifests: %v", err)
		}
	} else {
		err = uploadManifests(ctx, m, actManifest)
		if err != nil {
			utils.Fatalf("had an error uploading the manifests: %v", err)
		}
	}
}

func printManifests(rootAccessManifest, actManifest *api.Manifest) error {
	js, err := json.Marshal(rootAccessManifest)
	if err != nil {
		return err
	}
	fmt.Println(string(js))

	if actManifest != nil {
		js, err := json.Marshal(actManifest)
		if err != nil {
			return err
		}
		fmt.Println(string(js))
	}
	return nil
}

func uploadManifests(ctx *cli.Context, rootAccessManifest, actManifest *api.Manifest) error {
	bzzapi := strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
	client := client.NewClient(bzzapi)

	var (
		key string
		err error
	)
	if actManifest != nil {
		key, err = client.UploadManifest(actManifest, false)
		if err != nil {
			return err
		}

		rootAccessManifest.Entries[0].Access.Act = key
	}
	key, err = client.UploadManifest(rootAccessManifest, false)
	if err != nil {
		return err
	}
	fmt.Println(key)
	return nil
}

//makepasswordlist从global--password标志指定的文件中读取密码行
//
//此函数是utils.makepasswordlist的分支，用于查找子命令的CLI上下文。
//函数ctx.setglobal未设置可访问的全局标志值
//
func makePasswordList(ctx *cli.Context) []string {
	path := ctx.GlobalString(utils.PasswordFileFlag.Name)
	if path == "" {
		path = ctx.String(utils.PasswordFileFlag.Name)
		if path == "" {
			return nil
		}
	}
	text, err := ioutil.ReadFile(path)
	if err != nil {
		utils.Fatalf("Failed to read password file: %v", err)
	}
	lines := strings.Split(string(text), "\n")
//
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines
}

