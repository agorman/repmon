package repmon

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestRepMon(t *testing.T) {
	config, err := OpenConfig("./testdata/repmon.yaml")
	assert.Nil(t, err)

	db, err := sql.Open("mysql", config.MySQL.DSN())
	assert.Nil(t, err)
	defer db.Close()

	checker := NewMySQLReplicationChecker(db, config)

	notifier := NewEmailNotifier(config)

	rm := New(config, checker, notifier)
	rm.Start()
	rm.Start()
	rm.Stop()
	rm.Stop()
}
