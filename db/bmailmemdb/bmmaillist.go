package bmailmemdb

import (
	"encoding/json"
	"github.com/realbmail/go-bas-mail-server/config"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type BMBlockMailList struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

var (
	bmbmlStore     *BMBlockMailList
	bmbmlStoreLock sync.Mutex
)

type BlockMailAddress struct {
	CreateTime int64 `json:"ct"`
	UpdateTime int64 `json:"ut"`
}

func newBMBlockMailList() *BMBlockMailList {
	cfg := config.GetBMSCfg()

	db := db.NewFileDb(cfg.GetBMMLSavePath()).Load()

	return &BMBlockMailList{NbsDbInter: db}
}

func GetBMBlockMailList() *BMBlockMailList {
	if bmbmlStore == nil {
		bmbmlStoreLock.Lock()
		defer bmbmlStoreLock.Unlock()

		if bmbmlStore == nil {
			bmbmlStore = newBMBlockMailList()
		}

	}

	return bmbmlStore
}

func (s *BMBlockMailList) Insert(mAddr string) error {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if _, err := s.NbsDbInter.Find(mAddr); err == nil {
		return err
	}
	now := tools.GetNowMsTime()

	bma := &BlockMailAddress{now, now}

	if v, err := json.Marshal(*bma); err != nil {
		return err
	} else {
		return s.NbsDbInter.Insert(mAddr, string(v))
	}

}

func (s *BMBlockMailList) Find(mAddr string) (*BlockMailAddress, error) {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if vs, err := s.NbsDbInter.Find(mAddr); err != nil {
		return nil, err
	} else {
		f := &BlockMailAddress{}
		err = json.Unmarshal([]byte(vs), f)
		if err != nil {
			return nil, err
		} else {
			return f, err
		}
	}
}

func (s *BMBlockMailList) Remove(mAddr string) {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.NbsDbInter.Delete(mAddr)

}

func (s *BMBlockMailList) Save() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.NbsDbInter.Save()
}

func (s *BMBlockMailList) Iterator() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.cursor = s.NbsDbInter.DBIterator()
}

func (s *BMBlockMailList) Next() (key string, meta *BlockMailAddress, r1 error) {
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
	meta = &BlockMailAddress{}

	if err := json.Unmarshal([]byte(v), meta); err != nil {
		return "", nil, err
	}

	key = k

	return

}
