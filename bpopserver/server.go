package bpopserver

import (
	"github.com/realbmail/go-bas-mail-server/bmtpserver"
	"github.com/realbmail/go-bmail-protocol/translayer"
	"sync"
)

var (
	bpopserverInst     bmtpserver.BMTPServerIntf
	bpopserverInstLock sync.Mutex
)

func GetBMTPServer() bmtpserver.BMTPServerIntf {
	if bpopserverInst == nil {
		bpopserverInstLock.Lock()
		bpopserverInstLock.Unlock()
		if bpopserverInst == nil {
			bpopserverInst = bmtpserver.NewServer2(translayer.BPOP3)
		}
	}

	return bpopserverInst
}
