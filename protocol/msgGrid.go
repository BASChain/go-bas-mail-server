package protocol

import (
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"encoding/json"

	"crypto/rand"
)

type MsgBody interface {
	UnPack(data []byte) error
	//Save2DB()	error
	Response() (*bmp.EnvelopeAck,error)
}


type CryptEnvelopeMsg struct {
	EpSyn *bmp.EnvelopeSyn
	CryptEp *bmp.CryptEnvelope
	RawEp *bmp.RawEnvelope
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
//
//func (cem *CryptEnvelopeMsg)Verify() bool  {
//	//todo
//
//	return true
//}


//func (cem *CryptEnvelopeMsg)Save2DB() error  {
//	//todo...
//
//	return nil
//}

func newSn() []byte  {
	sn := make([]byte, 16)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	return sn
}

func (cem *CryptEnvelopeMsg)Response() (*bmp.EnvelopeAck,error){
	ack:=&bmp.EnvelopeAck{}
	copy(ack.NextSN[:],newSn())
	ack.Hash = cem.EpSyn.Hash
	ack.Sig = newSn()

	cem.EpAck = ack

	return ack,nil

}



var MsgGrid = map[uint16]MsgBody{
	translayer.SEND_CRYPT_ENVELOPE:&CryptEnvelopeMsg{},
}






