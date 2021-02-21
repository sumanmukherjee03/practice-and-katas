package access_token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenConstants(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(expirationTime, 24, "expirationTime should be 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	assert := assert.New(t)
	at := GetNewAccessToken(1)
	assert.False(at.IsExpired(), "should not have expired")
	assert.NotEmpty(at.AccessToken, "should be empty")
	assert.NotZero(at.UserId, "UserId should be zero")
	assert.NotZero(at.ClientId, "ClientId should be zero")
}

func TestAccessTokenIsExpired(t *testing.T) {
	assert := assert.New(t)
	at0 := AccessToken{}
	assert.True(at0.IsExpired(), "access token without expiry should have expired")
	at0.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(at0.IsExpired(), "access token expiring 3 hours from now should not have expired")
}
