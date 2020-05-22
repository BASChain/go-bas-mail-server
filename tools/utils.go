package tools

import "crypto/rand"

const (
	SerialNumberLength int = 16
)

func NewSn(bytesCnt int) []byte {
	sn := make([]byte, bytesCnt)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	return sn
}
