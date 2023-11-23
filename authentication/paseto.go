package authentication

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type Role string

var AuthorizationHeaderKey = "Authorization"
var AuthorizationTypeBearer = "bearer"
var AuthorizationPayloadKey = "authorization_payload"

// List of role
var RoleAdmin Role = "ADMIN"
var RoleUser Role = "USER"

var errExpiredToken = errors.New("token has expired")

type PasetoAuthenticationPayload struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	AccountType string    `json:"account_type"`
	IssuedAt    int64     `json:"issued_at"`
	ExpiredAt   int64     `json:"expired_at"`
}

type PasetoAuthentication interface {
	CreateToken(payload *PasetoAuthenticationPayload) (string, error)
	VerifyToken(token string) (*PasetoAuthenticationPayload, error)
	PasetoGinMiddleware(roles []Role, envApp string) gin.HandlerFunc
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

	b, _ := hex.DecodeString(key)
	privateKey := ed25519.PrivateKey(b)

	b, _ = hex.DecodeString(key)
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

type errorResponse struct {
	Code     float64 `json:"code"`
	Messsage string  `json:"message"`
}

// AuthMiddleware creates a gin middleware for authorization
func (auth *PasetoAuthenticationCtx) PasetoGinMiddleware(roles []Role, envApp string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{
				Code:     401,
				Messsage: "authorization header is not provided",
			})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{
				Code:     401,
				Messsage: "invalid authorization header format",
			})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{
				Code:     401,
				Messsage: "unsupported authorization type " + authorizationType,
			})
			return
		}

		accessToken := fields[1]
		payload, err := auth.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{
				Code:     401,
				Messsage: err.Error(),
			})
			return
		}

		isAuthorized := false
		for _, role := range roles {
			if string(role) == payload.AccountType {
				isAuthorized = true
			}
		}

		if !isAuthorized {
			err := errors.New("forbiden access")
			ctx.AbortWithStatusJSON(http.StatusForbidden, err)
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
