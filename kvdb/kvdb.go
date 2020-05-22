package kvdb

import (
	"encoding/json"
	"errors"
	"github.com/BASChain/go-bas-mail-server/rsakey"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"sync"
)

type KVDBInterface interface {
	Store(k, v string) error
	Update(k, v string)
	GetK(k string) (v string)
	Find(k string) (v string, err error)
	DelK(k string)
	TraversDo(arg interface{}, f func(arg interface{}, k, v interface{}))
	Save()
	Load()
}

type DBV struct {
	path       string
	line       int
	data       string
	accessTime int64
}

type KV struct {
	K string `json:"k"`
	V string `json:"v,omitempty"`
	P string `json:"p,omitempty"`
	L int    `json:"l,omitempty"`
}

type KVS struct {
	Kvs []*KV `json:"kvs"`
}

type KVDB struct {
	lock     sync.Mutex
	password string
	savepath string
	db       map[string]*DBV
}

func NewKVDB(p string, savepath string) *KVDB {

	kv := &KVDB{}
	kv.password = p
	kv.db = make(map[string]*DBV, 0)
	kv.savepath = savepath
	kv.Load()

	return kv
}

func (v *DBV) GetAccessTime() int64 {
	return v.accessTime
}

func (v *DBV) GetData() string {
	return v.data
}

func (kvdb *KVDB) Store(k, v string) error {
	kvdb.lock.Lock()
	defer kvdb.lock.Unlock()

	if _, ok := kvdb.db[k]; ok {
		return errors.New("key duplicate")
	}

	dbv := &DBV{data: v}
	dbv.accessTime = tools.GetNowMsTime()

	kvdb.db[k] = dbv

	return nil
}

func (kvdb *KVDB) Update(k, v string) {
	kvdb.lock.Lock()
	defer kvdb.lock.Unlock()
	dbv := &DBV{data: v}
	dbv.accessTime = tools.GetNowMsTime()

	kvdb.db[k] = dbv

}

func (kvdb *KVDB) GetK(k string) (v string) {
	kvdb.lock.Lock()
	defer kvdb.lock.Unlock()

	if data, ok := kvdb.db[k]; !ok {
		return ""
	} else {
		return data.data
	}
}

func (kvdb *KVDB) Find(k string) (v string, err error) {
	kvdb.lock.Lock()
	defer kvdb.lock.Unlock()

	if data, ok := kvdb.db[k]; !ok {
		return "", errors.New("Not found")
	} else {
		return data.data, nil
	}
}

func (kvdb *KVDB) DelK(k string) {
	kvdb.lock.Lock()
	defer kvdb.lock.Unlock()

	delete(kvdb.db, k)
}

func (kvdb *KVDB) TraversDo(arg interface{}, f func(arg interface{}, k, v interface{})) {
	kvdb.lock.Lock()
	defer kvdb.lock.Unlock()

	for key, val := range kvdb.db {
		f(arg, key, val)
	}
}

func (kvdb *KVDB) Save() {
	if kvdb.savepath == "" {
		log.Println("save failed")
		return
	}

	kvs := &KVS{}

	for k, v := range kvdb.db {
		kv := &KV{}
		kv.K = k
		kv.V = v.data
		kv.P = v.path
		kv.L = v.line

		kvs.Kvs = append(kvs.Kvs, kv)
	}

	jb, err := json.Marshal(*kvs)
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println(kvdb.savepath, kvdb.password, string(jb))
	tosave := rsakey.AesEncrypt(jb, []byte(kvdb.password))
	err = tools.Save2File(tosave, kvdb.savepath)
	if err != nil {
		log.Println(err)
	}

}

func (kvdb *KVDB) Load() {

	log.Println(kvdb.savepath, kvdb.password)
	data, err := tools.OpenAndReadAll(kvdb.savepath)
	if err != nil {
		log.Println("no kvdb")
		return
	}

	tojson := rsakey.AesDecrypt(data, []byte(kvdb.password))

	if tojson == nil {
		log.Println("decrypt data from db failed")
		return
	}

	log.Println(string(tojson))

	kvs := &KVS{}

	err = json.Unmarshal(tojson, kvs)
	if err != nil {
		log.Println("unmarshal json failed", err)
		return
	}

	kvdb.lock.Lock()
	defer kvdb.lock.Unlock()

	for i := 0; i < len(kvs.Kvs); i++ {
		dbv := &DBV{}
		dbv.data = kvs.Kvs[i].V
		dbv.path = kvs.Kvs[i].P
		dbv.line = kvs.Kvs[i].L

		kvdb.db[kvs.Kvs[i].K] = dbv
	}

	return
}
