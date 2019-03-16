
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:43</date>
//</624342657077874688>

//由“stringer-type=nodeEvent”生成的代码；不要编辑。

package discv5

import "strconv"

const _nodeEvent_name = "pongTimeoutpingTimeoutneighboursTimeout"

var _nodeEvent_index = [...]uint8{0, 11, 22, 39}

func (i nodeEvent) String() string {
	i -= 264
	if i >= nodeEvent(len(_nodeEvent_index)-1) {
		return "nodeEvent(" + strconv.FormatInt(int64(i+264), 10) + ")"
	}
	return _nodeEvent_name[_nodeEvent_index[i]:_nodeEvent_index[i+1]]
}

