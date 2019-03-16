
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:51</date>
//</624342690145767424>

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

/*











*/


//

package whisperv6

import (
	"time"
)

//
const (
ProtocolVersion    = uint64(6) //
ProtocolVersionStr = "6.0"     //
ProtocolName       = "shh"     //

//
statusCode           = 0   //
messagesCode         = 1   //
powRequirementCode   = 2   //
bloomFilterExCode    = 3   //
p2pRequestCode       = 126 //
p2pMessageCode       = 127 //
	NumberOfMessageCodes = 128

SizeMask      = byte(3) //
	signatureFlag = byte(4)

TopicLength     = 4  //
signatureLength = 65 //
aesKeyLength    = 32 //
aesNonceLength  = 12 //
keyIDSize       = 32 //
BloomFilterSize = 64 //
	flagsLength     = 1

	EnvelopeHeaderLength = 20

MaxMessageSize        = uint32(10 * 1024 * 1024) //
	DefaultMaxMessageSize = uint32(1024 * 1024)
	DefaultMinimumPoW     = 0.2

padSizeLimit      = 256 //
	messageQueueLimit = 1024

	expirationCycle   = time.Second
	transmissionCycle = 300 * time.Millisecond

DefaultTTL           = 50 //
DefaultSyncAllowance = 10 //
)

//
//
//
//
//
//
type MailServer interface {
	Archive(env *Envelope)
	DeliverMail(whisperPeer *Peer, request *Envelope)
}

