package wallet

import (
	"github.com/BASChain/go-bmail-account"
	"crypto/rand"
	"sync"
)

type ServerWalletIntf interface {
	BCAddress() bmail.Address
}

type ServerWallet struct {
	Addr bmail.Address
}

var (
	serverWalletInst ServerWalletIntf
	serverWalletInstLock sync.Mutex
)

func NewAddr() []byte  {
	sn := make([]byte, 32)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	return sn
}

func (sw *ServerWallet)BCAddress() bmail.Address  {
	return sw.Addr
}

func NewWallet() ServerWalletIntf {
	sw:=&ServerWallet{}
	sw.Addr = bmail.ToAddress(NewAddr())

	return sw

}

func GetServerWallet() ServerWalletIntf{
	if serverWalletInst == nil{
		serverWalletInstLock.Lock()
		defer serverWalletInstLock.Unlock()
		if serverWalletInst == nil{
			serverWalletInst = NewWallet()
		}
	}

	return serverWalletInst
}



