package resourcetoken

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Sign(secret []byte, userID string, ttl time.Duration) string {
	expiry := strconv.FormatInt(time.Now().Add(ttl).Unix(), 10)
	payload := userID + ":" + expiry
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(payload + ":" + sig))
}

func Verify(secret []byte, token string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("invalid token encoding")
	}

	parts := strings.SplitN(string(raw), ":", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	userID, expiryStr, sigB64 := parts[0], parts[1], parts[2]

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid token expiry")
	}
	if time.Now().Unix() > expiry {
		return "", fmt.Errorf("token expired")
	}

	payload := userID + ":" + expiryStr
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sigB64), []byte(expectedSig)) {
		return "", fmt.Errorf("invalid token signature")
	}

	return userID, nil
}

func DeriveSecret(encryptionKey string) []byte {
	h := sha256.Sum256([]byte("courrier-resource-token:" + encryptionKey))
	return h[:]
}
