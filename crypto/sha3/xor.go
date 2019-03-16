
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:37</date>
//</624342628611133440>

//版权所有2015 The Go作者。版权所有。
//此源代码的使用受BSD样式的控制
//可以在许可文件中找到的许可证。

//+建设！AMD64！386！PPC64LE发动机

package sha3

var (
	xorIn            = xorInGeneric
	copyOut          = copyOutGeneric
	xorInUnaligned   = xorInGeneric
	copyOutUnaligned = copyOutGeneric
)

const xorImplementationUnaligned = "generic"

