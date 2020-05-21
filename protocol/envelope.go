package protocol

import (
	"github.com/BASChain/go-bmail-protocol/bmp"
	"encoding/json"
	"github.com/BASChain/go-bas-mail-server/wallet"
	"bytes"
	"github.com/BASChain/go-bas-mail-server/tools"
)

type CryptEnvelopeMsg struct {
	EpSyn *bmp.EnvelopeSyn
	CryptEp *bmp.CryptEnvelope
	//RawEp *bmp.RawEnvelope
	EpAck *bmp.EnvelopeAck
}

func (cem *CryptEnvelopeMsg)UnPack(data []byte) error  {
	cem.EpSyn = &bmp.EnvelopeSyn{}
	cem.EpSyn.Env = &bmp.CryptEnvelope{}
	if err:=json.Unmarshal(data,cem.EpSyn);err!=nil{
		return err
	}

	cem.CryptEp = cem.EpSyn.Env.(*bmp.CryptEnvelope)

	return nil
}

func (cem *CryptEnvelopeMsg)Verify() bool  {
	//todo

	return true
}

//func (cem *CryptEnvelopeMsg)Save2DB() error  {
//	//todo...
//
//	return nil
//}

func (cem *CryptEnvelopeMsg)Response() (WBody,error){
	ack:=&bmp.EnvelopeAck{}

	hash:=cem.CryptEp.Hash()
	if bytes.Compare(hash,cem.EpSyn.Hash) != 0{
		ack.ErrorCode = 1
	}else{
		copy(ack.NextSN[:],tools.NewSn(tools.SerialNumberLength))
		ack.Hash = cem.EpSyn.Hash
		ack.Sig = wallet.GetServerWallet().Sign(cem.EpSyn.Hash)
	}


	cem.EpAck = ack

	return ack,nil

}

