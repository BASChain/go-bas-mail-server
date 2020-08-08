package protocol

import (
	"bytes"
	"encoding/json"
	"github.com/realbmail/go-bas-mail-server/bmailcrypt"
	"github.com/realbmail/go-bmail-protocol/bpop"
	"github.com/realbmail/go-bmail-resolver"
	"github.com/mr-tron/base58"
	"log"
)

type CommandDeleteMsg struct {
	Sn        []byte
	CmdSyn    *bpop.CommandSyn
	CmdDelete *bpop.CmdDelete

	CmdDelAck *bpop.CmdDeleteAck
	CmdAck    *bpop.CommandAck
}

func (cdm *CommandDeleteMsg) UnPack(data []byte) error {
	cdm.CmdSyn = &bpop.CommandSyn{}
	cdm.CmdSyn.Cmd = &bpop.CmdDelete{}

	if err := json.Unmarshal(data, cdm.CmdSyn); err != nil {
		return err
	}

	cdm.CmdDelete = cdm.CmdSyn.Cmd.(*bpop.CmdDelete)

	return nil
}

func (cdm *CommandDeleteMsg) Verify() bool {
	if bytes.Compare(cdm.Sn, cdm.CmdSyn.SN[:]) != 0 {
		log.Println("sn not equals", base58.Encode(cdm.Sn), base58.Encode(cdm.CmdSyn.SN[:]))
		return false
	}

	addr, _ := resolver.NewEthResolver(true).BMailBCA(cdm.CmdDelete.MailAddr)
	if addr != cdm.CmdDelete.Owner {
		log.Println("addr not equals", addr, cdm.CmdDelete.Owner)
		return false
	}

	if !bmailcrypt.Verify(addr.ToPubKey(), cdm.CmdSyn.SN[:], cdm.CmdSyn.Sig) {
		log.Println("verify signature failed")
		return false
	}

	return true
}

func (cdm *CommandDeleteMsg) SetCurrentSn(sn []byte) {
	cdm.Sn = sn
}

func (cdm *CommandDeleteMsg) Dispatch() error {
	return nil
}

func (cdm *CommandDeleteMsg) Response() (WBody, error) {
	return nil, nil
}
