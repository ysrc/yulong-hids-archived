package task

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"math/big"
	"yulong-hids/daemon/common"
)

func setPublicKey() {
	var res map[string]string
	url := common.Proto + "://" + common.ServerIP + common.PUBLICKEY_API
	resp, err := common.HTTPClient.Get(url)
	if err != nil {
		log.Println("HTTP get publickey error:", err.Error())
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Update publickey ioutil.ReadAll error", err.Error())
		return
	}
	json.Unmarshal([]byte(result), &res)
	common.PublicKey = res["public"]
}

func pubKeyDecrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil, errors.New("data length error")
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil, errors.New("message too long for RSA public key size")
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, errors.New("data broken, first byte is not zero")
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, errors.New("data is not encrypted by the private key")
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}
func rsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(common.PublicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKeyDecrypt(pub.(*rsa.PublicKey), ciphertext)
}
