
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:49</date>
//</624342680939270144>

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

package storage

import (
	"errors"
)

const (
	ErrInit = iota
	ErrNotFound
	ErrIO
	ErrUnauthorized
	ErrInvalidValue
	ErrDataOverflow
	ErrNothingToReturn
	ErrCorruptData
	ErrInvalidSignature
	ErrNotSynced
	ErrPeriodDepth
	ErrCnt
)

var (
	ErrChunkNotFound    = errors.New("chunk not found")
	ErrFetching         = errors.New("chunk still fetching")
	ErrChunkInvalid     = errors.New("invalid chunk")
	ErrChunkForward     = errors.New("cannot forward")
	ErrChunkUnavailable = errors.New("chunk unavailable")
	ErrChunkTimeout     = errors.New("timeout")
)

