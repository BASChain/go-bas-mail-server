package bmailmemdb

import (
	"github.com/kprc/nbsnetwork/db"
	"sync"
	"net"
	"github.com/BASChain/go-bas-mail-server/db/dbcommon"
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/kprc/nbsnetwork/tools"
	"encoding/json"
	"log"
)

type BMAccessServerDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	asCursor *db.DBCusor
}

var (
	bmasStore *BMAccessServerDb
	bmasStoreLock sync.Mutex
)

type MxServerAddr struct {
	Weight int 	`json:"w"`
	IPAddr net.IP 	`json:"i"`

}

type AccessServer struct {
	TopDomain string `json:"td"`
	PublicKey dbcommon.PKey	`json:"pk"`
	Addr []*MxServerAddr `json:"a"`
	CreateTime int64 	`json:"ct"`
	UpdateTime int64	`json:"ut"`
}

func newBMAccessServer() *BMAccessServerDb {
	cfg := config.GetBMSCfg()

	db:=db.NewFileDb(cfg.GetBMSMSavePath()).Load()

	return &BMAccessServerDb{NbsDbInter:db}

}

func GetBMAccessServer()  *BMAccessServerDb {
	if bmasStore == nil {
		bmasStoreLock.Lock()
		defer bmasStoreLock.Unlock()

		if bmasStore == nil {
			bmasStore = newBMAccessServer()
		}

	}

	return bmasStore
}

func (s *BMAccessServerDb)Insert(tld string,pk dbcommon.PKey,addr []*MxServerAddr) error {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if _,err:=s.NbsDbInter.Find(tld);err==nil{
		return err
	}

	now:=tools.GetNowMsTime()
	as:=&AccessServer{tld,pk,addr,now,now}

	if v,err:=json.Marshal(*as);err!=nil{
		return err
	}else{
		return s.NbsDbInter.Insert(tld,string(v))
	}
}

func (s *BMAccessServerDb)UpdatePK(tld string,pk dbcommon.PKey) (old *AccessServer,err error)  {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	var vs string
	if vs,err = s.NbsDbInter.Find(tld);err!=nil{
		return
	}

	old = &AccessServer{}
	if err = json.Unmarshal([]byte(vs),old);err!=nil{
		return nil,err
	}

	newas := AccessServer{}
	newas = *old

	newas.UpdateTime = tools.GetNowMsTime()

	var v[]byte
	if v,err = json.Marshal(newas);err!=nil{
		return
	}else{
		s.NbsDbInter.Update(tld,string(v))
	}
	return
}

func (s *BMAccessServerDb)AddMXIP(tld string,ip net.IP,weight int) error {
	return s.UpdateMXIP(tld,ip,weight)

}

func (s *BMAccessServerDb)UpdateMXIP(tld string,ip net.IP,weight int) error {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	var as *AccessServer

	var vs string
	var err error
	now := tools.GetNowMsTime()
	if vs,err = s.NbsDbInter.Find(tld);err!=nil{
		as = &AccessServer{TopDomain:tld,CreateTime:now,UpdateTime:now}
	}else{
		as = &AccessServer{}
		if err = json.Unmarshal([]byte(vs),as);err!=nil{
			log.Println("unmarshal AccessServer info failed in function UpdateMxIP")
			return err
		}
		as.UpdateTime = now
	}

	var mxs *MxServerAddr
	for i:=0;i<len(as.Addr);i++{
		m:=as.Addr[i]
		if m.IPAddr.Equal(ip){
			mxs = m
			break
		}
	}

	if mxs == nil{
		m:=&MxServerAddr{IPAddr:ip,Weight:weight}
		as.Addr = append(as.Addr,m)
	}else{
		mxs.Weight = weight
	}

	var v []byte
	if v,err = json.Marshal(as);err!=nil{
		return err
	}else{
		s.NbsDbInter.Update(tld,string(v))
	}

	return nil
}

func (s *BMAccessServerDb)DelMXIP(tld string,ip net.IP) error {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if vs,err:=s.NbsDbInter.Find(tld);err!=nil{
		return nil
	}else{
		as := &AccessServer{}
		if err = json.Unmarshal([]byte(vs),as);err!=nil{
			log.Println("unmarshal AccessServer info failed in function DelMXIP")
			return err
		}

		var i int

		for i=0;i<len(as.Addr);i++{
			m:=as.Addr[i]
			if m.IPAddr.Equal(ip){
				break
			}
		}

		if i == len(as.Addr){
			return nil
		}

		as.Addr = append(as.Addr[:i],as.Addr[i+1:]...)
		var v []byte
		if v,err = json.Marshal(as);err!=nil{
			log.Println("marshal AccessServer info failed in function DelMXIP")
			return err
		}else{
			s.NbsDbInter.Update(tld,string(v))

			return nil
		}
	}

}

func (s *BMAccessServerDb)Remove(tld string) (del *AccessServer,err error)  {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	var vs string
	if vs,err=s.NbsDbInter.Find(tld);err!=nil{
		return nil,nil
	}

	del = &AccessServer{}
	err = json.Unmarshal([]byte(vs),del)

	s.NbsDbInter.Delete(tld)

	return
}

func (s *BMAccessServerDb)Find(tld string) (*AccessServer,error)  {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if vs,err:=s.NbsDbInter.Find(tld);err!=nil{
		return nil,err
	}else{
		f := &AccessServer{}
		err = json.Unmarshal([]byte(vs),f)
		if err!=nil{
			return nil,err
		}else{
			return f,err
		}
	}

}

func (s *BMAccessServerDb)Save()  {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()


	s.NbsDbInter.Save()

}


func (s *BMAccessServerDb)Iterator()  {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.asCursor = s.NbsDbInter.DBIterator()
}


func (s *BMAccessServerDb)Next() (tld string,meta *AccessServer,r1 error)  {
	if s.asCursor == nil{
		return
	}
	s.dbLock.Lock()
	s.dbLock.Unlock()
	k,v:=s.asCursor.Next()
	if tld == ""{
		s.dbLock.Unlock()
		return "",nil,nil
	}
	s.dbLock.Unlock()
	meta = &AccessServer{}

	if err := json.Unmarshal([]byte(v),meta);err!=nil{
		return "",nil,err
	}

	tld = k

	return

}
