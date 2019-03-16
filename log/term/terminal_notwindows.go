
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:41</date>
//</624342647263203328>

//基于ssh/终端：
//版权所有2011 Go作者。版权所有。
//此源代码的使用受BSD样式的控制
//可以在许可文件中找到的许可证。

//+建立Linux，！appengine darwin freebsd openbsd netbsd

package term

import (
	"syscall"
	"unsafe"
)

//如果给定的文件描述符是终端，则istty返回true。
func IsTty(fd uintptr) bool {
	var termios Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

