package repmon

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMySQLReplicationChecker(t *testing.T) {
	config, err := OpenConfig("./testdata/repmon.yaml")
	assert.Nil(t, err)

	db, err := sql.Open("mysql", config.MySQL.DSN())
	assert.Nil(t, err)
	defer db.Close()

	checker := NewMySQLReplicationChecker(db, config)
	err = checker.Replicating()
	assert.Error(t, err)
}
