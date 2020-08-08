package instdb

import (
	"crypto/rand"
	"github.com/realbmail/go-bas-mail-server/kvdb"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Session struct {
	sn     string
	pubkey string
}

var (
	sessdb     kvdb.KVDBInterface
	sessdblock sync.Mutex
	sessdbquit chan int
	sesswg     sync.WaitGroup
)

func (sess *Session) GetSn() string {
	return sess.sn
}

func (sess *Session) GetPubKey() string {
	return sess.pubkey
}

func (sess *Session) SetSn(sn string) {
	sess.sn = sn
}

func (sess *Session) SetPubKey(pk string) {
	sess.pubkey = pk
}

func GetSessionDB() kvdb.KVDBInterface {
	if sessdb != nil {
		return sessdb
	}

	sessdblock.Lock()
	defer sessdblock.Unlock()

	if sessdb != nil {
		return sessdb
	}

	sessdb = kvdb.NewKVDB("session", "")
	sessdbquit = make(chan int, 1)

	return sessdb
}

func GenSerialNumber() string {

	sn := make([]byte, 32, 32)

	for {
		if _, err := rand.Read(sn); err != nil {
			continue
		}
		break
	}

	return base58.Encode(sn)
}

func NewSession(pubkey string) *Session {
	sess := &Session{}

	sess.sn = GenSerialNumber()
	sess.pubkey = pubkey

	return sess
}

func InserSession(sess *Session) error {
	store := GetSessionDB()

	if sess == nil {
		return errors.New("sess is null")
	}

	if _, err := store.Find(sess.sn); err == nil {
		return errors.New("have save session")
	}

	store.Store(sess.sn, sess.pubkey)

	return nil
}

func FindSession(sess *Session) (r *Session, err error) {

	if sess == nil || sess.sn == "" {
		return nil, errors.New("sess is null or searial number error")
	}

	store := GetSessionDB()

	if v, err := store.Find(sess.sn); err != nil {
		return nil, errors.New("no session in db")
	} else {
		r := &Session{}
		r.sn = sess.sn
		r.pubkey = v

		return r, nil
	}
}

func DelSession(sess *Session) {
	if sess == nil {
		return
	}
	store := GetSessionDB()

	store.DelK(sess.sn)
}

func SessionTimeOut() {
	sesswg.Add(1)
	defer sesswg.Done()
	store := GetSessionDB()

	type SessTimeOut struct {
		to      []string
		curtime int64
	}

	lasttime := tools.GetNowMsTime()

	for {
		select {
		case <-sessdbquit:
			return
		default:
			//nothing to do
		}

		curtime := tools.GetNowMsTime()

		if curtime-lasttime < 30000 {
			continue
		}

		lasttime = curtime

		st := &SessTimeOut{}
		st.curtime = curtime

		store.TraversDo(st, func(arg interface{}, k, v interface{}) {
			s := arg.(*SessTimeOut)

			dbv := v.(*kvdb.DBV)

			if s.curtime-dbv.GetAccessTime() > 300000 {
				s.to = append(s.to, k.(string))
			}

		})

		for i := 0; i < len(st.to); i++ {
			store.DelK(st.to[i])
		}

		time.Sleep(time.Second * 1)
	}
}

func SessTimeOutStop() {
	sessdbquit <- 1
	sesswg.Wait()
}
