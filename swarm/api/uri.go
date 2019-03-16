
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342670319292416>

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

package api

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

//
//
var hashMatcher = regexp.MustCompile("^([0-9A-Fa-f]{64})([0-9A-Fa-f]{64})?$")

//
type URI struct {
//
//
//
//
//
//
//
//
	Scheme string

//
//
	Addr string

//
	addr storage.Address

//
	Path string
}

func (u *URI) MarshalJSON() (out []byte, err error) {
	return []byte(`"` + u.String() + `"`), nil
}

func (u *URI) UnmarshalJSON(value []byte) error {
	uri, err := Parse(string(value))
	if err != nil {
		return err
	}
	*u = *uri
	return nil
}

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
func Parse(rawuri string) (*URI, error) {
	u, err := url.Parse(rawuri)
	if err != nil {
		return nil, err
	}
	uri := &URI{Scheme: u.Scheme}

//
	switch uri.Scheme {
	case "bzz", "bzz-raw", "bzz-immutable", "bzz-list", "bzz-hash", "bzz-resource":
	default:
		return nil, fmt.Errorf("unknown scheme %q", u.Scheme)
	}

//
//
	if u.Host != "" {
		uri.Addr = u.Host
		uri.Path = strings.TrimLeft(u.Path, "/")
		return uri, nil
	}

//
//
	parts := strings.SplitN(strings.TrimLeft(u.Path, "/"), "/", 2)
	uri.Addr = parts[0]
	if len(parts) == 2 {
		uri.Path = parts[1]
	}
	return uri, nil
}
func (u *URI) Resource() bool {
	return u.Scheme == "bzz-resource"
}

func (u *URI) Raw() bool {
	return u.Scheme == "bzz-raw"
}

func (u *URI) Immutable() bool {
	return u.Scheme == "bzz-immutable"
}

func (u *URI) List() bool {
	return u.Scheme == "bzz-list"
}

func (u *URI) Hash() bool {
	return u.Scheme == "bzz-hash"
}

func (u *URI) String() string {
	return u.Scheme + ":/" + u.Addr + "/" + u.Path
}

func (u *URI) Address() storage.Address {
	if u.addr != nil {
		return u.addr
	}
	if hashMatcher.MatchString(u.Addr) {
		u.addr = common.Hex2Bytes(u.Addr)
		return u.addr
	}
	return nil
}

