package cmdclient

import (
	"fmt"
	"github.com/realbmail/go-bas-mail-server/app/cmdcommon"
	"github.com/realbmail/go-bas-mail-server/app/cmdpb"
	"log"
	"strings"
)

func StringOpCmdSend(addr string, cmd int32, reqs string) {
	if addr == "" || strings.Contains(addr, "127.0.0.1") {
		if _, err := cmdcommon.IsProcessStarted(); err != nil {
			log.Println(err)
			return
		}
	}

	request := &cmdpb.StringOP{}
	request.Op = cmd
	request.Param = reqs

	cc := NewCmdClient(addr)

	cc.DialToCmdServer()
	defer cc.Close()

	client := cmdpb.NewStringopsrvClient(cc.conn.c)

	if resp, err := client.StringOpDo(cc.conn.ctx, request); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Message)
	}

}
