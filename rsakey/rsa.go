package rsakey

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools"
	"io/ioutil"
	"os"
	"path"
)

func GenerateKeyPair(bitsCnt int) (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := rsa.GenerateKey(rand.Reader, bitsCnt)

	if err != nil {
		return nil, nil
	}

	return priv, &priv.PublicKey

}

func KeyGenerated(ph string) bool {
	privFileName := path.Join(ph, "priv.key")

	return tools.FileExists(privFileName)
}

func Save2File(savePath string, privKey *rsa.PrivateKey, password string) error {
	if savePath == "" {
		return errors.New("Path is none")
	}

	if privKey == nil {
		return errors.New("Private key is none")
	}

	keypath := savePath
	if !tools.FileExists(keypath) {
		os.MkdirAll(keypath, 0755)
	}

	//save private key
	pb := x509.MarshalPKCS1PrivateKey(privKey)

	//aes encrypt
	pembyte := AesEncrypt(pb, []byte(password))

	block := &pem.Block{Type: "Priv", Bytes: pembyte}

	if f, err := os.OpenFile(path.Join(keypath, "priv.key"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755); err != nil {
		return err
	} else {
		if err = pem.Encode(f, block); err != nil {
			return err
		}
	}

	//save public key
	pubKey := &privKey.PublicKey
	pubbytes := x509.MarshalPKCS1PublicKey(pubKey)
	block = &pem.Block{Type: "Pub", Bytes: pubbytes}
	if f, err := os.OpenFile(path.Join(keypath, "pub.key"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755); err != nil {
		return err
	} else {
		if err = pem.Encode(f, block); err != nil {
			return err
		}
	}

	return nil

}

func LoadRSAKey(savePath string, passwd []byte) (priv *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	//load privKey
	var block *pem.Block
	if privKey, err := ioutil.ReadFile(path.Join(savePath, "priv.key")); err != nil {
		return nil, nil, errors.New("read priv.key error")
	} else {
		block, _ = pem.Decode(privKey)
		if block == nil {
			return nil, nil, errors.New("recover privKey error")
		}
	}

	pb := AesDecrypt(block.Bytes, passwd)
	if pb == nil {
		return nil, nil, errors.New("password error")
	}

	if priv, err = x509.ParsePKCS1PrivateKey(pb); err != nil {
		return nil, nil, errors.New("Parse privKey error")
	}

	pub = &priv.PublicKey

	//if pubkey,err:=ioutil.ReadFile(path.Join(savePath,"pub.key"));err!=nil{
	//	return nil,nil,errors.New("read pub.key error")
	//}else {
	//	block, _ = pem.Decode(pubkey)
	//	if block == nil{
	//		return nil,nil,errors.New("recover pubKey error")
	//	}
	//
	//}
	//
	//if pub,err = x509.ParsePKCS1PublicKey(block.Bytes);err!=nil{
	//	return nil,nil,errors.New("Parse PubKey error")
	//}

	return

}

func ParsePubKey(pubkey []byte) (pub *rsa.PublicKey, err error) {
	return x509.ParsePKCS1PublicKey(pubkey)
}

func PubKeyToBytes(pub *rsa.PublicKey) (pubkey []byte, err error) {
	if pub == nil {
		return nil, errors.New("Parameter error")
	}

	pubkey = x509.MarshalPKCS1PublicKey(pub)

	return

}

func EncryptRSA(data []byte, pub *rsa.PublicKey) (encData []byte, err error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func DecryptRsa(encData []byte, priv *rsa.PrivateKey) (data []byte, err error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, encData)
}

func SignRSA(data []byte, priv *rsa.PrivateKey) (signData []byte, err error) {
	hash256 := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash256[:])
}

func VerifyRSA(data []byte, sig []byte, pub *rsa.PublicKey) error {
	hash256 := sha256.Sum256(data)

	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash256[:], sig)

}

func PubKeyBytes2Addr(pubk []byte) string {
	hash256 := sha256.Sum256(pubk)

	return "58" + base58.Encode(hash256[:])
}

func PubKey2Addr(pk *rsa.PublicKey) string {
	if pk == nil {
		return ""
	}
	pkb, _ := PubKeyToBytes(pk)

	return PubKeyBytes2Addr(pkb)
}
