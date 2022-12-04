package repmon

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	// LogPath is the directory on disk where repmon logs will be stored. Defaults to /var/log/repmon.
	LogPath string `yaml:"log_path"`

	// LogLevel sets the level of logging. Valid levels are: panic, fatal, trace, debug, warn, info, and error. Defaults to error
	LogLevel string `yaml:"log_level"`

	// Frequency describes how often the replication status will be checked
	Frequency string `yaml:"frequency"`
	frequency time.Duration

	// SlaveSQLRunningField is the field checked in the database determine if slave sql is running. Defaults to Slave_SQL_Running and likely never
	// needs to be changed.
	SlaveSQLRunningField string `yaml:"slave_sql_running_field"`

	// SlaveIORunningField is the field checked in the database to determine if slave io is running. Defaults to Slave_IO_Running and likely never
	// needs to be changed.
	SlaveIORunningField string `yaml:"slave_io_running_field"`

	// SecondsBehindMasterField is the field checked in the database to determine if the seconds replica is behind master.
	// Defaults to Seconds_Behind_Master and likely never needs to be changed.
	SecondsBehindMasterField string `yaml:"seconds_behind_master_field"`

	// SecondsBehindMasterThreshold is the seconds behind master where we will send a notification email. It must be greather
	// than 0 and defaults to 36000 (10 hours).
	SecondsBehindMasterThreshold int `yaml:"seconds_behind_master_threshold"`

	HTTP  *HTTP  `yaml:"http"`
	MySQL *MySQL `yaml:"mysql"`
	Email *Email `yaml:"email"`
}

func (c *Config) validate() error {
	if c.LogPath == "" {
		c.LogPath = "/var/log/repmon.log"
	}

	if c.LogLevel == "" {
		c.LogLevel = "error"
		log.SetLevel(log.ErrorLevel)
	} else {
		switch c.LogLevel {
		case "panic":
			log.SetLevel(log.PanicLevel)
		case "fatal":
			log.SetLevel(log.FatalLevel)
		case "trace":
			log.SetLevel(log.TraceLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		default:
			return fmt.Errorf("Invalid log_level: %s", c.LogLevel)
		}
	}

	if c.Frequency == "" {
		c.Frequency = "1h"
	}
	var err error
	c.frequency, err = time.ParseDuration(c.Frequency)
	if err != nil {
		return fmt.Errorf("Config: invalid frequency: %w", err)
	}

	if c.SlaveIORunningField == "" {
		c.SlaveIORunningField = "Slave_IO_Running"
	}

	if c.SlaveSQLRunningField == "" {
		c.SlaveSQLRunningField = "Slave_SQL_Running"
	}

	if c.SecondsBehindMasterField == "" {
		c.SecondsBehindMasterField = "Seconds_Behind_Master"
	}

	if c.SecondsBehindMasterThreshold == 0 {
		c.SecondsBehindMasterThreshold = 36000
	}

	if c.HTTP != nil {
		if c.HTTP.Addr == "" {
			c.HTTP.Addr = "127.0.0.1"
		}

		if c.HTTP.Port == 0 {
			c.HTTP.Port = 4040
		}
	}

	if c.Email == nil {
		return errors.New("Missing required email configuration")
	}

	if c.Email.Host == "" {
		return errors.New("Missing required host entry for email")
	}

	if c.Email.Port == 0 {
		c.Email.Port = 25
	}

	// StartTLS takes presidence over SSL
	if c.Email.StartTLS {
		c.Email.SSL = false
	}

	if c.Email.Subject == "" {
		c.Email.Subject = "Database Replication Failure"
	}

	if c.Email.From == "" {
		return errors.New("Missing required from entry for email")
	}

	if len(c.Email.To) == 0 {
		return errors.New("Missing required to entry for email")
	}

	if c.MySQL == nil {
		return errors.New("Missing required mysql configuration")
	}

	if c.MySQL.Host == "" {
		return errors.New("Missing required host entry for mysql")
	}

	if c.MySQL.Port == 0 {
		c.MySQL.Port = 3306
	}

	if c.MySQL.User == "" {
		return errors.New("Missing required user entry for mysql")
	}

	if c.MySQL.Pass == "" {
		return errors.New("Missing required pass entry for mysql")
	}

	return nil

}

type MySQL struct {
	// Host is the hostname or IP of the MySQL server.
	Host string `yaml:"host"`

	// Port is the port of the MySQL server. Ignored if SocketPath is not the empty string
	Port int `yaml:"port"`

	// User is the username used to authenticate.
	User string `yaml:"user"`

	// Pass is the password used to authenticate.
	Pass string `yaml:"pass"`

	// SocketPath is the path to the unix socket. If set to the empty string then a TCP connection is used instead.
	SocketPath string `yaml:"socket_path"`
}

// DSN returns the go database dsn.
func (m *MySQL) DSN() string {
	if m.SocketPath != "" {
		return fmt.Sprintf("%s:%s@unix(%s)/", m.User, m.Pass, m.SocketPath)
	} else {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/", m.User, m.Pass, m.Host, m.Port)
	}
}

// HTTP defines the configuration for http health checks.
type HTTP struct {
	// The address the http server will listen on.
	Addr string `yaml:"addr"`

	// The port the http server will listen on.
	Port int `yaml:"port"`
}

type Email struct {
	// Host is the hostname or IP of the SMTP server.
	Host string `yaml:"host"`

	// Port is the port of the SMTP server.
	Port int `yaml:"port"`

	// User is the username used to authenticate.
	User string `yaml:"user"`

	// Pass is the password used to authenticate.
	Pass string `yaml:"pass"`

	// StartTLS enables TLS security. If both StartTLS and SSL are true then StartTLS will be used.
	StartTLS bool `yaml:"starttls"`

	// Skip verifying the server's certificate chain and host name.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`

	// SSL enables SSL security. If both StartTLS and SSL are true then StartTLS will be used.
	SSL bool `yaml:"ssl"`

	// Optional subject field for notification emails
	Subject string `yaml:"subject"`

	// From is the email address the email will be sent from.
	From string `yaml:"from"`

	// To is an array of email addresses for which emails will be sent.
	To []string `yaml:"to"`
}

// OpenConfig returns a new Config option by reading the YAML file at path. If the file
// doesn't exist, can't be read, is invalid YAML, or doesn't match the repmon spec then
// an error is returned.
func OpenConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	return config, config.validate()
}
