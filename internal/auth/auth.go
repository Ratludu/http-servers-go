package auth

import (
	"errors"
	"time"

	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPw), nil
}

func CheckPasswordHash(password, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err

}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))

	return ss, err

}

func ValidateJWT(tokenstring, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return parsedId, nil

}

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token := authHeader[7:]
		return token, nil
	}

	return "", errors.New("Token not foud")
}
