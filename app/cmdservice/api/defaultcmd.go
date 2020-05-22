package api

import (
	"context"

	"time"

	"encoding/json"

	"github.com/BASChain/go-bas-mail-server/config"

	"github.com/BASChain/go-bas-mail-server/app/cmdcommon"
	"github.com/BASChain/go-bas-mail-server/app/cmdpb"
	"github.com/BASChain/go-bas-mail-server/bmtpserver"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"strconv"
)

type CmdDefaultServer struct {
	Stop func()
}

func (cds *CmdDefaultServer) DefaultCmdDo(ctx context.Context,
	request *cmdpb.DefaultRequest) (*cmdpb.DefaultResp, error) {
	if request.Reqid == cmdcommon.CMD_STOP {
		return cds.stop()
	}

	if request.Reqid == cmdcommon.CMD_CONFIG_SHOW {
		return cds.configShow()
	}

	if request.Reqid == cmdcommon.CMD_PK_SHOW {
		return cds.showAccout()
	}

	if request.Reqid == cmdcommon.CMD_RUN {
		return cds.serverRun()
	}

	resp := &cmdpb.DefaultResp{}

	resp.Message = "no cmd found"

	return resp, nil
}

func (cds *CmdDefaultServer) stop() (*cmdpb.DefaultResp, error) {

	go func() {
		time.Sleep(time.Second * 2)
		cds.Stop()
	}()
	resp := &cmdpb.DefaultResp{}
	resp.Message = "server stoped"
	return resp, nil
}

func encapResp(msg string) *cmdpb.DefaultResp {
	resp := &cmdpb.DefaultResp{}
	resp.Message = msg

	return resp
}

func (cds *CmdDefaultServer) configShow() (*cmdpb.DefaultResp, error) {
	cfg := config.GetBMSCfg()

	bapc, err := json.MarshalIndent(*cfg, "", "\t")
	if err != nil {
		return encapResp("Internal error"), nil
	}

	return encapResp(string(bapc)), nil
}

func (cds *CmdDefaultServer) showAccout() (*cmdpb.DefaultResp, error) {
	cfg := config.GetBMSCfg()

	return encapResp("Account: " + bmail.ToAddress(cfg.PubKey).String()), nil
}

func (cds *CmdDefaultServer) serverRun() (*cmdpb.DefaultResp, error) {

	if config.GetBMSCfg().PubKey == nil || config.GetBMSCfg().PrivKey == nil {
		return encapResp("bmtp need account"), nil
	}

	go bmtpserver.GetBMTPServer().StartTCPServer()

	return encapResp("bmtp server start at: " + strconv.Itoa(int(translayer.BMTP_PORT))), nil

}
