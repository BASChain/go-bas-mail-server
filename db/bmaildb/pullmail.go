package bmaildb

import (
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/kprc/nbsnetwork/db"
	"sync"
)

var (
	pullMailDbStore     *BMSaveMailDb
	pullMailDbStoreLock sync.Mutex
)

func newBMPullMailDb() *BMSaveMailDb {
	cfg := config.GetBMSCfg()

	db := db.NewFileDb(cfg.GetPullMailDbPath()).Load()

	return &BMSaveMailDb{NbsDbInter: db}

}

func GetBMPullMailDb() *BMSaveMailDb {
	if pullMailDbStore == nil {
		pullMailDbStoreLock.Lock()
		defer pullMailDbStoreLock.Unlock()

		if pullMailDbStore == nil {
			pullMailDbStore = newBMPullMailDb()
		}

	}
	return pullMailDbStore
}
