
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:37</date>
//</624342632125960192>


package dashboard

import "time"

//默认配置包含仪表板的默认设置。
var DefaultConfig = Config{
	Host:    "localhost",
	Port:    8080,
	Refresh: 5 * time.Second,
}

//配置包含仪表板的配置参数。
type Config struct {
//主机是启动仪表板服务器的主机接口。如果这样
//字段为空，将不启动任何仪表板。
	Host string `toml:",omitempty"`

//端口是启动仪表板服务器的TCP端口号。这个
//默认的零值是/有效的，将随机选择端口号（有用
//对于季节性节点）。
	Port int `toml:",omitempty"`

//refresh是数据更新的刷新率，通常会收集图表条目。
	Refresh time.Duration `toml:",omitempty"`
}

