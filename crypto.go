package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// HEX

func Hex2Bin(hexStr string) (bytes []byte, err error) {
	bytes, err = hex.DecodeString(hexStr)
	return
}

func EncodeHex(bytes []byte) string {
	return fmt.Sprintf("%x", bytes)
}

// BASE64

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		err := errors.New("failed to base64 decode payload")
		return nil, err
	}
	return data, nil
}

// AES

func GenerateAesSecret() (hexKey string, err error) {
	keySize := 32
	key := make([]byte, keySize)

	io.ReadFull(rand.Reader, key)
	hexKey = EncodeHex(key)
	return
}

func EncryptAes(hexKey string, text string) string {
	key, err := Hex2Bin(hexKey)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(text))
	iv := make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext, []byte(text))

	var ivHex = EncodeHex(iv)
	var ciphertextBase64 = EncodeBase64(ciphertext)

	return ivHex + " " + ciphertextBase64
}

func DecryptAes(hexKey string, raw string) (string, error) {
	key, err := Hex2Bin(hexKey)
	if err != nil {
		panic(err)
	}

	words := strings.Fields(raw)

	if len(words) == 0 {
		panic("failed to decrypt aes due to invalid raw format (zero words)")
	}

	ivHex := words[0]
	iv, _ := Hex2Bin(ivHex)

	var ciphertext string
	if len(words) == 1 {
		ciphertext = ""
	} else if len(words) == 2 {
		ciphertext = words[1]
	} else {
		panic("failed to decrypt aes due to invalid raw format (more than two words)")
	}

	text, err := DecodeBase64(ciphertext)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)

	return string(text), nil
}

// RSA

const RSAKeySize = 3072

var defaultLabel = []byte{}

func GeneratePemToFile(filename string) (err error) {
	prvPem, pubPem, err := generatePemKeys()
	if err != nil {
		return
	}
	err = exportPem(filename, prvPem, pubPem)
	return
}

func ImportPemFromFile(filename string) (prv *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	prv, _, err = importKeyFromFile(filename)
	if err != nil {
		return
	}

	_, pub, err = importKeyFromFile(filename + ".pub")
	if err != nil {
		return
	}

	return
}

func DecryptRsa(prv *rsa.PrivateKey, ct []byte) (pt []byte, err error) {
	hash := sha256.New()
	pt, err = rsa.DecryptOAEP(hash, rand.Reader, prv, ct, defaultLabel)
	return
}

func EncryptRsa(pub *rsa.PublicKey, pt []byte) (ct []byte, err error) {
	if len(ct) > maxMessageLength(pub) {
		err = fmt.Errorf("message is too long")
		return
	}

	hash := sha256.New()
	ct, err = rsa.EncryptOAEP(hash, rand.Reader, pub, pt, defaultLabel)
	return
}

func Sign(prv *rsa.PrivateKey, msg string) (signature string, err error) {
	h := sha256.New()
	h.Write([]byte(msg))
	d := h.Sum(nil)
	sigBin, _ := rsa.SignPSS(rand.Reader, prv, crypto.SHA256, d, nil)
	signature = EncodeBase64(sigBin)
	return
}

func Verify(pub *rsa.PublicKey, msg string, signature string) (err error) {
	sig, _ := DecodeBase64(signature)
	h := sha256.New()
	h.Write([]byte(msg))
	d := h.Sum(nil)
	return rsa.VerifyPSS(pub, crypto.SHA256, d, sig, nil)
}

// private stuff

func importKeyFromFile(filename string) (prv *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	cert, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	for {
		var blk *pem.Block
		blk, cert = pem.Decode(cert)
		if blk == nil {
			break
		}
		switch blk.Type {
		case "RSA PRIVATE KEY":
			prv, err = x509.ParsePKCS1PrivateKey(blk.Bytes)
			return
		case "RSA PUBLIC KEY":
			var in interface{}
			in, err = x509.ParsePKIXPublicKey(blk.Bytes)
			if err != nil {
				return
			}
			pub = in.(*rsa.PublicKey)
			return
		}
		if cert == nil || len(cert) == 0 {
			break
		}
	}
	return
}

func generatePrivatePem(prv *rsa.PrivateKey) (prvPem string, err error) {
	cert := x509.MarshalPKCS1PrivateKey(prv)
	blk := new(pem.Block)
	blk.Type = "RSA PRIVATE KEY"
	blk.Bytes = cert

	var b bytes.Buffer
	err = pem.Encode(&b, blk)
	if err != nil {
		return
	}

	prvPem = b.String()
	return
}

func generatePublicPem(pub *rsa.PublicKey) (pubPem string, err error) {
	cert, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return
	}

	blk := new(pem.Block)
	blk.Type = "RSA PUBLIC KEY"
	blk.Bytes = cert

	var b bytes.Buffer
	err = pem.Encode(&b, blk)
	if err != nil {
		return
	}

	pubPem = b.String()
	return
}

func generatePemKeys() (prvPem string, pubPem string, err error) {
	key, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		return
	}

	prvPem, err = generatePrivatePem(key)
	if err != nil {
		return
	}

	pubPem, err = generatePublicPem(&key.PublicKey)
	if err != nil {
		return
	}

	return
}

func importKeyFromString(str string) (prv *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	cert := []byte(str)
	for {
		var blk *pem.Block
		blk, cert = pem.Decode(cert)
		if blk == nil {
			break
		}
		switch blk.Type {
		case "RSA PRIVATE KEY":
			prv, err = x509.ParsePKCS1PrivateKey(blk.Bytes)
			return
		case "RSA PUBLIC KEY":
			var in interface{}
			in, err = x509.ParsePKIXPublicKey(blk.Bytes)
			if err != nil {
				return
			}
			pub = in.(*rsa.PublicKey)
			return
		}
		if cert == nil || len(cert) == 0 {
			break
		}
	}
	return
}

func exportPem(filename string, prvPem string, pubPem string) (err error) {
	privateKeyFile, err := os.Create(filename)
	if err != nil {
		return
	}

	privateKeyFile.WriteString(prvPem)
	privateKeyFile.Sync()

	publicKeyFile, err := os.Create(filename + ".pub")
	if err != nil {
		return
	}

	publicKeyFile.WriteString(pubPem)
	publicKeyFile.Sync()
	return
}

func maxMessageLength(key *rsa.PublicKey) int {
	if key == nil {
		return 0
	}
	return (key.N.BitLen() / 8) - (2 * sha256.Size) - 2
}
