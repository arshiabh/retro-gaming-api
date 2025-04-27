package auth

import "github.com/golang-jwt/jwt/v5"

//define interface to use in main app for dip
type Authenticator interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
