
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:41</date>
//</624342647116402688>

package term

import (
	"syscall"
)

const ioctlReadTermios = syscall.TIOCGETA

//Go 1.2不包括FreeBSD的Termios。这应该在1.3中添加，并且可以与达尔文航站楼合并。
type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed uint32
	Ospeed uint32
}

