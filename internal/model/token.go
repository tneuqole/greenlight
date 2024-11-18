package model

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/tneuqole/greenlight/internal/validator"
)

const ScopeActivation = "activation"

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:] // convert array to slice with [:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, plaintext string) {
	v.Check(plaintext != "", "token", "must be provided")
	v.Check(len(plaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(t *Token) error {
	q := `INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), Q_TIMEOUT)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, q, t.Hash, t.UserID, t.Expiry, t.Scope)
	return err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	q := `DELETE FROM tokens
	WHERE scope=$1 AND user_id=$2`

	ctx, cancel := context.WithTimeout(context.Background(), Q_TIMEOUT)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, q, scope, userID)
	return err
}
