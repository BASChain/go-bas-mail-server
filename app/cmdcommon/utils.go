package cmdcommon

import (
	"errors"
	"github.com/realbmail/go-bas-mail-server/config"
	"github.com/kprc/nbsnetwork/tools"
)

func IsProcessCanStarted() (bool, error) {

	cfg := config.PreLoad()

	if cfg == nil {
		return true, nil
	}

	ip, port, err := tools.GetIPPort(cfg.CmdListenPort)
	if err != nil {

		return false, errors.New("Command line listen address error")
	}

	if tools.CheckPortUsed("tcp", ip, uint16(port)) {

		return false, errors.New("Process have started")
	}

	return true, nil
}

func IsProcessStarted() (bool, error) {
	if !config.IsInitialized() {
		return false, errors.New("need to initialize config file first")
	}

	cfg := config.PreLoad()
	if cfg == nil {
		return false, errors.New("load config failed")
	}

	ip, port, err := tools.GetIPPort(cfg.CmdListenPort)
	if err != nil {

		return false, errors.New("Command line listen address error")
	}

	if tools.CheckPortUsed("tcp", ip, uint16(port)) {
		return true, nil
	}

	return false, errors.New("process is not started")

}
