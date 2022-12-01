package repmon

// Notifier defines a notification method.
type Notifier interface {
	Notify(error) error
}
