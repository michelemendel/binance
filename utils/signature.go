package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
)

func Signature(secretKey, query, body string, timestamp int64) string {
	payload := query + body + fmt.Sprintf("&timestamp=%d", timestamp)
	signed := HMACSign(secretKey, payload)
	urlEncoded := urlEncode(signed)
	return urlEncoded
}

// todo: Need this, do we?
// func base64Encode(str string) string {
// 	return base64.StdEncoding.EncodeToString([]byte(str))
// }

func HMACSign(secret, data string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func urlEncode(str string) string {
	return url.QueryEscape(str)
}
