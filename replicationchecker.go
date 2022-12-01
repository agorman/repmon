package repmon

// ReplicationChecker defines a replication check for a database
type ReplicationChecker interface {
	Replicating() error
}
