package protocol

import (
	"bytes"
	"encoding/json"
	"github.com/BASChain/go-bas-mail-server/bmailcrypt"
	"github.com/BASChain/go-bas-mail-server/db/bmaildb"
	"github.com/BASChain/go-bas-mail-server/db/savefile"
	"github.com/BASChain/go-bas-mail-server/tools"
	"github.com/BASChain/go-bas-mail-server/wallet"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/BASChain/go-bmail-resolver"
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"log"
	"time"
)

type CryptEnvelopeMsg struct {
	Sn      []byte
	EpSyn   *bmp.EnvelopeSyn
	CryptEp *bmp.BMailEnvelope
	//RawEp *bmp.RawEnvelope
	EpAck *bmp.EnvelopeAck
}

func (cem *CryptEnvelopeMsg) UnPack(data []byte) error {
	cem.EpSyn = &bmp.EnvelopeSyn{}
	cem.EpSyn.Env = &bmp.BMailEnvelope{}
	if err := json.Unmarshal(data, cem.EpSyn); err != nil {
		return err
	}

	cem.CryptEp = cem.EpSyn.Env

	return nil
}

func (cem *CryptEnvelopeMsg) Verify() bool {
	if bytes.Compare(cem.Sn, cem.EpSyn.SN[:]) != 0 {
		log.Println("sn not equals ", base58.Encode(cem.Sn), base58.Encode(cem.EpSyn.SN[:]))
		return false
	}

	addr, _ := resolver.NewEthResolver(true).BMailBCA(cem.CryptEp.FromName)
	if addr != cem.CryptEp.FromAddr {
		log.Println("addr not equals", addr, cem.CryptEp.FromAddr)
		return false
	}

	if !bmailcrypt.Verify(addr.ToPubKey(), cem.EpSyn.SN[:], cem.EpSyn.Sig) {
		log.Println("verify signature failed")
		return false
	}

	return true
}

func (cem *CryptEnvelopeMsg) SetCurrentSn(sn []byte) {
	cem.Sn = sn
}

func (cem *CryptEnvelopeMsg) Dispatch() error {

	var size int

	h := cem.CryptEp

	h.DateSince1970 = uint64(time.Now().UnixNano() / 1e6)

	//save mail
	if data, err := json.Marshal(*cem.CryptEp); err == nil {
		size = len(data)
		savefile.Save2File(uuid.MustParse(h.Eid), data)
	} else {
		return err
	}

	//save meta
	mcdb := bmaildb.GetBMMailContentDb()

	if err := mcdb.Insert(uuid.MustParse(h.Eid), cem.CryptEp); err != nil {
		return err
	}

	//save index
	smdb := bmaildb.GetBMSendMailDb()
	smdb.Insert(h.FromName, size, uuid.MustParse(h.Eid), int64(h.DateSince1970))
	pmdb := bmaildb.GetBMPullMailDb()

	for i := 0; i < len(cem.CryptEp.RCPTs); i++ {
		rcpt := cem.CryptEp.RCPTs[i]
		pmdb.Insert(rcpt.ToName, size, uuid.MustParse(h.Eid), int64(h.DateSince1970))

	}

	mcdb.IncRef(uuid.MustParse(h.Eid))
	mcdb.IncRef(uuid.MustParse(h.Eid))

	return nil
}

func (cem *CryptEnvelopeMsg) Response() (WBody, error) {
	ack := &bmp.EnvelopeAck{}

	hash := cem.CryptEp.Hash()
	if bytes.Compare(hash, cem.EpSyn.Hash) != 0 {
		ack.ErrorCode = 1
	} else {
		copy(ack.NextSN[:], tools.NewSn(tools.SerialNumberLength))
		ack.Hash = cem.EpSyn.Hash
		ack.Sig = wallet.GetServerWallet().Sign(cem.EpSyn.Hash)
	}

	cem.EpAck = ack

	return ack, nil

}
