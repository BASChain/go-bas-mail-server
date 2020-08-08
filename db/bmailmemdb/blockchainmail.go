package bmailmemdb

import (
	"encoding/base64"
	"encoding/json"
	"github.com/realbmail/go-bas-mail-server/config"
	"github.com/realbmail/go-bas-mail-server/db/dbcommon"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type BMailMetaDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

var (
	bmmStore     *BMailMetaDb
	bmmStoreLock sync.Mutex
)

type BMailMeta struct {
	BMailAddress string        `json:"bma"`
	PublicKey    dbcommon.PKey `json:"pk"`
	CreateTime   int64         `json:"ct"`
	UpdateTime   int64         `json:"ut"`
}

func newBMMDb() *BMailMetaDb {
	cfg := config.GetBMSCfg()

	db := db.NewFileDb(cfg.GetBMMSavePath()).Load()

	return &BMailMetaDb{NbsDbInter: db}
}

func GetBMMStore() *BMailMetaDb {
	if bmmStore == nil {
		bmmStoreLock.Lock()
		defer bmmStoreLock.Unlock()
		if bmmStore == nil {
			bmmStore = newBMMDb()
		}
	}

	return bmmStore
}

func (s *BMailMetaDb) Insert(mHash dbcommon.Hash, mAddress string, pk dbcommon.PKey) error {

	k := base64.StdEncoding.EncodeToString(mHash[:])

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if _, err := s.NbsDbInter.Find(k); err == nil {
		return err
	}

	now := tools.GetNowMsTime()

	bmm := &BMailMeta{mAddress, pk, now, now}

	if v, err := json.Marshal(*bmm); err != nil {
		return err
	} else {
		return s.NbsDbInter.Insert(k, string(v))
	}
}

func (s *BMailMetaDb) Update(mHash dbcommon.Hash, mAddress string, pk dbcommon.PKey) (old *BMailMeta, err error) {
	k := base64.StdEncoding.EncodeToString(mHash[:])

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	var vs string
	if vs, err = s.NbsDbInter.Find(k); err != nil {
		return
	}

	old = &BMailMeta{}
	if err = json.Unmarshal([]byte(vs), old); err != nil {
		return nil, err
	}

	newbmm := BMailMeta{}
	newbmm = *old
	newbmm.UpdateTime = tools.GetNowMsTime()

	var v []byte
	if v, err = json.Marshal(newbmm); err != nil {
		return
	} else {
		s.NbsDbInter.Update(k, string(v))
	}
	return
}

func (s *BMailMetaDb) Remove(mHash dbcommon.Hash) (del *BMailMeta, err error) {
	k := base64.StdEncoding.EncodeToString(mHash[:])

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	var vs string
	if vs, err = s.NbsDbInter.Find(k); err != nil {
		return nil, nil
	}

	del = &BMailMeta{}
	err = json.Unmarshal([]byte(vs), del)

	s.NbsDbInter.Delete(k)

	return

}

func (s *BMailMetaDb) Find(mHash dbcommon.Hash) (*BMailMeta, error) {
	k := base64.StdEncoding.EncodeToString(mHash[:])

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if vs, err := s.NbsDbInter.Find(k); err != nil {
		return nil, err
	} else {
		f := &BMailMeta{}
		err = json.Unmarshal([]byte(vs), f)
		if err != nil {
			return nil, err
		} else {
			return f, err
		}
	}
}

func (s *BMailMetaDb) Save() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.NbsDbInter.Save()
}

func (s *BMailMetaDb) Iterator() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.cursor = s.NbsDbInter.DBIterator()
}

func (s *BMailMetaDb) Next() (mHash *dbcommon.Hash, meta *BMailMeta, r1 error) {
	if s.cursor == nil {
		return
	}
	s.dbLock.Lock()
	s.dbLock.Unlock()
	k, v := s.cursor.Next()
	if k == "" {
		s.dbLock.Unlock()
		return nil, nil, nil
	}
	s.dbLock.Unlock()
	meta = &BMailMeta{}

	if err := json.Unmarshal([]byte(v), meta); err != nil {
		return nil, nil, err
	}

	if sk, err := base64.StdEncoding.DecodeString(k); err != nil {
		return nil, nil, err
	} else {
		h := dbcommon.Hash{}
		copy(h[:], sk)
		mHash = &h
	}

	return

}
