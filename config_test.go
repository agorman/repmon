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

func TestConfigDefaults(t *testing.T) {
	config, err := OpenConfig("./testdata/defaults.yaml")
	assert.Nil(t, err)

	config.HTTP = &HTTP{}

	err = config.validate()
	assert.Nil(t, err)
	assert.Equal(t, config.LogPath, "/var/log/repmon.log")
	assert.Equal(t, config.LogLevel, "error")
	assert.Equal(t, config.Frequency, "1h")
	assert.Equal(t, config.SlaveIORunningField, "Slave_IO_Running")
	assert.Equal(t, config.SlaveSQLRunningField, "Slave_SQL_Running")
	assert.Equal(t, config.SecondsBehindMasterField, "Seconds_Behind_Master")
	assert.Equal(t, config.SecondsBehindMasterThreshold, 36000)

	assert.Equal(t, config.Email.Port, 25)
	assert.Equal(t, config.Email.StartTLS, false)
	assert.Equal(t, config.Email.SSL, false)
	assert.Equal(t, config.Email.Subject, "Database Replication Failure")

	assert.Equal(t, config.MySQL.Port, 3306)
	assert.Equal(t, config.MySQL.SocketPath, "")

	assert.Equal(t, config.HTTP.Addr, "127.0.0.1")
	assert.Equal(t, config.HTTP.Port, 4040)
}

func TestConfigInvalidFrequency(t *testing.T) {
	config, err := OpenConfig("./testdata/defaults.yaml")
	assert.Nil(t, err)

	config.Frequency = "Invalid"

	err = config.validate()
	assert.Error(t, err)
}

func TestConfigRequired(t *testing.T) {
	config := &Config{}

	err := config.validate()
	assert.Error(t, err)

	config.Email = &Email{}
	err = config.validate()
	assert.Error(t, err)

	config.Email.Host = "smtp.me.com"
	err = config.validate()
	assert.Error(t, err)

	config.Email.From = "me@me.com"
	err = config.validate()
	assert.Error(t, err)

	config.Email.To = []string{"you@me.com"}
	err = config.validate()
	assert.Error(t, err)

	config.MySQL = &MySQL{}
	err = config.validate()
	assert.Error(t, err)

	config.MySQL.Host = "127.0.0.1"
	err = config.validate()
	assert.Error(t, err)

	config.MySQL.User = "user"
	err = config.validate()
	assert.Error(t, err)

	config.MySQL.Pass = "pass"
	err = config.validate()
	assert.Nil(t, err)
}

func TestConfigBadPath(t *testing.T) {
	_, err := OpenConfig("./testdata/notexist.yaml")
	assert.Error(t, err)
}
