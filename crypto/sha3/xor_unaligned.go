
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:37</date>
//</624342628749545472>

//版权所有2015 The Go作者。版权所有。
//此源代码的使用受BSD样式的控制
//可以在许可文件中找到的许可证。

//+构建AMD64 386 PPC64LE
//+建设！应用程序引擎

package sha3

import "unsafe"

func xorInUnaligned(d *state, buf []byte) {
	bw := (*[maxRate / 8]uint64)(unsafe.Pointer(&buf[0]))
	n := len(buf)
	if n >= 72 {
		d.a[0] ^= bw[0]
		d.a[1] ^= bw[1]
		d.a[2] ^= bw[2]
		d.a[3] ^= bw[3]
		d.a[4] ^= bw[4]
		d.a[5] ^= bw[5]
		d.a[6] ^= bw[6]
		d.a[7] ^= bw[7]
		d.a[8] ^= bw[8]
	}
	if n >= 104 {
		d.a[9] ^= bw[9]
		d.a[10] ^= bw[10]
		d.a[11] ^= bw[11]
		d.a[12] ^= bw[12]
	}
	if n >= 136 {
		d.a[13] ^= bw[13]
		d.a[14] ^= bw[14]
		d.a[15] ^= bw[15]
		d.a[16] ^= bw[16]
	}
	if n >= 144 {
		d.a[17] ^= bw[17]
	}
	if n >= 168 {
		d.a[18] ^= bw[18]
		d.a[19] ^= bw[19]
		d.a[20] ^= bw[20]
	}
}

func copyOutUnaligned(d *state, buf []byte) {
	ab := (*[maxRate]uint8)(unsafe.Pointer(&d.a[0]))
	copy(buf, ab[:])
}

var (
	xorIn   = xorInUnaligned
	copyOut = copyOutUnaligned
)

const xorImplementationUnaligned = "unaligned"

