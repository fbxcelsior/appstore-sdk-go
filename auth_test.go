package appstore

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Auth_AuthToken_IsNotExpiredSuccess(t *testing.T) {
	expired := time.Now().Unix() + 1000
	token := AuthToken{ExpiresAt: expired}
	assert.True(t, token.IsNotExpired())
}

func Test_Auth_AuthToken_IsNotExpiredFail(t *testing.T) {
	expired := time.Now().Unix() - 1000
	token := AuthToken{ExpiresAt: expired}
	assert.False(t, token.IsNotExpired())
}

func Test_Auth_AuthToken_IsValidSuccess(t *testing.T) {
	expired := time.Now().Unix() + 1000
	token := AuthToken{ExpiresAt: expired, Token: "qwerty"}
	assert.True(t, token.IsValid())
}

func Test_Auth_AuthToken_IsValidFailExpired(t *testing.T) {
	expired := time.Now().Unix() - 1000
	token := AuthToken{ExpiresAt: expired, Token: "qwerty"}
	assert.False(t, token.IsValid())
}

func Test_Auth_AuthToken_IsValidFailEmpty(t *testing.T) {
	expired := time.Now().Unix() + 1000
	token := AuthToken{ExpiresAt: expired}
	assert.False(t, token.IsValid())
}

func Test_Auth_NewTokenBuilder(t *testing.T) {
	auth := NewTokenBuilder(buildStubConfig())
	assert.NotEmpty(t, auth)
	assert.NotEmpty(t, auth.cfg)
}

func Test_Auth_TokenBuilder_BuildPayload(t *testing.T) {
	cfg := buildStubConfig()
	auth := NewTokenBuilder(cfg)
	payload := auth.BuildPayload()
	assert.NotEmpty(t, payload)
	assert.NotEmpty(t, payload.ExpiresAt)
	assert.Equal(t, cfg.Token.Audience, payload.Audience)
	assert.Equal(t, cfg.IssuerId, payload.Issuer)
}

func Test_Auth_TokenBuilder_BuildJWTToken(t *testing.T) {
	cfg := buildStubConfig()
	auth := NewTokenBuilder(cfg)
	payload := auth.BuildPayload()
	jwtToken := auth.BuildJWTToken(payload)
	assert.NotEmpty(t, jwtToken)
	assert.Equal(t, jwtToken.Claims, payload)
	assert.Equal(t, jwtToken.Method, jwt.SigningMethodES256)
	assert.Equal(t, jwtToken.Header["typ"], cfg.Token.Type)
	assert.Equal(t, jwtToken.Header["alg"], cfg.Token.Algo)
	assert.Equal(t, jwtToken.Header["kid"], cfg.KeyId)
}

func Test_Auth_TokenBuilder_BuildAuthTokenSuccess(t *testing.T) {
	cfg := buildStubConfig()
	auth := NewTokenBuilder(cfg)
	token, _ := auth.BuildAuthToken()
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, token.Token)
	assert.NotEmpty(t, token.ExpiresAt)
	assert.True(t, token.IsValid())
	assert.True(t, token.IsNotExpired())
}

func Test_Auth_TokenBuilder_BuildAuthTokenError(t *testing.T) {
	cfg := buildStubConfig()
	cfg.PrivateKey = "stubs/auth/keys/fail.p8"
	auth := NewTokenBuilder(cfg)
	_, err := auth.BuildAuthToken()
	assert.Error(t, err)
	assert.Equal(t, "PrivateKey.DecodePem: AuthKey must be a valid .p8 PEM file", err.Error())
}
