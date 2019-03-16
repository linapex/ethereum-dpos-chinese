
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:27</date>
//</624342586135416832>


//+构建iOS Linux，ARM64 Windows！达尔文！FreeBSD！Linux！NETBSD！索拉里斯

//这是目录监视的后备实现。
//它用于不受支持的平台。

package keystore

type watcher struct{ running bool }

func newWatcher(*accountCache) *watcher { return new(watcher) }
func (*watcher) start()                 {}
func (*watcher) close()                 {}

