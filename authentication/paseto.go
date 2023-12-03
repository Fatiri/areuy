package authentication

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Fatiri/areuy/exception"
	"github.com/aead/chacha20poly1305"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type Role string

var AuthorizationHeaderKey = "Authorization"
var AuthorizationTypeBearer = "bearer"
var AuthorizationPayloadKey = "authorization_payload"

type PasetoAuthenticationGinPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  int64     `json:"issued_at"`
	ExpiredAt int64     `json:"expired_at"`
}

type PasetoAuthenticationGin interface {
	CreateToken(payload *PasetoAuthenticationGinPayload) (string, error)
	VerifyToken(token string) (*PasetoAuthenticationGinPayload, *exception.Response)
	PasetoGinMiddleware(roles []string) gin.HandlerFunc
}

type PasetoAuthenticationGinCtx struct {
	paseto       *paseto.V2
	symmetricKey []byte
	privateKey   ed25519.PrivateKey
	publicKey    ed25519.PublicKey
	mode         string
	TokenExpired time.Duration
}

func NewPasetoAuthenticationGin(key, mode string) PasetoAuthenticationGin {
	if len(key) != chacha20poly1305.KeySize {
		log.Panic(fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize))
	}

	b, _ := hex.DecodeString(key)
	privateKey := ed25519.PrivateKey(b)

	b, _ = hex.DecodeString(key)
	publicKey := ed25519.PublicKey(b)

	auth := &PasetoAuthenticationGinCtx{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
		privateKey:   privateKey,
		publicKey:    publicKey,
		mode:         mode,
	}

	return auth
}

// CreateToken create new token
func (auth *PasetoAuthenticationGinCtx) CreateToken(payload *PasetoAuthenticationGinPayload) (string, error) {
	timeAsiaJakarta, _ := time.LoadLocation("Asia/Jakarta")
	start := time.Now().In(timeAsiaJakarta).UTC()
	end := start.Add(auth.TokenExpired)

	payload.IssuedAt = start.Unix()
	payload.ExpiredAt = end.Unix()
	if strings.EqualFold(auth.mode, "Production") {
		return auth.paseto.Sign(auth.privateKey, &payload, &payload)
	}

	return auth.paseto.Encrypt(auth.symmetricKey, &payload, payload)
}

// VerifyToken will verify token payload
func (auth *PasetoAuthenticationGinCtx) VerifyToken(token string) (*PasetoAuthenticationGinPayload, *exception.Response) {
	payload := &PasetoAuthenticationGinPayload{}

	if strings.EqualFold(auth.mode, "production") {
		err := auth.paseto.Verify(token, auth.publicKey, payload, nil)
		if err != nil {
			return nil, exception.Error(nil, exception.Message{
				Id: "Token akses tidak valid",
				En: "Invalid authorization token",
			}, auth.mode)
		}
	} else {
		err := auth.paseto.Decrypt(token, auth.symmetricKey, payload, nil)
		if err != nil {
			return nil, exception.Error(nil, exception.Message{
				Id: "Token akses tidak valid",
				En: "Invalid authorization token",
			}, auth.mode)
		}
	}

	exp := time.Unix(payload.ExpiredAt, 0)
	if time.Now().After(exp) {
		return nil, exception.Error(nil, exception.Message{
			Id: "Akses telah kedaluwarsa",
			En: "Access has expired",
		}, auth.mode)
	}

	return payload, nil
}

// AuthMiddleware creates a gin middleware for authorization
func (auth *PasetoAuthenticationGinCtx) PasetoGinMiddleware(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.Error(nil, exception.Message{
				Id: "Authorization header tidak tersedia",
				En: "Authorization header is not provided",
			}, auth.mode))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.Error(nil, exception.Message{
				Id: "Authorization token tidak tersedia",
				En: "Authorization token is not provided",
			}, auth.mode))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.Error(nil, exception.Message{
				Id: "Tipe Authorization tidak valid",
				En: "Authorization type is not valid",
			}, auth.mode))
			return
		}

		accessToken := fields[1]
		payload, err := auth.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		isRoleAuthorized := false
		for _, role := range roles {
			if role == payload.Role {
				isRoleAuthorized = true
				break
			}
		}

		if !isRoleAuthorized {
			ctx.AbortWithStatusJSON(http.StatusForbidden, exception.Error(nil, exception.Message{
				Id: fmt.Sprintf("Role %s akses di tolak", payload.Role),
				En: fmt.Sprintf("Role %s access denied", payload.Role),
			}, auth.mode))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
