package api

import (
	"context"
	"github.com/realbmail/go-bas-mail-server/app/cmdcommon"
	"github.com/realbmail/go-bas-mail-server/app/cmdpb"
	"github.com/realbmail/go-bas-mail-server/bmailcrypt"
	"github.com/realbmail/go-bas-mail-server/config"
	"github.com/realbmail/go-bmail-account"
)

type CmdStringOPSrv struct {
}

func (cso *CmdStringOPSrv) StringOpDo(cxt context.Context, so *cmdpb.StringOP) (*cmdpb.DefaultResp, error) {
	msg := ""
	switch so.Op {
	case cmdcommon.CMD_ACCOUNT_CREATE:
		msg = createAccount(so.Param)
	case cmdcommon.CMD_ACCOUNT_LOAD:
		msg = loadAccount(so.Param)
	//	msg = GetRecords(so.Param)
	//case cmdcommon.CMD_DEAL:
	//	msg = GetDeal(so.Param)
	//case cmdcommon.CMD_ORDER:
	//	msg = GetOrder(so.Param)
	default:
		return encapResp("Command Not Found"), nil
	}

	return encapResp(msg), nil
}

func createAccount(passwd string) string {
	err := bmailcrypt.GenEd25519KeyAndSave(passwd)
	if err != nil {
		return "create account failed"
	}

	bmailcrypt.LoadKey(passwd)

	addr := bmail.ToAddress(config.GetBMSCfg().PubKey).String()

	return "Address: " + addr
}

func loadAccount(passwd string) string {

	bmailcrypt.LoadKey(passwd)

	addr := bmail.ToAddress(config.GetBMSCfg().PubKey).String()

	return "load account success! \r\nAddress: " + addr
}
