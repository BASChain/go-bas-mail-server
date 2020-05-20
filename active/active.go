package active

import (
	"encoding/json"
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/BASChain/go-bas-mail-server/rsakey"
	"github.com/btcsuite/btcutil/base58"
	//"github.com/rickeyliao/ServiceAgent/common"
	"log"
)

type WhiteListReq struct {
	PubKey string `json:"pubkey"`
	Sn     string `json:"sn"`
	Sig    string `json:"sig"`
	Step   int    `json:"step"`
}

type WhiteListResp struct {
	Sn       string `json:"sn"`
	Step     int    `json:"step"`
	State    int    `json:"state"`
	ServerPk string `json:"serverpk"`
}

func ActiveVPN() {
	cfg := config.GetBMSCfg()

	req := &WhiteListReq{}

	req.Step = 1

	var bpk []byte
	var err error

	bpk, err = rsakey.PubKeyToBytes(cfg.PubKey)
	if err != nil {
		log.Println("public key error")
		return
	}

	pk := base58.Encode(bpk)

	req.PubKey = pk

	jsonstr, err := json.Marshal(*req)
	if err != nil {
		log.Println("Marshall json error", err,string(jsonstr))
		return
	}
	var r string
	var code int
	//r, code, err = common.Post1("http://"+cfg.RemoteServer+":"+strconv.Itoa(cfg.MgtHttpPort)+"/ajax/chg", string(jsonstr), false)

	if err != nil || code != 200 {
		log.Println("step 1 failed", err)
		return
	}

	//log.Println("reponse1:",string(r))
	resp := &WhiteListResp{}

	err = json.Unmarshal([]byte(r), resp)
	if err != nil {
		log.Println("step 1 response failed", err)
		return
	}
	bsn := base58.Decode(resp.Sn)

	var bsig []byte

	bsig, err = rsakey.SignRSA(bsn, cfg.PrivKey)
	if err != nil {
		log.Println("sig failed")
		return
	}

	sig := base58.Encode(bsig)

	req = &WhiteListReq{}

	req.Sn = resp.Sn
	req.Step = 2
	req.Sig = sig

	jsonstr, err = json.Marshal(*req)
	if err != nil {
		log.Println("Marshall json error,step 2", err,string(jsonstr))
		return
	}

	//r, code, err = common.Post1("http://"+cfg.RemoteServer+":"+strconv.Itoa(cfg.MgtHttpPort)+"/ajax/chg", string(jsonstr), false)

	if err != nil || code != 200 {
		log.Println("step 2 failed", err)
		return
	}

	log.Println("Success")

}
