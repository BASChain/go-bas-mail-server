package rsakey

import (
	"bytes"
	"crypto/aes"

	"crypto/cipher"
)

//func main() {
//	orig := "hello world"
//	key := "123456781234567812345678"
//
//
//	encryptCode := AesEncrypt(orig, key)
//	fmt.Println("Encrypt :",encryptCode)
//
//	decryptCode := AesDecrypt(encryptCode, key)
//	fmt.Println("Decrypt :",decryptCode)
//}

func appendToBlockSize(k []byte) []byte {
	l := len(k)

	for i := 0; i < aes.BlockSize-l; i++ {
		k = append(k, '0')
	}

	return k
}

func AesEncrypt(orig []byte, key []byte) []byte {

	k := appendToBlockSize(key)

	block, _ := aes.NewCipher(k)

	blockSize := block.BlockSize()

	origData := PKCS7Padding(orig, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])

	cryted := make([]byte, len(origData))

	blockMode.CryptBlocks(cryted, origData)

	return cryted

}

func AesDecrypt(cryted []byte, key []byte) []byte {

	k := appendToBlockSize(key)

	block, _ := aes.NewCipher(k)

	blockSize := block.BlockSize()

	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])

	orig := make([]byte, len(cryted))

	blockMode.CryptBlocks(orig, cryted)

	orig = PKCS7UnPadding(orig)
	return orig
}

func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
