package repmon

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// RepMon monitors a database to ensure that replication is running.
type RepMon struct {
	config             *Config
	replicationChecker ReplicationChecker
	notifier           Notifier
	running            bool
	stopc              chan struct{}
	donec              chan struct{}
}

// New returns a new RepMon instance.
func New(config *Config, replicationChecker ReplicationChecker, notifier Notifier) *RepMon {
	return &RepMon{
		config:             config,
		replicationChecker: replicationChecker,
		notifier:           notifier,
		stopc:              make(chan struct{}),
		donec:              make(chan struct{}),
	}
}

// Start starts monitoring the configured database.
func (r *RepMon) Start() {
	if r.running {
		return
	}

	r.running = true

	go r.loop()
}

// Stop stops monitoring the configured database.
func (r *RepMon) Stop() {
	if !r.running {
		return
	}

	r.stopc <- struct{}{}
	<-r.donec

	r.running = false
}

func (r *RepMon) loop() {
	ticker := time.NewTicker(r.config.frequency)

	log.Infof("RepMon started: checking slave replication every %s", r.config.Frequency)

	for {
		select {
		case <-ticker.C:
			if err := r.replicationChecker.Replicating(); err != nil {
				r.notifier.Notify(err)
			}
		case <-r.stopc:
			ticker.Stop()
			r.donec <- struct{}{}
			log.Info("RepMon shutdown")
			return
		}
	}
}
