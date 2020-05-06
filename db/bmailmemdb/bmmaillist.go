package bmailmemdb

import (
	"github.com/kprc/nbsnetwork/db"
	"sync"
)

type BMBlockMailList struct {
	mailBlackList   db.NbsDbInter
	mailBlLstLock sync.Mutex
	mblCursor *db.DBCusor
}