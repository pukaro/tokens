package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"crypto/sha512"

	"golang.org/x/crypto/ed25519"
)

const lifetimeAcc = time.Hour
const lifetimeRefresh = time.Hour * 24 * 90

var header = base64.StdEncoding.EncodeToString([]byte(`{"alg":"SHA512","typ":"JWT"}`))
var PrivateKey, PublicKey = func() (PrivKey ed25519.PrivateKey, PubKey ed25519.PublicKey) {
	PrivKey, _ = base64.StdEncoding.DecodeString(`S+KyUM/Ewk/N1AhsNlN2ff6ZHPQoGMPWA3A0Mww/ufYE+C5A1PIZv7orFIqFIaDzA9K/HqGFIBzy9mmF/0VZ5Q==`)
	PubKey = ed25519.PublicKey(PrivKey[32:])
	return
}()

type PayloadAccess struct {
	UserID int64  `json:"user_id"`
	Exp    int64  `json:"exp"`
	GUID   string `json:"GUID"`
}
type PayloadRefresh struct {
	Access string `json:"access"`
	Exp    int64  `json:"exp"`
	GUID   string `json:"GUID"`
}

func (data PayloadAccess) GenerateAccess() (string, error) {
	if data.GUID == "" {
		buff := make([]byte, 16)
		if _, err := rand.Read(buff); err != nil {
			return "", err
		}
		data.GUID = base64.StdEncoding.EncodeToString(buff)
	}
	if data.Exp == 0 {
		data.Exp = time.Now().Add(lifetimeAcc).Unix()
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	payload := base64.StdEncoding.EncodeToString(dataJSON)
	token := header + "." + payload
	hash := sha512.Sum512([]byte(token))
	sig := base64.StdEncoding.EncodeToString(hash[:])
	token += "." + sig
	return token, err
}

// Генерирует Refresh token на основе передаваемого Access token.
func GenerateRefresh(accessToken string) (string, error) {
	if accessToken == "" {
		return "", errors.New("необходимо добавить Access token")
	}

	buff := make([]byte, 16)
	if _, err := rand.Read(buff); err != nil {
		return "", err
	}

	// Вместо подпись можно брать хеш Access token
	signTokenAcc := ed25519.Sign(PrivateKey, []byte(accessToken))
	payload := PayloadRefresh{
		Access: base64.StdEncoding.EncodeToString(signTokenAcc),
		Exp:    time.Now().Add(lifetimeRefresh).Unix(),
		GUID:   base64.StdEncoding.EncodeToString(buff),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Вместо подпись payload, можно его хешировать
	singTokenRef := ed25519.Sign(PrivateKey, payloadJSON)
	refreshToken := base64.StdEncoding.EncodeToString(payloadJSON)
	refreshToken += "." + base64.StdEncoding.EncodeToString(singTokenRef)
	return refreshToken, nil
}
