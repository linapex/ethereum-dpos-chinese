
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:44</date>
//</624342660945022976>

package adapters

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"
)

//wsaddrpattern是一个regex，用于从节点的
//日志
var wsAddrPattern = regexp.MustCompile(`ws://[\D::] +）

func matchWSAddr(str string) (string, bool) {
	if !strings.Contains(str, "WebSocket endpoint opened") {
		return "", false
	}

	return wsAddrPattern.FindString(str), true
}

//findwsaddr通过读卡器r扫描，查找
//WebSocket地址信息。
func findWSAddr(r io.Reader, timeout time.Duration) (string, error) {
	ch := make(chan string)

	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			addr, ok := matchWSAddr(s.Text())
			if ok {
				ch <- addr
			}
		}
		close(ch)
	}()

	var wsAddr string
	select {
	case wsAddr = <-ch:
		if wsAddr == "" {
			return "", errors.New("empty result")
		}
	case <-time.After(timeout):
		return "", errors.New("timed out")
	}

	return wsAddr, nil
}

