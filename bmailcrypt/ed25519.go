package bmailcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-account"
	"github.com/BASChain/go-account/edwards25519"

	"github.com/kprc/nbsnetwork/tools"
	"golang.org/x/crypto/curve25519"
	"io"
	"log"

	"github.com/realbmail/go-bmail-account"

	"github.com/realbmail/go-bas-mail-server/config"
)

type KeyJson struct {
	PubKey    string `json:"pub_key"`
	CipherKey string `json:"cipher_key"`
}

func KeyIsGenerated() bool {
	cfg := config.GetBMSCfg()
	if cfg == nil {
		log.Fatal("Can't Get Config file")
		return false
	}

	if tools.FileExists(cfg.GetKeyPath()) {
		return true
	}

	return false
}

func LoadKey(password string) {
	cfg := config.GetBMSCfg()

	data, err := tools.OpenAndReadAll(cfg.GetKeyPath())
	if err != nil {
		log.Fatal("Load From key file error")
		return
	}

	kj := &KeyJson{}

	err = json.Unmarshal(data, kj)
	if err != nil {
		log.Fatal("Load From json error")
		return
	}

	pk := bmail.Address(kj.PubKey).ToPubKey()
	var priv ed25519.PrivateKey
	priv, err = account.DecryptSubPriKey(pk, kj.CipherKey, password)
	if err != nil {
		log.Fatal("Decrypt PrivKey failed")
		return
	}

	cfg.PubKey = pk
	cfg.PrivKey = priv

	return

}

func GenEd25519KeyAndSave(password string) error {

	var (
		priv ed25519.PrivateKey
		pub  ed25519.PublicKey
		err  error
	)
	cnt := 0
	for {
		cnt++
		pub, priv, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			if cnt > 10 {
				return err
			}
			continue
		} else {
			break
		}
	}

	var cipherTxt string
	cipherTxt, err = account.EncryptSubPriKey(priv, pub, password)
	if err != nil {
		return err
	}

	kj := &KeyJson{PubKey: bmail.ToAddress(pub[:]).String(), CipherKey: cipherTxt}

	cfg := config.GetBMSCfg()

	var data []byte
	data, err = json.Marshal(*kj)
	err = tools.Save2File(data, cfg.GetKeyPath())
	if err != nil {
		return err
	}

	return nil
}

func Sign(priv ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(priv, message)
}

func Verify(pub ed25519.PublicKey, message, sig []byte) bool {
	return ed25519.Verify(pub, message, sig)
}

func EncryptWithIV(key, iv, plainTxt []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainTxt))

	copy(cipherText[:aes.BlockSize], iv)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTxt)

	return cipherText, nil
}

func Encrypt(key []byte, plainTxt []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainTxt))

	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTxt)

	return cipherText, nil
}

func Decrypt(key []byte, cipherTxt []byte) ([]byte, error) {

	_, plainTxt, err := DecryptAndIV(key, cipherTxt)

	return plainTxt, err

}

func DecryptAndIV(key []byte, cipherTxt []byte) (iv, plainTxt []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	if len(cipherTxt) < aes.BlockSize {
		return nil, nil, fmt.Errorf("cipher text too short")
	}

	iv = cipherTxt[:aes.BlockSize]
	cipherTxt = cipherTxt[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherTxt, cipherTxt)

	return iv, cipherTxt, nil
}

func DeriveKey(seed []byte) (pub ed25519.PublicKey, priv ed25519.PrivateKey) {
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, privateKey[32:])

	return publicKey, privateKey
}

var KP = struct {
	S int
	N int
	R int
	P int
	L int
}{
	S: 8,
	N: 1 << 15,
	R: 8,
	P: 1,
	L: 32,
}
var (
	EConvertCurvePubKey = fmt.Errorf("convert ed25519 public key to curve25519 public key failed")
)

func GenerateAesKey(peerPub []byte, key ed25519.PrivateKey) ([]byte, error) {
	var priKey [32]byte
	var privateKeyBytes [64]byte
	copy(privateKeyBytes[:], key)
	PrivateKeyToCurve25519(&priKey, &privateKeyBytes)

	var curvePub, pubKey [32]byte
	copy(pubKey[:], peerPub)
	if ok := PublicKeyToCurve25519(&curvePub, &pubKey); !ok {
		return nil, EConvertCurvePubKey
	}
	return curve25519.X25519(priKey[:], curvePub[:])
}

func populateKey(data []byte) (ed25519.PublicKey, ed25519.PrivateKey) {
	pri := ed25519.PrivateKey(data)
	pub := pri.Public().(ed25519.PublicKey)
	return pub, pri
}

func PrivateKeyToCurve25519(curve25519Private *[32]byte, privateKey *[64]byte) {
	h := sha512.New()
	h.Write(privateKey[:32])
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(curve25519Private[:], digest)
}

func edwardsToMontgomeryX(outX, y *edwards25519.FieldElement) {
	var oneMinusY edwards25519.FieldElement
	edwards25519.FeOne(&oneMinusY)
	edwards25519.FeSub(&oneMinusY, &oneMinusY, y)
	edwards25519.FeInvert(&oneMinusY, &oneMinusY)

	edwards25519.FeOne(outX)
	edwards25519.FeAdd(outX, outX, y)

	edwards25519.FeMul(outX, outX, &oneMinusY)
}

func PublicKeyToCurve25519(curve25519Public *[32]byte, publicKey *[32]byte) bool {
	var A edwards25519.ExtendedGroupElement
	if !A.FromBytes(publicKey) {
		return false
	}

	var x edwards25519.FieldElement
	edwardsToMontgomeryX(&x, &A.Y)
	edwards25519.FeToBytes(curve25519Public, &x)
	return true
}
