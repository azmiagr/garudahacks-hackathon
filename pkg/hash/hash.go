package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
)

type Interface interface {
	HashNIK(nik string) string
}

type hasher struct {
	secretKey string
}

func Init() Interface {
	secretKey := os.Getenv("NIK_HASH_SECRET")
	if secretKey == "" {
		log.Fatalf("error init hash: NIK_HASH_SECRET is not set")
	}

	return &hasher{secretKey: secretKey}
}

func (h *hasher) HashNIK(nik string) string {
	mac := hmac.New(sha256.New, []byte(h.secretKey))
	mac.Write([]byte(nik))
	return hex.EncodeToString(mac.Sum(nil))
}
