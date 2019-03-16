
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:39</date>
//</624342639830896640>


package filter

import (
	"testing"
	"time"
)

//简单测试，检查基线匹配/不匹配过滤是否有效。
func TestFilters(t *testing.T) {
	fm := New()
	fm.Start()

//注册两个筛选器以捕获已发布的数据
	first := make(chan struct{})
	fm.Install(Generic{
		Str1: "hello",
		Fn: func(data interface{}) {
			first <- struct{}{}
		},
	})
	second := make(chan struct{})
	fm.Install(Generic{
		Str1: "hello1",
		Str2: "hello",
		Fn: func(data interface{}) {
			second <- struct{}{}
		},
	})
//发布只应与第一个筛选器匹配的事件
	fm.Notify(Generic{Str1: "hello"}, true)
	fm.Stop()

//确保只有Mathing过滤器启动
	select {
	case <-first:
	case <-time.After(100 * time.Millisecond):
		t.Error("matching filter timed out")
	}
	select {
	case <-second:
		t.Error("mismatching filter fired")
	case <-time.After(100 * time.Millisecond):
	}
}

