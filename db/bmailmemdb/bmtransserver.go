package bmailmemdb

import (
	"github.com/kprc/nbsnetwork/db"
	"sync"
)

type BMBlockTransServer struct {
	blockTransferServer db.NbsDbInter
	blockTransSrvLock sync.Mutex
	btsCursor *db.DBCusor
}
