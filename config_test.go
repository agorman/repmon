package repmon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config, err := OpenConfig("./testdata/repmon.yaml")
	assert.Nil(t, err)

	assert.Equal(t, config.LogPath, "/var/log/remon_log.log")
	assert.Equal(t, config.LogLevel, "debug")
	assert.Equal(t, config.Frequency, "10s")

	assert.NotNil(t, config.HTTP)
	assert.Equal(t, config.HTTP.Addr, "127.0.0.1")
	assert.Equal(t, config.HTTP.Port, 1234)

	assert.NotNil(t, config.MySQL)
	assert.Equal(t, config.MySQL.Host, "127.0.0.1")
	assert.Equal(t, config.MySQL.Port, 3308)
	assert.Equal(t, config.MySQL.User, "user")
	assert.Equal(t, config.MySQL.Pass, "pass")

	assert.NotNil(t, config.Email)
	assert.Equal(t, config.Email.Host, "smtp.me.com")
	assert.Equal(t, config.Email.Port, 25)
	assert.Equal(t, config.Email.User, "me")
	assert.Equal(t, config.Email.Pass, "pass")
	assert.Equal(t, config.Email.StartTLS, true)
	assert.Equal(t, config.Email.InsecureSkipVerify, false)
	assert.Equal(t, config.Email.SSL, false)
	assert.Equal(t, config.Email.From, "me@me.com")
	assert.Contains(t, config.Email.To, "they@me.com")
	assert.Contains(t, config.Email.To, "them@me.com")
	assert.Equal(t, config.Email.Subject, "REP FAILED!")
}
