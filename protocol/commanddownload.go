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
	"github.com/BASChain/go-bmail-protocol/bpop"
	"github.com/BASChain/go-bmail-resolver"
	"github.com/google/uuid"
)

type CommandDownloadMsg struct {
	Sn          []byte
	CmdSyn      *bpop.CommandSyn
	CmdDownload *bpop.CmdDownload

	CmdDownAck *bpop.CmdDownloadAck
	CmdAck     *bpop.CommandAck
}

func (cdm *CommandDownloadMsg) UnPack(data []byte) error {
	cdm.CmdSyn = &bpop.CommandSyn{}
	cdm.CmdSyn.Cmd = &bpop.CmdDownload{}

	if err := json.Unmarshal(data, cdm.CmdSyn); err != nil {
		return err
	}

	cdm.CmdDownload = cdm.CmdSyn.Cmd.(*bpop.CmdDownload)

	return nil
}

func (cdm *CommandDownloadMsg) Verify() bool {
	if bytes.Compare(cdm.Sn, cdm.CmdSyn.SN[:]) != 0 {
		return false
	}

	addr, _ := resolver.NewEthResolver(true).BMailBCA(cdm.CmdDownload.MailAddr)
	if addr != cdm.CmdDownload.Owner {
		return false
	}

	if !bmailcrypt.Verify(addr.ToPubKey(), cdm.CmdSyn.SN[:], cdm.CmdSyn.Sig) {
		return false
	}

	return true

}

func (cdm *CommandDownloadMsg) SetCurrentSn(sn []byte) {
	cdm.Sn = sn
}

func (cdm *CommandDownloadMsg) Dispatch() error {

	return nil
}

func RecoverFromFile(eid uuid.UUID) (cep *bmp.CryptEnvelope, err error) {
	data, err := savefile.ReadFromFile(eid)

	if err != nil {
		return nil, err
	}

	cep = &bmp.CryptEnvelope{}
	err = json.Unmarshal(data, cep)
	if err != nil {
		return nil, err
	}
	return

}

func (cdm *CommandDownloadMsg) Response() (WBody, error) {

	cdm.CmdAck = &bpop.CommandAck{}
	cdm.CmdDownAck = &bpop.CmdDownloadAck{}
	cdm.CmdAck.CmdCxt = cdm.CmdDownAck

	pmdb := bmaildb.GetBMPullMailDb()

	copy(cdm.CmdAck.NextSN[:], tools.NewSn(tools.SerialNumberLength))

	sm, err := pmdb.Find(cdm.CmdDownload.Owner.String())
	if err != nil {
		cdm.CmdAck.ErrorCode = bpop.EC_No_Mail
		return cdm.CmdAck, nil
	}

	cnt := cdm.CmdDownload.MailCnt
	if cnt <= 0 {
		cnt = bpop.DefaultMailCount
	}

	total := 0

	for i := len(sm.Smi) - 1; i >= 0; i-- {
		if sm.Smi[i].CreateTime < cdm.CmdDownload.BeforeTime {
			cep, err := RecoverFromFile(sm.Smi[i].Eid)
			if err != nil {
				continue
			}

			cdm.CmdDownAck.CryptEps = append(cdm.CmdDownAck.CryptEps, *cep)
			total++
			if total >= cnt {
				break
			}
		}
	}

	if len(cdm.CmdDownAck.CryptEps) == 0 {
		cdm.CmdAck.ErrorCode = bpop.EC_No_Mail
		return cdm.CmdAck, nil
	}

	cdm.CmdAck.Hash = cdm.CmdAck.CmdCxt.Hash()
	cdm.CmdAck.Sig = wallet.GetServerWallet().Sign(cdm.CmdAck.Hash)

	return cdm.CmdAck, nil

}
