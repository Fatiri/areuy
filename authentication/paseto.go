package authentication

import (
	"crypto/ed25519"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Fatiri/areuy/exception"
	"github.com/aead/chacha20poly1305"
	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
)

type Role string

var AuthorizationHeaderKey = "Authorization"
var AuthorizationTypeBearer = "bearer"
var AuthorizationPayloadKey = "authorization_payload"

type PasetoAuthenticationGinPayload struct {
	ID        string `json:"id,omitempty"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiredAt int64  `json:"expired_at"`
}

type pasetoAuthenticationGinPayloadPublic struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiredAt int64  `json:"expired_at"`
}

type PasetoAuthenticationGin interface {
	CreateToken(payload *PasetoAuthenticationGinPayload, access string) (string, error)
	VerifyToken(token string) (*PasetoAuthenticationGinPayload, *exception.Response)
	PasetoGinMiddleware(roles []string) gin.HandlerFunc
}

type PasetoAuthenticationGinCtx struct {
	paseto       *paseto.V2
	SymmetricKey []byte
	PrivateKey   ed25519.PrivateKey
	PublicKey    ed25519.PublicKey
	Mode         string
}

func NewPasetoAuthenticationGin(ctx PasetoAuthenticationGinCtx) PasetoAuthenticationGin {
	if len(ctx.SymmetricKey) != chacha20poly1305.KeySize {
		log.Panic(fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize))
	}

	return &PasetoAuthenticationGinCtx{
		paseto:       paseto.NewV2(),
		SymmetricKey: ctx.SymmetricKey,
		PrivateKey:   ctx.PrivateKey,
		PublicKey:    ctx.PublicKey,
		Mode:         ctx.Mode,
	}
}

// CreateToken create new token
func (auth *PasetoAuthenticationGinCtx) CreateToken(payload *PasetoAuthenticationGinPayload, access string) (string, error) {
	var IPayload interface{}
	if strings.ToLower(access) == "public" {
		IPayload = pasetoAuthenticationGinPayloadPublic{
			Username:  payload.Username,
			Role:      payload.Role,
			IssuedAt:  payload.IssuedAt,
			ExpiredAt: payload.ExpiredAt,
		}
	} else {
		IPayload = payload
	}

	if strings.ToLower(auth.Mode) == "production" {
		return auth.paseto.Sign(auth.PrivateKey, &payload, IPayload)
	}

	return auth.paseto.Encrypt(auth.SymmetricKey, &payload, IPayload)
}

// VerifyToken will verify token payload
func (auth *PasetoAuthenticationGinCtx) VerifyToken(token string) (*PasetoAuthenticationGinPayload, *exception.Response) {
	payload := &PasetoAuthenticationGinPayload{}

	if strings.ToLower(auth.Mode) == "production" {
		err := auth.paseto.Verify(token, auth.PublicKey, payload, nil)
		if err != nil {
			return nil, exception.Error(nil, exception.Message{
				Id: "Token akses tidak valid",
				En: "Invalid authorization token",
			}, auth.Mode)
		}
	} else {
		err := auth.paseto.Decrypt(token, auth.SymmetricKey, payload, nil)
		if err != nil {
			return nil, exception.Error(nil, exception.Message{
				Id: "Token akses tidak valid",
				En: "Invalid authorization token",
			}, auth.Mode)
		}
	}

	exp := time.Unix(payload.ExpiredAt, 0)
	if time.Now().After(exp) {
		return nil, exception.Error(nil, exception.Message{
			Id: "Akses telah kedaluwarsa",
			En: "Access has expired",
		}, auth.Mode)
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
			}, auth.Mode))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.Error(nil, exception.Message{
				Id: "Authorization token tidak tersedia",
				En: "Authorization token is not provided",
			}, auth.Mode))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.Error(nil, exception.Message{
				Id: "Tipe Authorization tidak valid",
				En: "Authorization type is not valid",
			}, auth.Mode))
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
			}, auth.Mode))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
