package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agorman/repmon"
	"github.com/etherlabsio/healthcheck/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	conf := flag.String("conf", "/etc/repmon.yaml", "Path to the repmon configuration file")
	debug := flag.Bool("debug", false, "Log to STDOUT")
	flag.Parse()

	config, err := repmon.OpenConfig(*conf)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", config.MySQL.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	replicationChecker := repmon.NewMySQLReplicationChecker(db, config)

	notifier := repmon.NewEmailNotifier(config)

	if !*debug {
		logfile := &lumberjack.Logger{
			Filename:   config.LogPath,
			MaxSize:    10,
			MaxBackups: 4,
		}
		log.SetOutput(logfile)
	}

	rm := repmon.New(config, replicationChecker, notifier)
	rm.Start()
	defer rm.Stop()

	errc := make(chan error, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	if config.HTTP != nil {
		http.Handle("/live", healthcheck.Handler(
			healthcheck.WithTimeout(5*time.Second),
			healthcheck.WithChecker(
				"live", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return nil
					},
				),
			),
		))

		http.Handle("/replicate", healthcheck.Handler(
			healthcheck.WithChecker(
				"replicate", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return replicationChecker.Replicating()
					},
				),
			),
		))

		go func() { errc <- http.ListenAndServe(fmt.Sprintf("%s:%d", config.HTTP.Addr, config.HTTP.Port), nil) }()
	}

	select {
	case s := <-sig:
		log.Warnf("Received signal %s, exiting", s)
	case e := <-errc:
		log.Errorf("Run error: %s", e)
	}
}
