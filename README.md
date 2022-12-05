[![Build Status](https://github.com/agorman/repmon/workflows/repmon/badge.svg)](https://github.com/agorman/repmon/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/agorman/repmon)](https://goreportcard.com/report/github.com/agorman/repmon)
[![GoDoc](https://godoc.org/github.com/agorman/repmon?status.svg)](https://godoc.org/github.com/agorman/repmon)
[![codecov](https://codecov.io/gh/agorman/repmon/branch/main/graph/badge.svg)](https://codecov.io/gh/agorman/repmon)

# Repmon


Repmon is a simple database replication monitoring tool. Repmon will send notifications based on database replication failure. Repmon also optionally support HTTP healthchecks for database replication.


# Supported Databases


- MySQL 5.5 (possibly other MySQL versions but untested)


# Supported Notifications


- Email


# HTTP Healthcheck


Repmon can optionally listen for HTTP health probes at /healthcheck. It will return a 200 status code if replication is
running and a 503 otherwise.


# How does it work?


1. Download the [latest release](https://github.com/agorman/repmon/releases).
2. Create a YAML configuration file
3. Run it `repmon -conf repmon.yaml`


# Configuration file


The YAML file defines repmon's operation.


## Full config example

~~~
log_path: /var/log/repmon.log
log_level: error
frequency: 1h
http:
  addr: 0.0.0.0
  port: 4040
mysql:
  host: 127.0.0.1
  port: 3306
  user: user
  pass: pass
email:
  host: mail.me.com
  port: 587
  user: me
  pass: pass
  starttls: true
  ssl: false
  subject: Database Replication Failure
  from: me@me.com
  to:
    - you@me.com
~~~


## Global Options


**log_path** - File on disk where repmon logs will be stored. Defaults to /var/log/repmon.log.

**log_level** - Sets the log level. Valid levels are: panic, fatal, trace, debug, warn, info, and error. Defaults to error.

**frequency** - How often repmon will check the database to ensure replication is running. Defaults to 1h.


## HTTP


**addr** - The listening address for the HTTP server. Default to 127.0.0.1

**port** - The listening port for the HTTP server. Default to 4040


## MySQL


**host** - The hostname or IP of the MySQL server.

**port** - The port of the MySQL server.

**user** - The username used to authenticate.

**pass** - The password used to authenticate.

**socket_path** - Connect to the MySQL database through a socket file rather than a port.


## Email


**host** - The hostname or IP of the SMTP server.

**port** - The port of the SMTP server.

**user** - The username used to authenticate.

**pass** - The password used to authenticate.

**start_tls** - StartTLS enables TLS security. If both StartTLS and SSL are true then StartTLS will be used.

**insecure_skip_verify** - When using TLS skip verifying the server's certificate chain and host name.

**ssl** - SSL enables SSL security. If both StartTLS and SSL are true then StartTLS will be used.

**from** - The email address the email will be sent from.

**to** - An array of email addresses for which emails will be sent.


# Flags


**-conf** - Path to the repmon configuration file

**-debug** - Log to STDOUT


## Road Map


- Docker Image
- Systemd service file
- Create rpm
- Create deb
- Support for more databases
- Support for more notifiers