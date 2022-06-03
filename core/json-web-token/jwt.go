package json_web_token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mecamon/shoppingify-server/models"
	"strings"
	"time"
)

var hmacSampleSecret = make([]byte, 32)

func Generate(ID int64, email string) (string, error) {
	issuedAt := &jwt.NumericDate{Time: time.Now()}
	expiresAt := &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.CustomClaims{
		TokenType: "level1",
		ID:        ID,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    email,
			ExpiresAt: expiresAt,
			NotBefore: issuedAt,
			IssuedAt:  issuedAt,
		},
	})

	signedToken, err := token.SignedString(hmacSampleSecret)
	return signedToken, err
}

func Validate(tokenStr string) (*models.CustomClaims, error) {
	cleanTokenString := strings.Replace(tokenStr, "Bearer ", "", -1)
	token, err := jwt.ParseWithClaims(cleanTokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(*models.CustomClaims)
	return claims, nil
}
