package repmon

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// MySQLReplicationChecker checks if replication is running on a given MySQL 5.5 database.
type MySQLReplicationChecker struct {
	db     *sql.DB
	config *Config
}

// MySQLReplicationChecker returns a new MySQL replication checker.
func NewMySQLReplicationChecker(db *sql.DB, config *Config) *MySQLReplicationChecker {
	return &MySQLReplicationChecker{
		db:     db,
		config: config,
	}
}

// Replicating checks if replication is running. It returns an error if replication
// isn't running or the replica is behind by more than the configured seconds threshold.
func (d *MySQLReplicationChecker) Replicating() error {
	rows, err := d.db.Query("SHOW SLAVE STATUS")
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	var slaveIORunning bool
	var slaveSQLRunning bool
	var secondsBehindMaster int

	values := make([]interface{}, len(columns))
	for rows.Next() {
		for i := 0; i < len(columns); i++ {
			values[i] = new(sql.RawBytes)
		}
		err := rows.Scan(values...)
		if err != nil {
			return err
		}

		// sql.RawBytes values must be checked before the next call to rows.Next
		for i, value := range values {
			rb, ok := value.(*sql.RawBytes)
			if !ok {
				continue
			}

			stringVal := string(*rb)

			log.Debug(columns[i], " => ", stringVal)

			switch columns[i] {
			case d.config.SlaveIORunningField:
				if strings.ToLower(stringVal) == "yes" {
					slaveIORunning = true
				}
			case d.config.SlaveSQLRunningField:
				if strings.ToLower(stringVal) == "yes" {
					slaveSQLRunning = true
				}
			case d.config.SecondsBehindMasterField:
				sbm, err := strconv.Atoi(stringVal)
				if err != nil {
					return fmt.Errorf("Unable to read Seconds_Behind_Master. Got %s", stringVal)
				}
				secondsBehindMaster = sbm
			}
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if !slaveIORunning {
		return fmt.Errorf("%s not equal to yes", d.config.SlaveIORunningField)
	}

	if !slaveSQLRunning {
		return fmt.Errorf("%s not equal to yes", d.config.SlaveSQLRunningField)
	}

	if secondsBehindMaster > d.config.SecondsBehindMasterThreshold {
		return fmt.Errorf("seconds behind master %d exceeds threshold %d", secondsBehindMaster, d.config.SecondsBehindMasterThreshold)
	}

	return nil
}
