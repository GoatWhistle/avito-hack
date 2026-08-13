package moreless

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

const photoTokenPrefix = "mlp"

type PhotoSigner struct {
	secret []byte
}

func NewPhotoSigner(secret string) PhotoSigner {
	return PhotoSigner{secret: []byte(secret)}
}

func (s PhotoSigner) Enabled() bool {
	return len(s.secret) > 0
}

func (s PhotoSigner) Sign(roundID, displayID string) string {
	if !s.Enabled() {
		return ""
	}

	body := roundID + ":" + displayID
	encoded := base64.RawURLEncoding.EncodeToString([]byte(body))

	return photoTokenPrefix + "." + encoded + "." + s.mac(encoded)
}

func (s PhotoSigner) Verify(token string) (string, string, bool) {
	if !s.Enabled() {
		return "", "", false
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] != photoTokenPrefix {
		return "", "", false
	}

	if !hmac.Equal([]byte(parts[2]), []byte(s.mac(parts[1]))) {
		return "", "", false
	}

	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	roundID, displayID, ok := strings.Cut(string(raw), ":")
	if !ok {
		return "", "", false
	}

	return roundID, displayID, true
}

func (s PhotoSigner) mac(payload string) string {
	sum := hmac.New(sha256.New, s.secret)
	sum.Write([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(sum.Sum(nil))
}
