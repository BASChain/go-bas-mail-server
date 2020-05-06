package bmailmemdb

import (
	"github.com/kprc/nbsnetwork/db"
	"sync"
)

type BMBlockTransServer struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}
