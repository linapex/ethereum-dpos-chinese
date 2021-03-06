
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:26</date>
//</624342584537387008>


package accounts

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

//DefaultRootDerivationPath是自定义派生终结点的根路径
//被追加。因此，第一个帐户将位于m/44'/60'/0'/0，第二个帐户将位于
//在m/44'/60'/0'/1等处。
var DefaultRootDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}

//DefaultBaseDerivationPath是自定义派生终结点的基本路径
//是递增的。因此，第一个帐户将位于m/44'/60'/0'/0，第二个帐户将位于
//在m/44'/60'/0'/1等处。
var DefaultBaseDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}

//DefaultLedgerBasederivationPath是自定义派生终结点的基本路径
//是递增的。因此，第一个帐户将位于m/44'/60'/0'/0，第二个帐户将位于
//在m/44'/60'/0'/1等处。
var DefaultLedgerBaseDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}

//派生路径表示层次结构的计算机友好版本
//确定的钱包帐户派生路径。
//
//BIP-32规范https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
//定义派生路径的形式：
//
//M/用途/硬币类型/账户/更改/地址索引
//
//BIP-44规范https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
//定义加密货币的“用途”为44'（或0x800002c），以及
//slip-44 https://github.com/satoshilabs/slip/blob/master/slip-0044.md分配
//以太坊的“硬币”类型为“60”（或0x8000003c）。
//
//根据规范，以太坊的根路径为m/44'/60'/0'/0
//来自https://github.com/ethereum/eips/issues/84，尽管它不是石头做的
//然而，帐户应该增加最后一个组件还是
//那。我们将使用更简单的方法来增加最后一个组件。
type DerivationPath []uint32

//ParseDerivationPath将用户指定的派生路径字符串转换为
//内部二进制表示。
//
//完整的派生路径需要以“m/”前缀开头，相对派生
//路径（将附加到默认根路径）不能有前缀
//在第一个元素前面。空白被忽略。
func ParseDerivationPath(path string) (DerivationPath, error) {
	var result DerivationPath

//处理绝对或相对路径
	components := strings.Split(path, "/")
	switch {
	case len(components) == 0:
		return nil, errors.New("empty derivation path")

	case strings.TrimSpace(components[0]) == "":
		return nil, errors.New("ambiguous path: use 'm/' prefix for absolute paths, or no leading '/' for relative ones")

	case strings.TrimSpace(components[0]) == "m":
		components = components[1:]

	default:
		result = append(result, DefaultRootDerivationPath...)
	}
//其余所有组件都是相对的，逐个附加
	if len(components) == 0 {
return nil, errors.New("empty derivation path") //空的相对路径
	}
	for _, component := range components {
//忽略任何用户添加的空白
		component = strings.TrimSpace(component)
		var value uint32

//处理硬化路径
		if strings.HasSuffix(component, "'") {
			value = 0x80000000
			component = strings.TrimSpace(strings.TrimSuffix(component, "'"))
		}
//处理非硬化部件
		bigval, ok := new(big.Int).SetString(component, 0)
		if !ok {
			return nil, fmt.Errorf("invalid component: %s", component)
		}
		max := math.MaxUint32 - value
		if bigval.Sign() < 0 || bigval.Cmp(big.NewInt(int64(max))) > 0 {
			if value == 0 {
				return nil, fmt.Errorf("component %v out of allowed range [0, %d]", bigval, max)
			}
			return nil, fmt.Errorf("component %v out of allowed hardened range [0, %d]", bigval, max)
		}
		value += uint32(bigval.Uint64())

//追加并重复
		result = append(result, value)
	}
	return result, nil
}

//字符串实现Stringer接口，转换二进制派生路径
//它的规范表示。
func (path DerivationPath) String() string {
	result := "m"
	for _, component := range path {
		var hardened bool
		if component >= 0x80000000 {
			component -= 0x80000000
			hardened = true
		}
		result = fmt.Sprintf("%s/%d", result, component)
		if hardened {
			result += "'"
		}
	}
	return result
}

