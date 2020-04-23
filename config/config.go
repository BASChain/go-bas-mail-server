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
	BMC_HomeDir      = ".bmc"
	BMC_CFG_FileName = "bmc.json"
)

type BMCConfig struct {
	MgtHttpPort   int             `json:"mgthttpport"`
	KeyPath       string          `json:"keypath"`
	CmdListenPort string          `json:"cmdlistenport"`
	PrivKey       *rsa.PrivateKey `json:"-"`
	PubKey        *rsa.PublicKey  `json:"-"`
	PKAddr        string          `json:"-"`
	RemoteServer  string          `json:"remoteserver"`
}

var (
	bmccfgInst     *BMCConfig
	bmccfgInstLock sync.Mutex
)

func (bc *BMCConfig) InitCfg() *BMCConfig {
	bc.MgtHttpPort = 50818
	bc.KeyPath = "/keystore"
	bc.CmdListenPort = "127.0.0.1:59527"

	return bc
}

func (bc *BMCConfig) Load() *BMCConfig {
	if !tools.FileExists(GetBMCCFGFile()) {
		return nil
	}

	jbytes, err := tools.OpenAndReadAll(GetBMCCFGFile())
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

func newSSCCfg() *BMCConfig {

	bc := &BMCConfig{}

	bc.InitCfg()

	return bc
}

func GetBMCCfg() *BMCConfig {
	if bmccfgInst == nil {
		bmccfgInstLock.Lock()
		defer bmccfgInstLock.Unlock()
		if bmccfgInst == nil {
			bmccfgInst = newSSCCfg()
		}
	}

	return bmccfgInst
}

func PreLoad() *BMCConfig {
	bc := &BMCConfig{}

	return bc.Load()
}

func LoadFromCfgFile(file string) *BMCConfig {
	bc := &BMCConfig{}

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

	bmccfgInstLock.Lock()
	defer bmccfgInstLock.Unlock()
	bmccfgInst = bc

	return bc

}

func LoadFromCmd(initfromcmd func(cmdbc *BMCConfig) *BMCConfig) *BMCConfig {
	bmccfgInstLock.Lock()
	defer bmccfgInstLock.Unlock()

	lbc := newSSCCfg().Load()

	if lbc != nil {
		bmccfgInst = lbc
	} else {
		lbc = newSSCCfg()
	}

	bmccfgInst = initfromcmd(lbc)

	return bmccfgInst
}

func GetBMCHomeDir() string {
	curHome, err := tools.Home()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(curHome, BMC_HomeDir)
}

func GetBMCCFGFile() string {
	return path.Join(GetBMCHomeDir(), BMC_CFG_FileName)
}

func (bc *BMCConfig) Save() {
	jbytes, err := json.MarshalIndent(*bc, " ", "\t")

	if err != nil {
		log.Println("Save BASD Configuration json marshal failed", err)
	}

	if !tools.FileExists(GetBMCHomeDir()) {
		os.MkdirAll(GetBMCHomeDir(), 0755)
	}

	err = tools.Save2File(jbytes, GetBMCCFGFile())
	if err != nil {
		log.Println("Save BASD Configuration to file failed", err)
	}

}

func IsInitialized() bool {
	if tools.FileExists(GetBMCCFGFile()) {
		return true
	}

	return false
}

func (bc *BMCConfig) GetKeyPath() string {
	return path.Join(GetBMCHomeDir(), bc.KeyPath)

}

func (bc *BMCConfig) SetPrivKey(priv *rsa.PrivateKey) {
	bc.PrivKey = priv
}

func (bc *BMCConfig) SetPubKey(pub *rsa.PublicKey) {
	bc.PubKey = pub
}
func (bc *BMCConfig) GetPrivKey() (priv *rsa.PrivateKey) {
	return bc.PrivKey
}

func (bc *BMCConfig) GetPubKey() (pub *rsa.PublicKey) {
	return bc.PubKey
}
