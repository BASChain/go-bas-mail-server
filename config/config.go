package config

import (
	"crypto/rsa"
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
	MgtHttpPort   int             `json:"mgthttpport"`
	KeyPath       string          `json:"keypath"`
	CmdListenPort string          `json:"cmdlistenport"`
	PrivKey       *rsa.PrivateKey `json:"-"`
	PubKey        *rsa.PublicKey  `json:"-"`
	PKAddr        string          `json:"-"`
	RemoteServer  string          `json:"remoteserver"`
}

var (
	bmscfgInst     *BMSConfig
	bmscfgInstLock sync.Mutex
)

func (bc *BMSConfig) InitCfg() *BMSConfig {
	bc.MgtHttpPort = 50818
	bc.KeyPath = "/keystore"
	bc.CmdListenPort = "127.0.0.1:59527"

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

func (bc *BMSConfig) SetPrivKey(priv *rsa.PrivateKey) {
	bc.PrivKey = priv
}

func (bc *BMSConfig) SetPubKey(pub *rsa.PublicKey) {
	bc.PubKey = pub
}
func (bc *BMSConfig) GetPrivKey() (priv *rsa.PrivateKey) {
	return bc.PrivKey
}

func (bc *BMSConfig) GetPubKey() (pub *rsa.PublicKey) {
	return bc.PubKey
}
