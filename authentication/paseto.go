package authentication

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

var (
	errExpiredToken = errors.New("token has expired")
)

type PasetoAuthenticationPayload struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	PositionID     string    `json:"position_id"`
	BranchOfficeID string    `json:"branch_office_id"`
	AccountType    string    `json:"account_type"`
	IssuedAt       int64     `json:"issued_at"`
	ExpiredAt      int64     `json:"expired_at"`
}

type PasetoAuthentication interface {
	CreateToken(payload *PasetoAuthenticationPayload) (string, error)
	VerifyToken(token string) (*PasetoAuthenticationPayload, error)
}

type PasetoAuthenticationCtx struct {
	paseto       *paseto.V2
	symmetricKey []byte
	privateKey   ed25519.PrivateKey
	publicKey    ed25519.PublicKey
	mode         string
}

func NewPasetoAuthentication(key, mode string) PasetoAuthentication {
	if len(key) != chacha20poly1305.KeySize {
		log.Panic(fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize))
	}

	b, _ := hex.DecodeString("b4cbfb43df4ce210727d953e4a713307fa19bb7d9f85041438d9e11b942a37741eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	privateKey := ed25519.PrivateKey(b)

	b, _ = hex.DecodeString("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	publicKey := ed25519.PublicKey(b)

	auth := &PasetoAuthenticationCtx{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
		privateKey:   privateKey,
		publicKey:    publicKey,
		mode:         mode,
	}

	return auth
}

// CreateToken create new token
func (auth *PasetoAuthenticationCtx) CreateToken(payload *PasetoAuthenticationPayload) (string, error) {
	timeAsiaJakarta, _ := time.LoadLocation("Asia/Jakarta")
	start := time.Now().In(timeAsiaJakarta).UTC()
	end := start.Add(time.Hour * time.Duration(24))

	payload.IssuedAt = start.Unix()
	payload.ExpiredAt = end.Unix()
	if strings.EqualFold(auth.mode, "production") {
		return auth.paseto.Sign(auth.privateKey, &payload, &payload)
	}

	return auth.paseto.Encrypt(auth.symmetricKey, &payload, payload)
}

// VerifyToken will verify token payload
func (auth *PasetoAuthenticationCtx) VerifyToken(token string) (*PasetoAuthenticationPayload, error) {
	payload := &PasetoAuthenticationPayload{}

	if strings.EqualFold(auth.mode, "production") {
		err := auth.paseto.Verify(token, auth.publicKey, payload, nil)
		if err != nil {
			return nil, err
		}
	} else {
		err := auth.paseto.Decrypt(token, auth.symmetricKey, payload, nil)
		if err != nil {
			return nil, err
		}
	}

	exp := time.Unix(payload.ExpiredAt, 0)
	if time.Now().After(exp) {
		return nil, errExpiredToken
	}

	return payload, nil
}
