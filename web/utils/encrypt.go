package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base32"
	"encoding/pem"
	"os"
	"path"
	"strings"
	"time"
	"yulong-hids/web/settings"
)

// GetPassword create password like "google auth"
func GetPassword(key string) uint32 {
	key = strings.Replace(key, " ", "", -1)
	key = strings.ToUpper(key)
	keyDecode, err := base32.StdEncoding.DecodeString(key)
	if err == nil {
		epochSeconds := time.Now().Unix()
		return oneTimePassword(keyDecode, toBytes(epochSeconds/30))
	}
	return 0
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

//GenRsaKey create key file
func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privateFile, err := os.Create(path.Join(settings.FilePath, settings.PrivateKeyName))
	defer privateFile.Close()
	if err != nil {
		return err
	}
	err = pem.Encode(privateFile, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	publicFile, err := os.Create(path.Join(settings.FilePath, settings.PublicKeyName))
	defer publicFile.Close()
	if err != nil {
		return err
	}
	err = pem.Encode(publicFile, block)
	if err != nil {
		return err
	}
	return nil
}
