package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPasswords(t *testing.T) {

	cases := []struct {
		name     string
		password string
	}{
		{
			name:     "General",
			password: "Password123",
		},
		{
			name:     "Complex",
			password: "Py$YE.T3ptt.`s8%",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := HashPassword(tc.password)
			if err != nil {
				t.Errorf("Error: %v", err)
				return
			}

			err = CheckPasswordHash(tc.password, got)
			if err != nil {
				t.Errorf("Error: %v", err)
				return
			}
		})
	}

}

func TestMakeJWTRegular(t *testing.T) {

	newUserId := uuid.New()
	tokenSecret := "mySecret"
	expiry, err := time.ParseDuration("1m")
	if err != nil {
		t.Fatal(err)
		return
	}

	ss, err := MakeJWT(newUserId, tokenSecret, expiry)
	if err != nil {
		t.Fatal(err)
		return
	}

	parsedId, err := ValidateJWT(ss, tokenSecret)
	if err != nil {
		t.Fatal(err)
		return
	}

	if newUserId != parsedId {
		t.Errorf("Expeced %s got %s", newUserId, parsedId)
	}

}
func TestMakeJWTWrongToken(t *testing.T) {

	newUserId := uuid.New()
	tokenSecret := "mySecret"
	expiry, err := time.ParseDuration("1m")
	if err != nil {
		t.Fatal(err)
		return
	}

	ss, err := MakeJWT(newUserId, tokenSecret, expiry)
	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = ValidateJWT(ss, "wrong token")
	if err == nil {
		t.Fatal(err)
		return
	}

}
func TestMakeJWTExpired(t *testing.T) {

	newUserId := uuid.New()
	tokenSecret := "mySecret"
	expiry, err := time.ParseDuration("1ms")
	if err != nil {
		t.Fatal(err)
		return
	}

	ss, err := MakeJWT(newUserId, tokenSecret, expiry)
	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = ValidateJWT(ss, tokenSecret)
	if err == nil {
		t.Fatal(err)
		return
	}

}
