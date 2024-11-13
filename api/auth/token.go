package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kacperhemperek/twitter-v2/api"
)

var (
	ErrMissingToken  = errors.New("bearer token is empty")
	ErrInvalidClaims = errors.New("unknown claims type cannot process")
)

type Claims struct {
	User *UserToken `json:"user,omitempty"`
	jwt.RegisteredClaims
}

type UserToken struct {
	ID    string `json:"id" mapstructure:"id"`
	Email string `json:"email" mapstructure:"email"`
	Name  string `json:"name" mapstructure:"name"`
}

func (t UserToken) SignAccessToken() (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		t.getClaims(time.Now().Add(time.Minute*15)),
	).SignedString(t.getSecret())
}

func (t UserToken) SignRefreshToken() (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		t.getClaims(time.Now().Add(time.Hour*24*30)),
	).SignedString(t.getSecret())
}

func (t UserToken) getSecret() []byte {
	return []byte(api.ENV.JWT_SECRET)
}

func (t UserToken) getClaims(exp time.Time) *Claims {
	return &Claims{
		User: &t,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
}

func ParseToken(tokenStr string) (*UserToken, error) {
	if tokenStr == "" {
		return nil, ErrMissingToken
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(api.ENV.JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims.User, nil
}

func NewUserToken(ID, email, name string) *UserToken {
	return &UserToken{
		ID:    ID,
		Email: email,
		Name:  name,
	}
}
