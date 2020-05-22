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
)

type CryptEnvelopeMsg struct {
	Sn      []byte
	EpSyn   *bmp.EnvelopeSyn
	CryptEp *bmp.CryptEnvelope
	//RawEp *bmp.RawEnvelope
	EpAck *bmp.EnvelopeAck
}

func (cem *CryptEnvelopeMsg) UnPack(data []byte) error {
	cem.EpSyn = &bmp.EnvelopeSyn{}
	cem.EpSyn.Env = &bmp.CryptEnvelope{}
	if err := json.Unmarshal(data, cem.EpSyn); err != nil {
		return err
	}

	cem.CryptEp = cem.EpSyn.Env.(*bmp.CryptEnvelope)

	return nil
}

func (cem *CryptEnvelopeMsg) Verify() bool {
	if bytes.Compare(cem.Sn, cem.EpSyn.SN[:]) != 0 {
		return false
	}

	addr, _ := resolver.NewEthResolver(true).BMailBCA(cem.CryptEp.EnvelopeHead.From)
	if addr != cem.CryptEp.FromAddr {
		return false
	}

	if !bmailcrypt.Verify(addr.ToPubKey(), cem.EpSyn.SN[:], cem.EpSyn.Sig) {
		return false
	}

	toaddr, _ := resolver.NewEthResolver(true).BMailBCA(cem.CryptEp.EnvelopeHead.To)
	if toaddr != cem.CryptEp.ToAddr {
		return false
	}

	return true
}

func (cem *CryptEnvelopeMsg) SetCurrentSn(sn []byte) {
	cem.Sn = sn
}

func (cem *CryptEnvelopeMsg) Dispatch() error {

	//save meta
	mcdb := bmaildb.GetBMMailContentDb()
	h := &(cem.CryptEp.EnvelopeHead)
	if err := mcdb.Insert(h.Eid, h.From, h.FromAddr, h.To, h.ToAddr); err != nil {
		return err
	}

	var size int

	if data, err := json.Marshal(*cem.CryptEp); err != nil {
		size = len(data)
		savefile.Save2File(h.Eid, data)
	}

	smdb := bmaildb.GetBMSendMailDb()
	smdb.Insert(h.FromAddr.String(), size, h.Eid)
	pmdb := bmaildb.GetBMPullMailDb()
	pmdb.Insert(h.ToAddr.String(), size, h.Eid)

	mcdb.IncRef(h.Eid)
	mcdb.IncRef(h.Eid)

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
