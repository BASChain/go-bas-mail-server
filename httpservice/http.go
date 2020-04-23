package httpservice

import (
	"context"
	"github.com/BASChain/go-bas-mail-server/config"
	"github.com/BASChain/go-bas-mail-server/httpservice/api"
	"log"
	"net/http"
	"strconv"
	"time"
)

var webserver *http.Server

func StartWebDaemon() {
	mux := http.NewServeMux()

	mux.Handle("/ajax/chg", api.NewWhiteList())

	addr := ":" + strconv.Itoa(config.GetBMSCfg().MgtHttpPort)

	log.Println("Web Server Start at", addr)

	webserver = &http.Server{Addr: addr, Handler: mux}

	log.Fatal(webserver.ListenAndServe())

}

func StopWebDaemon() {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	webserver.Shutdown(ctx)
}
