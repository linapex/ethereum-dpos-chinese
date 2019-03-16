
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:40</date>
//</624342640984330240>


//+建设！GO1.5

//go<1.5的跟踪方法没有OP实现。

package debug

import "errors"

func (*HandlerT) StartGoTrace(string) error {
	return errors.New("tracing is not supported on Go < 1.5")
}

func (*HandlerT) StopGoTrace() error {
	return errors.New("tracing is not supported on Go < 1.5")
}

