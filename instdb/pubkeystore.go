package instdb

import (
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/BASChain/go-bas-mail-server/kvdb"
	"sync"
)

var (
	pubkeystore kvdb.KVDBInterface
	pkstorelock sync.Mutex
)

func GetPKDB() kvdb.KVDBInterface {
	if pubkeystore != nil {
		return pubkeystore
	}

	pkstorelock.Lock()
	defer pkstorelock.Unlock()

	if pubkeystore != nil {
		return pubkeystore
	}

	pubkeystore = kvdb.NewKVDB("pubkeystore", config.GetBMSCfg().GetPKPath())

	return pubkeystore
}
