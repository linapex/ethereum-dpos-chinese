
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:31</date>
//</624342605513101312>


package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/swarm"
)

//testcliswarmexportimport执行以下测试：
//
//
//三。运行本地数据存储的导出
//
//5。导入导出的数据存储
//6。从第二个节点获取上载的随机文件
func TestCLISwarmExportImport(t *testing.T) {
	cluster := newTestCluster(t, 1)

//
	f, cleanup := generateRandomFile(t, 10000000)
	defer cleanup()

//用“swarm up”上传文件，并期望得到一个哈希值
	up := runSwarm(t, "--bzzapi", cluster.Nodes[0].URL, "up", f.Name())
	_, matches := up.ExpectRegexp(`[a-f\d]{64}`)
	up.ExpectExit()
	hash := matches[0]

	var info swarm.Info
	if err := cluster.Nodes[0].Client.Call(&info, "bzz_info"); err != nil {
		t.Fatal(err)
	}

	cluster.Stop()
	defer cluster.Cleanup()

//生成export.tar
	exportCmd := runSwarm(t, "db", "export", info.Path+"/chunks", info.Path+"/export.tar", strings.TrimPrefix(info.BzzKey, "0x"))
	exportCmd.ExpectExit()

//
	cluster2 := newTestCluster(t, 1)

	var info2 swarm.Info
	if err := cluster2.Nodes[0].Client.Call(&info2, "bzz_info"); err != nil {
		t.Fatal(err)
	}

//
	cluster2.Stop()
	defer cluster2.Cleanup()

//导入export.tar
	importCmd := runSwarm(t, "db", "import", info2.Path+"/chunks", info.Path+"/export.tar", strings.TrimPrefix(info2.BzzKey, "0x"))
	importCmd.ExpectExit()

//旋转第二个群集备份
	cluster2.StartExistingNodes(t, 1, strings.TrimPrefix(info2.BzzAccount, "0x"))

//尝试获取导入的文件
	res, err := http.Get(cluster2.Nodes[0].URL + "/bzz:/" + hash)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected HTTP status %d, got %s", 200, res.Status)
	}

//
	mustEqualFiles(t, f, res.Body)
}

func mustEqualFiles(t *testing.T, up io.Reader, down io.Reader) {
	h := md5.New()
	upLen, err := io.Copy(h, up)
	if err != nil {
		t.Fatal(err)
	}
	upHash := h.Sum(nil)
	h.Reset()
	downLen, err := io.Copy(h, down)
	if err != nil {
		t.Fatal(err)
	}
	downHash := h.Sum(nil)

	if !bytes.Equal(upHash, downHash) || upLen != downLen {
		t.Fatalf("downloaded imported file md5=%x (length %v) is not the same as the generated one mp5=%x (length %v)", downHash, downLen, upHash, upLen)
	}
}

func generateRandomFile(t *testing.T, size int) (f *os.File, teardown func()) {
//创建tmp文件
	tmp, err := ioutil.TempFile("", "swarm-test")
	if err != nil {
		t.Fatal(err)
	}

//
	teardown = func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}

//
	buf := make([]byte, 10000000)
	_, err = rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile(tmp.Name(), buf, 0755)

	return tmp, teardown
}

