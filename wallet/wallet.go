package wallet

import (
	"crypto/ed25519"
	"github.com/BASChain/go-bas-mail-server/bmailcrypt"
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/BASChain/go-bmail-account"
	"sync"
)

type ServerWalletIntf interface {
	BCAddress() bmail.Address
	Sign(message []byte) []byte
	//Verify(message,sig []byte) bool
}

type ServerWallet struct {
	Addr    bmail.Address
	PubKey  ed25519.PublicKey
	PrivKey ed25519.PrivateKey
}

var (
	serverWalletInst     ServerWalletIntf
	serverWalletInstLock sync.Mutex
)

func NewWallet() ServerWalletIntf {
	sw := &ServerWallet{}
	cfg := config.GetBMSCfg()
	sw.Addr = bmail.ToAddress(cfg.PubKey)
	sw.PubKey = cfg.PubKey
	sw.PrivKey = cfg.PrivKey

	return sw

}

func GetServerWallet() ServerWalletIntf {
	if serverWalletInst == nil {
		serverWalletInstLock.Lock()
		defer serverWalletInstLock.Unlock()
		if serverWalletInst == nil {
			serverWalletInst = NewWallet()
		}
	}

	return serverWalletInst
}

func (sw *ServerWallet) Sign(message []byte) []byte {
	return bmailcrypt.Sign(sw.PrivKey, message)
}

func (sw *ServerWallet) BCAddress() bmail.Address {
	return sw.Addr
}
