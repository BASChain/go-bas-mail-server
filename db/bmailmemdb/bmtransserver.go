package bmailmemdb

import (
	"encoding/json"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/realbmail/go-bas-mail-server/config"
	"sync"
)

type BMBlockTransServerDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

var (
	bmbtsStore     *BMBlockTransServerDb
	bmbtsStoreLock sync.Mutex
)

type BlockTransServer struct {
	CreateTime int64 `json:"ct"`
	UpdateTime int64 `json:"ut"`
}

func newBmbtsStore() *BMBlockTransServerDb {
	cfg := config.GetBMSCfg()
	db := db.NewFileDb(cfg.GetBMTransferSavePath()).Load()

	return &BMBlockTransServerDb{NbsDbInter: db}
}

func GetBMBlockTransStore() *BMBlockTransServerDb {
	if bmbtsStore == nil {
		bmbtsStoreLock.Lock()
		defer bmbtsStoreLock.Unlock()
		if bmbtsStore == nil {
			bmbtsStore = newBmbtsStore()
		}
	}

	return bmbtsStore
}

func (s *BMBlockTransServerDb) Insert(srvDomain string) error {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if _, err := s.NbsDbInter.Find(srvDomain); err == nil {
		return err
	}
	now := tools.GetNowMsTime()

	bts := &BlockTransServer{now, now}

	if v, err := json.Marshal(*bts); err != nil {
		return err
	} else {
		return s.NbsDbInter.Insert(srvDomain, string(v))
	}

}

func (s *BMBlockTransServerDb) Find(srvDomain string) (*BlockTransServer, error) {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if vs, err := s.NbsDbInter.Find(srvDomain); err != nil {
		return nil, err
	} else {
		f := &BlockTransServer{}
		err = json.Unmarshal([]byte(vs), f)
		if err != nil {
			return nil, err
		} else {
			return f, err
		}
	}
}

func (s *BMBlockTransServerDb) Remove(srvDomain string) {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.NbsDbInter.Delete(srvDomain)

}

func (s *BMBlockTransServerDb) Save() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.NbsDbInter.Save()
}

func (s *BMBlockTransServerDb) Iterator() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.cursor = s.NbsDbInter.DBIterator()
}

func (s *BMBlockTransServerDb) Next() (key string, meta *BlockTransServer, r1 error) {
	if s.cursor == nil {
		return
	}
	s.dbLock.Lock()
	s.dbLock.Unlock()
	k, v := s.cursor.Next()
	if k == "" {
		s.dbLock.Unlock()
		return "", nil, nil
	}
	s.dbLock.Unlock()
	meta = &BlockTransServer{}

	if err := json.Unmarshal([]byte(v), meta); err != nil {
		return "", nil, err
	}

	key = k

	return

}
