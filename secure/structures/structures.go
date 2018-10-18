package structures

import "github.com/dgrijalva/jwt-go"

type TokenManagerI interface {
	CreateToken(tk *TokenClaims, cipher string) (string, error)
	ValidateToken(token string, cipher string) (*TokenClaims, error)
}

type TokenClaims struct {
	Email string `json:"email"`
	ID    int64  `json:"id"`
	jwt.StandardClaims
}
