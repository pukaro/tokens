package tokens

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
)

// Позволяет извлечь payload из Access token
func GetPayloadAccess(accessToken string) (payload PayloadAccess, err error) {
	payloadBase64 := regexp.MustCompile(`\..[^.]*`).FindString(accessToken)
	payloadJSON, err := base64.StdEncoding.DecodeString(payloadBase64[1:])
	if err != nil {
		return
	}
	err = json.Unmarshal(payloadJSON, &payload)
	return
}

// Позволяет извлечь payload из Refresh token
func GetPayliadRefresh(refreshToken string) (payload PayloadRefresh, err error) {
	payloadRefBase64 := regexp.MustCompile("^.[^.]*").FindString(refreshToken)
	payloadRefJSON, err := base64.StdEncoding.DecodeString(payloadRefBase64)
	if err != nil {
		return
	}
	err = json.Unmarshal(payloadRefJSON, &payload)
	return
}
