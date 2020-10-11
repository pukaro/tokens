package tokens

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/ed25519"
)

type TokenStatus struct {
	Valid bool
	Alive bool
}

// Проверяет access token на изменения, параметр userId является не обезательным.
func CheckAccessToken(accessToken string, userID int64) (status TokenStatus, err error) {
	payload, err := GetPayloadAccess(accessToken)
	if err != nil {
		return TokenStatus{}, err
	}
	if payload.UserID != userID && userID != 0 {
		return TokenStatus{}, errors.New("id пользователя был изменен")
	}

	newToken, err := payload.GenerateAccess()
	if err != nil {
		return TokenStatus{}, err
	}
	if newToken != accessToken {
		return TokenStatus{}, errors.New("access token был изменен")
	}
	if payload.Exp < time.Now().Unix() {
		return TokenStatus{Valid: true}, err
	}
	return TokenStatus{true, true}, err
}

// Проверяет только refresh token на измения.
func CheckRefreshToken(refreshToken, accessToken string) (status TokenStatus, err error) {
	var payload PayloadRefresh
	payloadBase64 := regexp.MustCompile("^.[^.]*").FindString(refreshToken)
	payloadJSON, err := base64.StdEncoding.DecodeString(payloadBase64)
	if err != nil {
		return
	}
	if err = json.Unmarshal(payloadJSON, &payload); err != nil {
		return
	}

	signAccess, err := base64.StdEncoding.DecodeString(payload.Access)
	if err != nil {
		return
	}
	if !ed25519.Verify(PublicKey, []byte(accessToken), signAccess) {
		return TokenStatus{}, errors.New("данные access токена в refresh токене были изменены")
	}

	signRefresh := ed25519.Sign(PrivateKey, payloadJSON)
	newToken := payloadBase64 + "." + base64.StdEncoding.EncodeToString(signRefresh)
	if refreshToken != newToken {
		return TokenStatus{}, errors.New("токен был изменен")
	}
	if payload.Exp < time.Now().Unix() {
		return TokenStatus{Valid: true}, err
	}
	return TokenStatus{true, true}, err
}

// Проверяет связанную пару Access, Refresh tokens.
// Если Access token неправельный, то сразу возвращается ошибка.
func CheckTokens(access, refrash string, userID int64) (statusAccess, statusRefresh TokenStatus, err error) {
	statusAccess, err = CheckAccessToken(access, userID)
	if err != nil || !statusAccess.Valid {
		return
	}
	statusRefresh, err = CheckRefreshToken(refrash, access)
	return
}
