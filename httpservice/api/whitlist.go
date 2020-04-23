package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
	"net/http"
	//"github.com/kprc/ssactiveserver/httpservice"
	"crypto/rsa"
	"github.com/kprc/ssactiveserver/config"
	"github.com/kprc/ssactiveserver/instdb"
	"github.com/kprc/ssactiveserver/kvdb"
	"github.com/kprc/ssactiveserver/rsakey"
	"os/exec"
	"strings"
	"log"
	"strconv"
)

type WhiteList struct {
}

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

func NewWhiteList() *WhiteList {
	return &WhiteList{}
}

func (wl *WhiteList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	var body []byte
	var err error

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	req := &WhiteListReq{}

	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	var resp *WhiteListResp
	var rpk string

	if req.Step == 1 {

		resp = Step1(req)
		if resp == nil {
			log.Println("reponse step 1 failed",string(body))
			w.WriteHeader(500)
			fmt.Fprintf(w, "{}")
			return
		}
	} else if req.Step == 2 {
		resp, rpk = Step2(req)
		if resp == nil {
			log.Println("step 2 failed")
			w.WriteHeader(500)
			fmt.Fprintf(w, "{}")
			return
		}
	} else {
		log.Println("step erro")
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	//log.Println(rpk)

	if req.Step == 2 {
		pkdb := instdb.GetPKDB()
		brpk := base58.Decode(rpk)
		s := sha256.Sum256(brpk)
		k := base58.Encode(s[:])
		ipold := pkdb.GetK(k)

		ipn, _ := getremoteaddr(r.RemoteAddr)

		if ipold == "" || ipold != ipn {

			err := changeWhitIP(k,ipold, ipn)
			if err == nil {
				pkdb.Update(k, ipn)
			}
			pkdb.Save()
		}
	}

	var bresp []byte

	bresp, err = json.Marshal(*resp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	//log.Println("response:",string(bresp))
	w.WriteHeader(200)
	w.Write(bresp)

}

func changeWhitIP(key,oldip string, newip string) error {
	pkdb := instdb.GetPKDB()

	type ChgIPs struct {
		oldip  string
		newip  string
		delold bool
		addnew bool
	}

	c := &ChgIPs{}
	c.oldip = oldip
	c.newip = newip
	if oldip == "" {
		c.delold = false
	} else {
		c.delold = true
	}
	c.addnew = true

	pkdb.TraversDo(c, func(arg interface{}, k, v interface{}) {
		chg := arg.(*ChgIPs)
		dbv := v.(*kvdb.DBV)

		if dbv.GetData() == chg.newip {
			chg.addnew = false
		}

		if chg.oldip != "" && k.(string)!=key && dbv.GetData() == chg.oldip {
			chg.delold = false
		}
	})


	if c.addnew {
		cmd := exec.Command("/bin/sh","-c",
							"/usr/bin/python2 -Es"+
							" /usr/bin/firewall-cmd "+
							"--permanent "+
							"--add-rich-rule='rule family=\"ipv4\" source address=\""+
							newip+
							"\" port port=\""+
							strconv.Itoa(config.GetSSSCfg().SSListenPort)+
							"\" protocol=\"tcp\" accept'")
		err := cmd.Run()
		if err != nil {
			log.Println("add",err)
			return err
		}
	}

	if c.delold {
		cmd := exec.Command("/bin/sh","-c",
							"/usr/bin/python2 -Es "+
							"/usr/bin/firewall-cmd "+
							"--permanent "+
							"--remove-rich-rule='rule family=\"ipv4\" source address=\""+
							oldip+
							"\" port port=\""+
							strconv.Itoa(config.GetSSSCfg().SSListenPort)+
							"\" protocol=\"tcp\" accept'")
		err := cmd.Run()
		if err != nil {
			log.Println("del",err)
			return nil
		}
	}

	if c.addnew || c.delold{
		cmd := exec.Command("/bin/sh","-c",
							"/usr/bin/python2 -Es "+
							"/usr/bin/firewall-cmd --reload")
		err := cmd.Run()
		if err != nil {
			log.Println("reload",err)
			return nil
		}
	}
	return nil
}

func getremoteaddr(remoteaddr string) (ip string, port string) {
	if remoteaddr == "" {
		return
	}

	ra := strings.Split(remoteaddr, ":")

	if len(ra) != 2 {
		return
	}

	ip = ra[0]
	port = ra[1]

	return
}

func Step1(req *WhiteListReq) *WhiteListResp {
	if req.PubKey == "" {
		return nil
	}
	bpk := base58.Decode(req.PubKey)
	if len(bpk) == 0 {
		return nil
	}
	s := sha256.Sum256(bpk)

	k := base58.Encode(s[:])

	db := instdb.GetPKDB()
	if _, err := db.Find(k); err != nil {
		return nil
	}

	resp := &WhiteListResp{}

	resp.Step = 1
	resp.State = 0
	sess := instdb.NewSession(req.PubKey)

	resp.Sn = sess.GetSn()

	//resp.ServerPk
	if bsrvpk, err := rsakey.PubKeyToBytes(config.GetSSSCfg().PubKey); err != nil {
		return nil
	} else {
		resp.ServerPk = base58.Encode(bsrvpk)
	}

	if err := instdb.InserSession(sess); err != nil {
		return nil
	}

	return resp
}

func Step2(req *WhiteListReq) (*WhiteListResp, string) {
	if req.Sn == "" {
		return nil, ""
	}

	sessp := &instdb.Session{}
	sessp.SetSn(req.Sn)

	var sess *instdb.Session
	var err error
	sess, err = instdb.FindSession(sessp)
	if err != nil {
		log.Println("sess not found")
		return nil, ""
	}
	var pk *rsa.PublicKey
	bpk := base58.Decode(sess.GetPubKey())
	pk, err = rsakey.ParsePubKey(bpk)
	if err != nil {
		log.Println("parse pub key failed")
		return nil, ""
	}

	bsn := base58.Decode(req.Sn)
	bsig := base58.Decode(req.Sig)

	err = rsakey.VerifyRSA(bsn, bsig, pk)
	if err != nil {
		log.Println("sign verify failed")
		return nil, ""
	}

	//del session
	instdb.DelSession(sess)

	resp := &WhiteListResp{}
	resp.Step = 2
	resp.State = 0

	return resp, sess.GetPubKey()

}
