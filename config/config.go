package config

import (
	"crypto/ed25519"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"os"
	"path"
	"sync"
)

const (
	BMS_HomeDir      = ".bms"
	BMS_CFG_FileName = "bms.json"
)

type BMSConfig struct {
	MgtHttpPort   int                `json:"mgthttpport"`
	KeyPath       string             `json:"keypath"`
	CmdListenPort string             `json:"cmdlistenport"`
	PrivKey       ed25519.PrivateKey `json:"-"`
	PubKey        ed25519.PublicKey  `json:"-"`
	RemoteServer  string             `json:"remoteserver"`
	PKStorePath   string             `json:"publickeypath"`
	SSListenPort  int                `json:"sslistenport"`

	DbPath          string `json:"dbpath"`
	BMailMetaDb     string `json:"bmailmetadb"`
	BMServerMetaDb  string `json:"bmservermetadb"`
	BMTransferDb    string `json:"bmtransferdb"`
	BMBlackListdb   string `json:"bmblacklist"`
	BMAttachDb      string `json:"bmattachdb"`
	BMSendMailDb    string `json:"bmsendmaildb"`
	BMMailContentDb string `json:"bmmailcontentdb"`
	BMPullMailDb    string `json:"bmpullmaildb"`

	FileStorePath string `json:"filestorepath"`
}

var (
	bmscfgInst     *BMSConfig
	bmscfgInstLock sync.Mutex
)

func (bc *BMSConfig) InitCfg() *BMSConfig {
	bc.MgtHttpPort = 50818
	bc.KeyPath = "/keystore"
	bc.CmdListenPort = "127.0.0.1:59529"
	bc.PKStorePath = "/pkstore"
	bc.SSListenPort = 50021

	bc.DbPath = "/db"
	bc.BMailMetaDb = "bmm.db"
	bc.BMServerMetaDb = "bmsm.db"
	bc.BMTransferDb = "bmtf.db"
	bc.BMBlackListdb = "bmbl.db"
	bc.BMAttachDb = "bmattch.db"
	bc.BMSendMailDb = "sendmail.db"
	bc.BMMailContentDb = "mailcontent.db"
	bc.BMPullMailDb = "pullmail.db"
	bc.FileStorePath = "mailstore"

	return bc
}

func (bc *BMSConfig) Load() *BMSConfig {
	if !tools.FileExists(GetBMSCFGFile()) {
		return nil
	}

	jbytes, err := tools.OpenAndReadAll(GetBMSCFGFile())
	if err != nil {
		log.Println("load file failed", err)
		return nil
	}

	err = json.Unmarshal(jbytes, bc)
	if err != nil {
		log.Println("load configuration unmarshal failed", err)
		return nil
	}

	return bc

}

func newBMSCfg() *BMSConfig {

	bc := &BMSConfig{}

	bc.InitCfg()

	return bc
}

func GetBMSCfg() *BMSConfig {
	if bmscfgInst == nil {
		bmscfgInstLock.Lock()
		defer bmscfgInstLock.Unlock()
		if bmscfgInst == nil {
			bmscfgInst = newBMSCfg()
		}
	}

	return bmscfgInst
}

func PreLoad() *BMSConfig {
	bc := &BMSConfig{}

	return bc.Load()
}

func LoadFromCfgFile(file string) *BMSConfig {
	bc := &BMSConfig{}

	bc.InitCfg()

	bcontent, err := tools.OpenAndReadAll(file)
	if err != nil {
		log.Fatal("Load Config file failed")
		return nil
	}

	err = json.Unmarshal(bcontent, bc)
	if err != nil {
		log.Fatal("Load Config From json failed")
		return nil
	}

	bmscfgInstLock.Lock()
	defer bmscfgInstLock.Unlock()
	bmscfgInst = bc

	return bc

}

func LoadFromCmd(initfromcmd func(cmdbc *BMSConfig) *BMSConfig) *BMSConfig {
	bmscfgInstLock.Lock()
	defer bmscfgInstLock.Unlock()

	lbc := newBMSCfg().Load()

	if lbc != nil {
		bmscfgInst = lbc
	} else {
		lbc = newBMSCfg()
	}

	bmscfgInst = initfromcmd(lbc)

	return bmscfgInst
}

func GetBMSHomeDir() string {
	curHome, err := tools.Home()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(curHome, BMS_HomeDir)
}

func GetBMSCFGFile() string {
	return path.Join(GetBMSHomeDir(), BMS_CFG_FileName)
}

func (bc *BMSConfig) Save() {
	jbytes, err := json.MarshalIndent(*bc, " ", "\t")

	if err != nil {
		log.Println("Save BASD Configuration json marshal failed", err)
	}

	if !tools.FileExists(GetBMSHomeDir()) {
		os.MkdirAll(GetBMSHomeDir(), 0755)
	}

	err = tools.Save2File(jbytes, GetBMSCFGFile())
	if err != nil {
		log.Println("Save BASD Configuration to file failed", err)
	}

}

func IsInitialized() bool {
	if tools.FileExists(GetBMSCFGFile()) {
		return true
	}

	return false
}

func (bc *BMSConfig) GetKeyPath() string {
	return path.Join(GetBMSHomeDir(), bc.KeyPath)

}

func (bc *BMSConfig) GetPKPath() string {
	return path.Join(GetBMSHomeDir(), bc.PKStorePath)
}

func (bc *BMSConfig) GetDbPath() string {
	dbpath := path.Join(GetBMSHomeDir(), bc.DbPath)
	if !tools.FileExists(dbpath) {
		os.MkdirAll(dbpath, 0755)
	}

	return dbpath
}

func (bc *BMSConfig) GetBMMSavePath() string {
	return path.Join(bc.GetDbPath(), bc.BMailMetaDb)
}

func (bc *BMSConfig) GetBMSMSavePath() string {
	return path.Join(bc.GetDbPath(), bc.BMServerMetaDb)
}

func (bc *BMSConfig) GetBMMLSavePath() string {
	return path.Join(bc.GetDbPath(), bc.BMBlackListdb)
}

func (bc *BMSConfig) GetBMTransferSavePath() string {
	return path.Join(bc.GetDbPath(), bc.BMTransferDb)
}

func (bc *BMSConfig) GetAttachmentSavePath() string {
	return path.Join(bc.GetDbPath(), bc.BMAttachDb)
}

func (bc *BMSConfig) GetSendMailDBPath() string {
	return path.Join(bc.GetDbPath(), bc.BMSendMailDb)
}

func (bc *BMSConfig) GetMailContentDBPath() string {
	return path.Join(bc.GetDbPath(), bc.BMMailContentDb)
}

func (bc *BMSConfig) GetPullMailDbPath() string {
	return path.Join(bc.GetDbPath(), bc.BMPullMailDb)
}

func (bc *BMSConfig) GetMailStorePath() string {
	dbpath := path.Join(GetBMSHomeDir(), bc.FileStorePath)
	if !tools.FileExists(dbpath) {
		os.MkdirAll(dbpath, 0755)
	}

	return dbpath
}
