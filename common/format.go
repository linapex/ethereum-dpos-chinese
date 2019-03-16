
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:32</date>
//</624342608939847680>


package common

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

//prettyDuration是时间的一个漂亮打印版本。Duration值会减少
//格式文本表示中不必要的精度。
type PrettyDuration time.Duration

var prettyDurationRe = regexp.MustCompile(`\.[0-9]+`)

//string实现了stringer接口，允许漂亮地打印持续时间
//数值四舍五入为三位小数。
func (d PrettyDuration) String() string {
	label := fmt.Sprintf("%v", time.Duration(d))
	if match := prettyDurationRe.FindString(label); len(match) > 4 {
		label = strings.Replace(label, match, match[:4], 1)
	}
	return label
}

