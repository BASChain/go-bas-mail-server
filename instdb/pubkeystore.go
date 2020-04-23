package instdb

import (
	"github.com/kprc/ssactiveserver/config"
	"github.com/kprc/ssactiveserver/kvdb"
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

	pubkeystore = kvdb.NewKVDB("pubkeystore", config.GetSSSCfg().GetPKPath())

	return pubkeystore
}
