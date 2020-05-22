package bmaildb

import (
	"encoding/base64"
	"encoding/json"
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/BASChain/go-bas-mail-server/db/dbcommon"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"sync"
	"sync/atomic"
)

type BMAttachmentDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor

	totalSize  int64
	totalCount int32
}

var (
	bmatmStore     *BMAttachmentDb
	bmatmStoreLock sync.Mutex
)

type AttahmentMeta struct {
	SendHash   dbcommon.Hash
	RecvHash   dbcommon.Hash
	SendPubKey dbcommon.PKey
	FileName   string
	FileTye    int
	FileSize   int
	FilePath   string
	CreateTime int64
	UpdateTime int64
}

func newBMAttachmentDB() *BMAttachmentDb {
	cfg := config.GetBMSCfg()
	db := db.NewFileDb(cfg.GetAttachmentSavePath()).Load()

	return (&BMAttachmentDb{NbsDbInter: db}).init()
}

func (a *BMAttachmentDb) init() *BMAttachmentDb {
	if a.NbsDbInter == nil {
		log.Fatal("init Attachment db failure")
		return nil
	}

	cursor := a.NbsDbInter.DBIterator()

	if cursor == nil {
		for {
			k, v := cursor.Next()
			if k == "" {
				break
			}

			at := &AttahmentMeta{}
			err := json.Unmarshal([]byte(v), at)
			if err != nil {
				log.Println("unmarshal " + v + " failed")
				continue
			}

			atomic.AddInt32(&a.totalCount, 1)
			atomic.AddInt64(&a.totalSize, int64(at.FileSize))

		}
	}

	return a
}

func GetBMAttachmentStore() *BMAttachmentDb {
	if bmatmStore == nil {
		bmatmStoreLock.Lock()
		defer bmatmStoreLock.Unlock()

		if bmatmStore == nil {
			bmatmStore = newBMAttachmentDB()
		}

	}

	return bmatmStore
}

func (s *BMAttachmentDb) Insert(fileHash dbcommon.Hash, am *AttahmentMeta) error {
	k := base64.StdEncoding.EncodeToString(fileHash[:])

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if _, err := s.NbsDbInter.Find(k); err == nil {
		return err
	}

	now := tools.GetNowMsTime()

	am.CreateTime = now
	am.UpdateTime = now

	return nil

}
