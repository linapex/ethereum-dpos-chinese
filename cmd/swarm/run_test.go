
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:31</date>
//</624342606255493120>


package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/docker/docker/pkg/reexec"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/internal/cmdtest"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/swarm"
)

func init() {
//如果我们在run swarm中被执行为“swarm测试”，就运行这个应用程序。
	reexec.Register("swarm-test", func() {
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}

func TestMain(m *testing.M) {
//检查我们是否被重新执行了
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

func runSwarm(t *testing.T, args ...string) *cmdtest.TestCmd {
	tt := cmdtest.NewTestCmd(t, nil)

//
//函数将阻止任何测试运行。
	tt.Run("swarm-test", args...)

	return tt
}

type testCluster struct {
	Nodes  []*testNode
	TmpDir string
}

//
//
//
//它。
//
//
//
//
//
//
//

func newTestCluster(t *testing.T, size int) *testCluster {
	cluster := &testCluster{}
	defer func() {
		if t.Failed() {
			cluster.Shutdown()
		}
	}()

	tmpdir, err := ioutil.TempDir("", "swarm-test")
	if err != nil {
		t.Fatal(err)
	}
	cluster.TmpDir = tmpdir

//
	cluster.StartNewNodes(t, size)

	if size == 1 {
		return cluster
	}

//
	for _, node := range cluster.Nodes {
		if err := node.Client.Call(nil, "admin_addPeer", cluster.Nodes[0].Enode); err != nil {
			t.Fatal(err)
		}
	}

//
outer:
	for _, node := range cluster.Nodes {
		var peers []*p2p.PeerInfo
		for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(50 * time.Millisecond) {
			if err := node.Client.Call(&peers, "admin_peers"); err != nil {
				t.Fatal(err)
			}
			if len(peers) == len(cluster.Nodes)-1 {
				continue outer
			}
		}
		t.Fatalf("%s only has %d / %d peers", node.Name, len(peers), len(cluster.Nodes)-1)
	}

	return cluster
}

func (c *testCluster) Shutdown() {
	for _, node := range c.Nodes {
		node.Shutdown()
	}
	os.RemoveAll(c.TmpDir)
}

func (c *testCluster) Stop() {
	for _, node := range c.Nodes {
		node.Shutdown()
	}
}

func (c *testCluster) StartNewNodes(t *testing.T, size int) {
	c.Nodes = make([]*testNode, 0, size)
	for i := 0; i < size; i++ {
		dir := filepath.Join(c.TmpDir, fmt.Sprintf("swarm%02d", i))
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}

		node := newTestNode(t, dir)
		node.Name = fmt.Sprintf("swarm%02d", i)

		c.Nodes = append(c.Nodes, node)
	}
}

func (c *testCluster) StartExistingNodes(t *testing.T, size int, bzzaccount string) {
	c.Nodes = make([]*testNode, 0, size)
	for i := 0; i < size; i++ {
		dir := filepath.Join(c.TmpDir, fmt.Sprintf("swarm%02d", i))
		node := existingTestNode(t, dir, bzzaccount)
		node.Name = fmt.Sprintf("swarm%02d", i)

		c.Nodes = append(c.Nodes, node)
	}
}

func (c *testCluster) Cleanup() {
	os.RemoveAll(c.TmpDir)
}

type testNode struct {
	Name       string
	Addr       string
	URL        string
	Enode      string
	Dir        string
	IpcPath    string
	PrivateKey *ecdsa.PrivateKey
	Client     *rpc.Client
	Cmd        *cmdtest.TestCmd
}

const testPassphrase = "swarm-test-passphrase"

func getTestAccount(t *testing.T, dir string) (conf *node.Config, account accounts.Account) {
//创建密钥
	conf = &node.Config{
		DataDir: dir,
		IPCPath: "bzzd.ipc",
		NoUSB:   true,
	}
	n, err := node.New(conf)
	if err != nil {
		t.Fatal(err)
	}
	account, err = n.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore).NewAccount(testPassphrase)
	if err != nil {
		t.Fatal(err)
	}

//在Windows上运行测试时使用唯一的ipcpath
	if runtime.GOOS == "windows" {
		conf.IPCPath = fmt.Sprintf("bzzd-%s.ipc", account.Address.String())
	}

	return conf, account
}

func existingTestNode(t *testing.T, dir string, bzzaccount string) *testNode {
	conf, _ := getTestAccount(t, dir)
	node := &testNode{Dir: dir}

//在Windows上运行测试时使用唯一的ipcpath
	if runtime.GOOS == "windows" {
		conf.IPCPath = fmt.Sprintf("bzzd-%s.ipc", bzzaccount)
	}

//指定端口
	ports, err := getAvailableTCPPorts(2)
	if err != nil {
		t.Fatal(err)
	}
	p2pPort := ports[0]
	httpPort := ports[1]

//启动节点
	node.Cmd = runSwarm(t,
		"--port", p2pPort,
		"--nodiscover",
		"--datadir", dir,
		"--ipcpath", conf.IPCPath,
		"--ens-api", "",
		"--bzzaccount", bzzaccount,
		"--bzznetworkid", "321",
		"--bzzport", httpPort,
		"--verbosity", "6",
	)
	node.Cmd.InputLine(testPassphrase)
	defer func() {
		if t.Failed() {
			node.Shutdown()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

//确保所有端口都有活动的侦听器
//这样下一个节点就不会得到相同的
//调用GetAvailableTCPPorts时
	err = waitTCPPorts(ctx, ports...)
	if err != nil {
		t.Fatal(err)
	}

//等待节点启动
	for start := time.Now(); time.Since(start) < 10*time.Second; time.Sleep(50 * time.Millisecond) {
		node.Client, err = rpc.Dial(conf.IPCEndpoint())
		if err == nil {
			break
		}
	}
	if node.Client == nil {
		t.Fatal(err)
	}

//加载信息
	var info swarm.Info
	if err := node.Client.Call(&info, "bzz_info"); err != nil {
		t.Fatal(err)
	}
	node.Addr = net.JoinHostPort("127.0.0.1", info.Port)
node.URL = "http://“+No.ADDR”

	var nodeInfo p2p.NodeInfo
	if err := node.Client.Call(&nodeInfo, "admin_nodeInfo"); err != nil {
		t.Fatal(err)
	}
node.Enode = fmt.Sprintf("enode://%s@127.0.0.1:%s”，nodeinfo.id，p2pport）

	return node
}

func newTestNode(t *testing.T, dir string) *testNode {

	conf, account := getTestAccount(t, dir)
	ks := keystore.NewKeyStore(path.Join(dir, "keystore"), 1<<18, 1)

	pk := decryptStoreAccount(ks, account.Address.Hex(), []string{testPassphrase})

	node := &testNode{Dir: dir, PrivateKey: pk}

//指定端口
	ports, err := getAvailableTCPPorts(2)
	if err != nil {
		t.Fatal(err)
	}
	p2pPort := ports[0]
	httpPort := ports[1]

//启动节点
	node.Cmd = runSwarm(t,
		"--port", p2pPort,
		"--nodiscover",
		"--datadir", dir,
		"--ipcpath", conf.IPCPath,
		"--ens-api", "",
		"--bzzaccount", account.Address.String(),
		"--bzznetworkid", "321",
		"--bzzport", httpPort,
		"--verbosity", "6",
	)
	node.Cmd.InputLine(testPassphrase)
	defer func() {
		if t.Failed() {
			node.Shutdown()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

//确保所有端口都有活动的侦听器
//这样下一个节点就不会得到相同的
//调用GetAvailableTCPPorts时
	err = waitTCPPorts(ctx, ports...)
	if err != nil {
		t.Fatal(err)
	}

//等待节点启动
	for start := time.Now(); time.Since(start) < 10*time.Second; time.Sleep(50 * time.Millisecond) {
		node.Client, err = rpc.Dial(conf.IPCEndpoint())
		if err == nil {
			break
		}
	}
	if node.Client == nil {
		t.Fatal(err)
	}

//加载信息
	var info swarm.Info
	if err := node.Client.Call(&info, "bzz_info"); err != nil {
		t.Fatal(err)
	}
	node.Addr = net.JoinHostPort("127.0.0.1", info.Port)
node.URL = "http://“+No.ADDR”

	var nodeInfo p2p.NodeInfo
	if err := node.Client.Call(&nodeInfo, "admin_nodeInfo"); err != nil {
		t.Fatal(err)
	}
node.Enode = fmt.Sprintf("enode://%s@127.0.0.1:%s”，nodeinfo.id，p2pport）
	node.IpcPath = conf.IPCPath

	return node
}

func (n *testNode) Shutdown() {
	if n.Cmd != nil {
		n.Cmd.Kill()
	}
}

//
//
//
//
//并保证同一港口将被运回
//不同的调用，因为侦听器在函数内关闭，
//
//可用端口。
func getAvailableTCPPorts(count int) (ports []string, err error) {
	for i := 0; i < count; i++ {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, err
		}
//在循环中延迟关闭以确保同一端口不会
//在下一个迭代中被选择
		defer l.Close()

		_, port, err := net.SplitHostPort(l.Addr().String())
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}
	return ports, nil
}

//
//
//
//遇到错误。
//另请参见waitcpport。
func waitTCPPorts(ctx context.Context, ports ...string) error {
	var err error
//在中分配的mu locks err变量
//其他Goroutines
	var mu sync.Mutex

//
//
//防止不必要的等待
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()

			e := waitTCPPort(ctx, port)

			mu.Lock()
			defer mu.Unlock()
			if e != nil && err == nil {
				err = e
				cancel()
			}
		}(port)
	}
	wg.Wait()

	return err
}

//
//ONA提供的端口。它最多有3分钟的超时时间，
//
//提供的上下文实例。拨号程序超时10秒
//在每次迭代中，连接被拒绝的错误将
//在100毫秒内重试。
func waitTCPPort(ctx context.Context, port string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	for {
		c, err := (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, "tcp", "127.0.0.1:"+port)
		if err != nil {
			if operr, ok := err.(*net.OpError); ok {
				if syserr, ok := operr.Err.(*os.SyscallError); ok && syserr.Err == syscall.ECONNREFUSED {
					time.Sleep(100 * time.Millisecond)
					continue
				}
			}
			return err
		}
		return c.Close()
	}
}

