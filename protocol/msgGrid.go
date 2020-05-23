package protocol

import (
	"github.com/BASChain/go-bmail-protocol/translayer"
)

type RBody interface {
	UnPack(data []byte) error
	Verify() bool
	SetCurrentSn(sn []byte)
	Dispatch() error
	Response() (WBody, error)
}

type WBody interface {
	MsgType() uint16
	GetBytes() ([]byte, error)
}

var MsgGrid = map[uint16]RBody{
	translayer.SEND_CRYPT_ENVELOPE: &CryptEnvelopeMsg{},
	translayer.RETR:                &CommandDownloadMsg{},
}
