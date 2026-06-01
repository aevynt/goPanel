package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"time"
)

// GenerateSecret generates a 16-character base32 secret
func GenerateSecret() (string, error) {
	bytes := make([]byte, 10)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes), nil
}

// VerifyTOTP checks if the 6-digit code is valid for the current time window, allowing clock drift of -1 to +1 intervals
func VerifyTOTP(secret string, code string) bool {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return false
	}

	currentTime := time.Now().Unix()
	interval := currentTime / 30

	// Check current, past, and next intervals
	for i := int64(-1); i <= 1; i++ {
		t := interval + i
		if fmt.Sprintf("%06d", calculateHOTP(key, t)) == code {
			return true
		}
	}
	return false
}

func calculateHOTP(key []byte, interval int64) int {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(interval))

	h := hmac.New(sha1.New, key)
	h.Write(buf)
	sum := h.Sum(nil)

	offset := sum[19] & 0x0f
	binCode := int(sum[offset]&0x7f)<<24 |
		int(sum[offset+1]&0xff)<<16 |
		int(sum[offset+2]&0xff)<<8 |
		int(sum[offset+3]&0xff)

	return binCode % 1000000
}
