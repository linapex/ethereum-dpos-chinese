
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:45</date>
//</624342664774422528>


//+构建darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package rpc

import (
	"context"
	"net"
	"os"
	"path/filepath"
)

//ipclisten将在给定的端点上创建一个Unix套接字。
func ipcListen(endpoint string) (net.Listener, error) {
//确保存在IPC路径，并删除以前的任何剩余部分
	if err := os.MkdirAll(filepath.Dir(endpoint), 0751); err != nil {
		return nil, err
	}
	os.Remove(endpoint)
	l, err := net.Listen("unix", endpoint)
	if err != nil {
		return nil, err
	}
	os.Chmod(endpoint, 0600)
	return l, nil
}

//newipcconnection将连接到给定端点上的UNIX套接字。
func newIPCConnection(ctx context.Context, endpoint string) (net.Conn, error) {
	return dialContext(ctx, "unix", endpoint)
}

