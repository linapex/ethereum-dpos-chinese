
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:49</date>
//</624342679110553600>

package sctx

import "context"

type ContextKey int

const (
	HTTPRequestIDKey ContextKey = iota
	requestHostKey
)

func SetHost(ctx context.Context, domain string) context.Context {
	return context.WithValue(ctx, requestHostKey, domain)
}

func GetHost(ctx context.Context) string {
	v, ok := ctx.Value(requestHostKey).(string)
	if ok {
		return v
	}
	return ""
}

