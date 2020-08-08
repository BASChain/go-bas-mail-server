package api

import (
	"context"

	"time"

	"encoding/json"

	"github.com/realbmail/go-bas-mail-server/config"

	"github.com/realbmail/go-bas-mail-server/app/cmdcommon"
	"github.com/realbmail/go-bas-mail-server/app/cmdpb"
	"github.com/realbmail/go-bas-mail-server/bmtpserver"
	"github.com/realbmail/go-bas-mail-server/bpopserver"
	"github.com/realbmail/go-bmail-account"
	"github.com/realbmail/go-bmail-protocol/translayer"
	"strconv"
	"sync"
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

var (
	runingFlag      bool
	runningOnceLock sync.Mutex
)

func (cds *CmdDefaultServer) serverRun() (*cmdpb.DefaultResp, error) {

	if config.GetBMSCfg().PubKey == nil || config.GetBMSCfg().PrivKey == nil {
		return encapResp("bmtp need account"), nil
	}

	if !runingFlag {
		runningOnceLock.Lock()
		defer runningOnceLock.Unlock()
		if !runingFlag {
			go bmtpserver.GetBMTPServer().StartTCPServer()
			go bpopserver.GetBMTPServer().StartTCPServer()
		}

		runingFlag = true
	}

	msg := "bmtp server start at: " + strconv.Itoa(int(translayer.BMTP_PORT))
	msg += "\r\nbpop server start at: " + strconv.Itoa(int(translayer.BPOP3))

	return encapResp(msg), nil

}
