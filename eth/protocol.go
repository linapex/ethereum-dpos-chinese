
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:39</date>
//</624342636404150272>


package eth

import (
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rlp"
)
//用于匹配协议版本和消息的常量
const (
	eth62 = 62
	eth63 = 63
)

//ProtocolName是在能力协商期间使用的协议的官方简称。
var ProtocolName = "eth"

//协议版本是ETH协议的支持版本（首先是主要版本）。
var ProtocolVersions = []uint{eth63, eth62}

//Protocollength是对应于不同协议版本的已实现消息数。
var ProtocolLengths = []uint64{17, 8}

const ProtocolMaxMsgSize = 10 * 1024 * 1024 //协议消息大小的最大上限

//ETH协议报文代码
const (
//属于ETH/62的协议消息
	StatusMsg          = 0x00
	NewBlockHashesMsg  = 0x01
	TxMsg              = 0x02
	GetBlockHeadersMsg = 0x03
	BlockHeadersMsg    = 0x04
	GetBlockBodiesMsg  = 0x05
	BlockBodiesMsg     = 0x06
	NewBlockMsg        = 0x07

//属于ETH/63的协议消息
	GetNodeDataMsg = 0x0d
	NodeDataMsg    = 0x0e
	GetReceiptsMsg = 0x0f
	ReceiptsMsg    = 0x10
)

type errCode int

const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrProtocolVersionMismatch
	ErrNetworkIdMismatch
	ErrGenesisBlockMismatch
	ErrNoStatusMsg
	ErrExtraStatusMsg
	ErrSuspendedPeer
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

//一旦旧代码用完，XXX就会更改
var errorToString = map[int]string{
	ErrMsgTooLarge:             "Message too long",
	ErrDecode:                  "Invalid message",
	ErrInvalidMsgCode:          "Invalid message code",
	ErrProtocolVersionMismatch: "Protocol version mismatch",
	ErrNetworkIdMismatch:       "NetworkId mismatch",
	ErrGenesisBlockMismatch:    "Genesis block mismatch",
	ErrNoStatusMsg:             "No status message",
	ErrExtraStatusMsg:          "Extra status message",
	ErrSuspendedPeer:           "Suspended peer",
}

type txPool interface {
//AddRemotes应该将给定的事务添加到池中。
	AddRemotes([]*types.Transaction) []error

//挂起应返回挂起的事务。
//该切片应由调用方可修改。
	Pending() (map[common.Address]types.Transactions, error)

//subscribenewtxsevent应返回的事件订阅
//newtxSevent并将事件发送到给定的通道。
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
}

//statusdata是状态消息的网络包。
type statusData struct {
	ProtocolVersion uint32
	NetworkId       uint64
	TD              *big.Int
	CurrentBlock    common.Hash
	GenesisBlock    common.Hash
}

//newblockhashesdata是块通知的网络包。
type newBlockHashesData []struct {
Hash   common.Hash //正在公布的一个特定块的哈希
Number uint64      //公布的一个特定区块的编号
}

//GetBlockHeadersData表示块头查询。
type getBlockHeadersData struct {
Origin  hashOrNumber //从中检索邮件头的块
Amount  uint64       //要检索的最大头数
Skip    uint64       //要在连续标题之间跳过的块
Reverse bool         //查询方向（假=上升到最新，真=下降到创世纪）
}

//hashornumber是用于指定源块的组合字段。
type hashOrNumber struct {
Hash   common.Hash //要从中检索头的块哈希（不包括数字）
Number uint64      //要从中检索头的块哈希（不包括哈希）
}

//encoderlp是一个专门的编码器，用于hashornumber只对
//两个包含联合字段。
func (hn *hashOrNumber) EncodeRLP(w io.Writer) error {
	if hn.Hash == (common.Hash{}) {
		return rlp.Encode(w, hn.Number)
	}
	if hn.Number != 0 {
		return fmt.Errorf("both origin hash (%x) and number (%d) provided", hn.Hash, hn.Number)
	}
	return rlp.Encode(w, hn.Hash)
}

//decoderlp是一种特殊的译码器，用于hashornumber对内容进行译码。
//分块散列或分块编号。
func (hn *hashOrNumber) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	origin, err := s.Raw()
	if err == nil {
		switch {
		case size == 32:
			err = rlp.DecodeBytes(origin, &hn.Hash)
		case size <= 8:
			err = rlp.DecodeBytes(origin, &hn.Number)
		default:
			err = fmt.Errorf("invalid input size %d for origin", size)
		}
	}
	return err
}

//newblockdata是块传播消息的网络包。
type newBlockData struct {
	Block *types.Block
	TD    *big.Int
}

//BlockBody表示单个块的数据内容。
type blockBody struct {
Transactions []*types.Transaction //块中包含的事务
Uncles       []*types.Header      //一个街区内的叔叔
}

//blockbodiesdata是用于块内容分发的网络包。
type blockBodiesData []*blockBody

