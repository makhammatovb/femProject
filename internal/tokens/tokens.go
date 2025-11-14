package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	PlainText string `json:"token"`
	Hash      []byte `json:"-"`
	UserID    int64  `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string `json:"-"`
}

const (
	ScopeAuthentication = "authentication"
)

func GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope: scope,
	}

	emptryBytes := make([]byte, 32)
	_, err := rand.Read(emptryBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptryBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}
