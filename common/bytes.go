
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:32</date>
//</624342607794802688>


//package common包含各种助手函数。
package common

import "encoding/hex"

//tohex返回b的十六进制表示形式，前缀为“0x”。
//对于空切片，返回值为“0x0”。
//
//已弃用：请改用hexutil.encode。
func ToHex(b []byte) string {
	hex := Bytes2Hex(b)
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}

//FromHex返回十六进制字符串s表示的字节。
//s的前缀可以是“0x”。
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

//CopyBytes返回所提供字节的精确副本。
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

//hashexprefix验证str以“0x”或“0x”开头。
func hasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

//IShexCharacter返回C的bool作为有效的十六进制。
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

//
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

//bytes2hex返回d的十六进制编码。
func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

//hex2bytes返回十六进制字符串str表示的字节。
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

//hex2bytesfixed返回指定固定长度flen的字节。
func Hex2BytesFixed(str string, flen int) []byte {
	h, _ := hex.DecodeString(str)
	if len(h) == flen {
		return h
	}
	if len(h) > flen {
		return h[len(h)-flen:]
	}
	hh := make([]byte, flen)
	copy(hh[flen-len(h):flen], h[:])
	return hh
}

//
func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded, slice)

	return padded
}

//LeftPadBytes将零个焊盘向左切片至长度L。
func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}

