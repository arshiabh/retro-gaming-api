package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// this struct implement authenticator that define auth so we can use it in app
type JWTAuth struct {
	secret string
}

func NewAuthentication(secret string) *JWTAuth {
	return &JWTAuth{
		secret: secret,
	}
}

func (a *JWTAuth) GenerateToken(userID int64) (string, error) {
	claim := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuth) ValidateToken(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(a.secret), nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["sub"])
	} else {
		fmt.Println(err)
	}
	return jwtToken, err
}
