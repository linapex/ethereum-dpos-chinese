
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:42</date>
//</624342651017105408>

//+构建CGO
//+建设！应用程序引擎

package metrics

import "runtime"

func numCgoCall() int64 {
	return runtime.NumCgoCall()
}

