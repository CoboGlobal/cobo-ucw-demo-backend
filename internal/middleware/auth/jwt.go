package auth

import (
	"context"
	"time"

	"cobo-ucw-backend/internal/conf"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/ucw.v1.UserControlWallet/Ping"] = struct{}{}
	whiteList["/ucw.v1.UserControlWallet/Login"] = struct{}{}
	whiteList["/ucw.v1.UserControlWallet/TransactionWebhook"] = struct{}{}
	whiteList["/ucw.v1.UserControlWallet/CoboCallback"] = struct{}{}
	whiteList["/ucw.v1.UserControlWallet/TssRequestWebhook"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u UserInfo) GetUserId() string {
	return u.ID
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (user UserInfo, ok bool) {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return UserInfo{}, ok
	}
	customClaims, ok := claims.(*CustomClaims)
	if !ok {
		return UserInfo{}, ok
	}
	return UserInfo{ID: customClaims.UserID}, ok
}

const TokenExpireDuration = time.Minute * 20

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwtv4.RegisteredClaims
}

type JwtMiddleware struct {
	ca *conf.UCW_Auth
}

func NewJwtMiddleware(ca *conf.UCW_Auth) *JwtMiddleware {
	return &JwtMiddleware{ca: ca}
}
func (j *JwtMiddleware) GenRegisteredClaims(userID string) (string, error) {
	claims := CustomClaims{
		userID,
		jwtv4.RegisteredClaims{
			ExpiresAt: jwtv4.NewNumericDate(time.Now().Add(j.ca.GetExpireDuration().AsDuration())),
			Issuer:    "cobo-tss",
		},
	}

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.ca.GetApiKey()))

	return tokenString, err
}

func (j *JwtMiddleware) Server() middleware.Middleware {
	return selector.Server(
		jwt.Server(
			func(token *jwtv4.Token) (interface{}, error) {
				return []byte(j.ca.GetApiKey()), nil
			}, jwt.WithSigningMethod(jwtv4.SigningMethodHS256), jwt.WithClaims(func() jwtv4.Claims {
				return &CustomClaims{}
			}),
		),
	).Match(NewWhiteListMatcher()).Build()
}
