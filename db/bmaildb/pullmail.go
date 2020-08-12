package bmaildb

import (
	"github.com/kprc/nbsnetwork/db"
	"github.com/realbmail/go-bas-mail-server/config"
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
