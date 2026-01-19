package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"net/url"
)

func TgGetHmacSha256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	hash := h.Sum(nil)
	return hash
}

func TgCheck(initData, loginToken string) (tgDate url.Values, res bool) {
	tgDate, err := UrlParamParse(initData)
	if err != nil {
		log.Error("utils.TgParseData", zap.Error(err), zap.String("initData", initData))
		return
	}
	dataCheckString, err := UrlParamSort(tgDate, nil, false, "hash")
	if err != nil {
		log.Error("Tg Data error", zap.Error(err))
		return
	}

	botTokenData := "WebAppData"
	secret := TgGetHmacSha256([]byte(botTokenData), []byte(loginToken))
	hash := hex.EncodeToString(TgGetHmacSha256([]byte(secret), []byte(dataCheckString)))
	if hash != tgDate.Get("hash") {
		log.Error("hash not true", zap.String("makeHash", hash), zap.String("hash", tgDate.Get("hash")))
		return
	}

	res = true

	return
}
