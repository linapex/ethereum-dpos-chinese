
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:35</date>
//</624342620939751424>


package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//目的地为每个契约存储一个映射（由代码散列键控）。
//地图包含JumpDest的每个位置的条目。
//指令。
type destinations map[common.Hash]bitvec

//检查代码在dest是否有JumpDest。
func (d destinations) has(codehash common.Hash, code []byte, dest *big.Int) bool {
//PC不能超过len（code），当然不能超过63位。
//在这种情况下，不要费心检查Jumpdest。
	udest := dest.Uint64()
	if dest.BitLen() >= 63 || udest >= uint64(len(code)) {
		return false
	}

	m, analysed := d[codehash]
	if !analysed {
		m = codeBitmap(code)
		d[codehash] = m
	}
	return OpCode(code[udest]) == JUMPDEST && m.codeSegment(udest)
}

//bitvec是一个位向量，它映射程序中的字节。
//未设置位表示字节为操作码，设置位表示
//这是数据（即pushxx的参数）。
type bitvec []byte

func (bits *bitvec) set(pos uint64) {
	(*bits)[pos/8] |= 0x80 >> (pos % 8)
}
func (bits *bitvec) set8(pos uint64) {
	(*bits)[pos/8] |= 0xFF >> (pos % 8)
	(*bits)[pos/8+1] |= ^(0xFF >> (pos % 8))
}

//代码段检查位置是否在代码段中。
func (bits *bitvec) codeSegment(pos uint64) bool {
	return ((*bits)[pos/8] & (0x80 >> (pos % 8))) == 0
}

//codeBitmap以代码的形式收集数据位置。
func codeBitmap(code []byte) bitvec {
//如果代码
//以push32结束，算法将把0推到
//位向量超出了实际代码的边界。
	bits := make(bitvec, len(code)/8+1+4)
	for pc := uint64(0); pc < uint64(len(code)); {
		op := OpCode(code[pc])

		if op >= PUSH1 && op <= PUSH32 {
			numbits := op - PUSH1 + 1
			pc++
			for ; numbits >= 8; numbits -= 8 {
bits.set8(pc) //八
				pc += 8
			}
			for ; numbits > 0; numbits-- {
				bits.set(pc)
				pc++
			}
		} else {
			pc++
		}
	}
	return bits
}

