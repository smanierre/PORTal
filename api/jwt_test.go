package api

import (
	"PORTal/testutils"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"testing"
	"time"
)

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte("test"), nil
}

func TestCreateAndValidateUserToken(t *testing.T) {
	k, _ := keyFunc(nil)
	key, _ := k.([]byte)
	normalMember := testutils.RandomMember(false)
	normalMember.ID = uuid.NewString()
	token, err := createToken(normalMember, time.Hour, key)
	if err != nil {
		t.Fatalf("Error when creating token: %s", err.Error())
	}
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, keyFunc)
	if err != nil {
		t.Fatalf("Error when validating token: %s", err.Error())
	}
	cc, ok := parsedToken.Claims.(*CustomClaims)
	if !ok {
		t.Fatal("Error when casting claims to CustomClaims")
	}
	if cc.Subject != normalMember.ID {
		t.Errorf("Expected subject to be: %s\nGot: %s", normalMember.ID, cc.Subject)
	}
	if cc.Admin != false {
		t.Error("Expected admin to be false, but got true")
	}
	if cc.ExpiresAt.Time.Before(time.Now().Add(59*time.Minute)) && cc.ExpiresAt.Time.After(time.Now().Add(61*time.Minute)) {
		t.Error("Expected JWT to last an hour, but it didn't")
	}
}

func TestCreateAndValidateAdminToken(t *testing.T) {
	k, _ := keyFunc(nil)
	key, _ := k.([]byte)
	adminMember := testutils.RandomMember(true)
	adminMember.ID = uuid.NewString()
	token, err := createToken(adminMember, time.Hour, key)
	if err != nil {
		t.Fatalf("Error when creating token: %s", err.Error())
	}
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, keyFunc)
	if err != nil {
		t.Fatalf("Error when validating token: %s", err.Error())
	}
	cc, ok := parsedToken.Claims.(*CustomClaims)
	if !ok {
		t.Fatal("Error when casting claims to CustomClaims")
	}
	if cc.Subject != adminMember.ID {
		t.Errorf("Expected subject to be: %s\nGot: %s", adminMember.ID, cc.Subject)
	}
	if cc.Admin != true {
		t.Error("Expected admin to be true, but got false")
	}
	if cc.ExpiresAt.Time.Before(time.Now().Add(59*time.Minute)) && cc.ExpiresAt.Time.After(time.Now().Add(61*time.Minute)) {
		t.Error("Expected JWT to last an hour, but it didn't")
	}
}

func TestExpiredJWT(t *testing.T) {
	k, _ := keyFunc(nil)
	key, _ := k.([]byte)
	adminMember := testutils.RandomMember(true)
	adminMember.ID = uuid.NewString()
	token, err := createToken(adminMember, time.Millisecond, key)
	if err != nil {
		t.Fatalf("Error when creating token: %s", err.Error())
	}
	_, err = jwt.ParseWithClaims(token, &CustomClaims{}, keyFunc)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("Expected token to be expired but got error: %s", err.Error())
	}
}
