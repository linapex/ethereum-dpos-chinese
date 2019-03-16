
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:47</date>
//</624342671875379200>

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

package multihash

import (
	"bytes"
	"math/rand"
	"testing"
)

//
func TestCheckMultihash(t *testing.T) {
	hashbytes := make([]byte, 32)
	c, err := rand.Read(hashbytes)
	if err != nil {
		t.Fatal(err)
	} else if c < 32 {
		t.Fatal("short read")
	}

	expected := ToMultihash(hashbytes)

	l, hl, _ := GetMultihashLength(expected)
	if l != 32 {
		t.Fatalf("expected length %d, got %d", 32, l)
	} else if hl != 2 {
		t.Fatalf("expected header length %d, got %d", 2, hl)
	}
	if _, _, err := GetMultihashLength(expected[1:]); err == nil {
		t.Fatal("expected failure on corrupt header")
	}
	if _, _, err := GetMultihashLength(expected[:len(expected)-2]); err == nil {
		t.Fatal("expected failure on short content")
	}
	dh, _ := FromMultihash(expected)
	if !bytes.Equal(dh, hashbytes) {
		t.Fatalf("expected content hash %x, got %x", hashbytes, dh)
	}
}

