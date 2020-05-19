package wallet

import "github.com/BASChain/go-bmail-account"

type ServerWallet interface {
	BCAddress() bmail.Address
}
