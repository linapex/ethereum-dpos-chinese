
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:41</date>
//</624342647074459648>

//基于ssh/终端：
//版权所有2013 Go作者。版权所有。
//此源代码的使用受BSD样式的控制
//可以在许可文件中找到的许可证。
//+建设！应用程序引擎

package term

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

type Termios syscall.Termios

