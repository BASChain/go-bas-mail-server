package bmaildb

import (
	"bytes"
	"encoding/json"
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/google/uuid"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"sync"
)

type BMSaveMailDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

var (
	sendMailStore     *BMSaveMailDb
	sendMailStoreLock sync.Mutex
)

func newBMSendMailDb() *BMSaveMailDb {
	cfg := config.GetBMSCfg()

	db := db.NewFileDb(cfg.GetSendMailDBPath()).Load()

	return &BMSaveMailDb{NbsDbInter: db}
}

func GetBMSendMailDb() *BMSaveMailDb {
	if sendMailStore == nil {
		sendMailStoreLock.Lock()
		defer sendMailStoreLock.Unlock()

		if sendMailStore == nil {
			sendMailStore = newBMSendMailDb()
		}
	}

	return sendMailStore
}

type MailItem struct {
	Eid        uuid.UUID `json:"eid"`
	Size       int       `json:"size"`
	CreateTime int64     `json:"create_time"`
}

type SaveMail struct {
	Owner     string      `json:"-"`
	TotalSize int64       `json:"total_size"`
	TotalCnt  int         `json:"total_cnt"`
	Smi       []*MailItem `json:"smi"`
}

func (smdb *BMSaveMailDb) Insert(owner string, size int, eid uuid.UUID, curTime int64) error {
	smdb.dbLock.Lock()
	defer smdb.dbLock.Unlock()

	sm := &SaveMail{}

	if v, err := smdb.NbsDbInter.Find(owner); err == nil {
		if err = json.Unmarshal([]byte(v), sm); err != nil {
			return err
		}
	}

	for i := 0; i < len(sm.Smi); i++ {
		if bytes.Compare(eid[:], sm.Smi[i].Eid[:]) == 0 {
			return errors.New("eid duplicated")
		}
	}

	sm.TotalSize += int64(size)
	sm.TotalCnt++

	smi := &MailItem{}

	if curTime <= 0 {
		curTime = tools.GetNowMsTime()
	}

	smi.CreateTime = curTime
	smi.Eid = eid
	smi.Size = size

	sm.Smi = append(sm.Smi, smi)

	if v, err := json.Marshal(*sm); err != nil {
		return err
	} else {
		smdb.NbsDbInter.Update(owner, string(v))
	}

	return nil
}

func (smdb *BMSaveMailDb) Find(owner string) (sm *SaveMail, err error) {
	smdb.dbLock.Lock()
	defer smdb.dbLock.Unlock()

	sm = &SaveMail{}

	if v, err := smdb.NbsDbInter.Find(owner); err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal([]byte(v), sm); err != nil {
			return nil, err
		}
		return sm, nil
	}

}

func (smdb *BMSaveMailDb) Delete(owner string, eid uuid.UUID) error {
	smdb.dbLock.Lock()
	defer smdb.dbLock.Unlock()

	sm := &SaveMail{}

	if v, err := smdb.NbsDbInter.Find(owner); err != nil {
		return nil
	} else {
		if err = json.Unmarshal([]byte(v), sm); err != nil {
			return err
		}

		var pos int = -1

		for i := 0; i < len(sm.Smi); i++ {
			if bytes.Compare(eid[:], sm.Smi[i].Eid[:]) == 0 {
				pos = i
				break
			}
		}

		if pos > -1 {
			smi := sm.Smi[pos]
			sm.TotalCnt--
			sm.TotalSize -= int64(smi.Size)

			smis := sm.Smi

			sm.Smi = []*MailItem{}

			for i := 0; i < len(smis); i++ {
				if i != pos {
					sm.Smi = append(sm.Smi, smis[i])
				}

			}
		}
		var vv []byte
		if vv, err = json.Marshal(*sm); err != nil {
			return err
		} else {
			smdb.NbsDbInter.Update(owner, string(vv))
		}

	}

	return nil
}
