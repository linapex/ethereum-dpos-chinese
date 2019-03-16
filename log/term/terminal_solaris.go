
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:41</date>
//</624342647389032448>

package term

import "golang.org/x/sys/unix"

//如果给定的文件描述符是终端，则istty返回true。
func IsTty(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETA)
	return err == nil
}

