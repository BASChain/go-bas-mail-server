package savefile

import (
	"github.com/google/uuid"

	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"github.com/realbmail/go-bas-mail-server/config"
	"os"
	"path"
)

func DeriveFilePath(eid uuid.UUID) (filePath string, fileName string) {

	s := sha256.Sum256(eid[:])

	s58 := base58.Encode(s[:])

	tlen := len(s58)

	cfg := config.GetBMSCfg()

	filePath = path.Join(cfg.GetMailStorePath(), s58[tlen-8:tlen-6], s58[tlen-6:tlen-4], s58[tlen-4:tlen-2], s58[tlen-2:])

	fileName = path.Join(filePath, eid.String())

	return
}

func Save2File(eid uuid.UUID, data []byte) error {

	if data == nil {
		return errors.New("param error")
	}

	filePath, fileName := DeriveFilePath(eid)

	if !tools.FileExists(filePath) {
		os.MkdirAll(filePath, 0755)
	}

	if tools.FileExists(fileName) {
		return errors.New("file exists")
	}

	return tools.Save2File(data, fileName)

}

func ReadFromFile(eid uuid.UUID) (data []byte, err error) {
	_, fileName := DeriveFilePath(eid)

	if !tools.FileExists(fileName) {
		return nil, errors.New("file not exists")
	}

	return tools.OpenAndReadAll(fileName)
}
