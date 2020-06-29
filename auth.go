package goshopee

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func VerifyPushMsg(url, requestBody, pKey, authorization string) (result bool) {
	calAuth := MakeAuthToken(url, requestBody, pKey)
	return authorization == calAuth
}

func MakeAuthToken(url, requestBody, pKey string) string {
	baseStr := url + "|" + requestBody
	h := hmac.New(sha256.New, []byte(pKey))
	h.Write([]byte(baseStr))
	calAuth := fmt.Sprintf("%x", h.Sum(nil))
	return calAuth
}
