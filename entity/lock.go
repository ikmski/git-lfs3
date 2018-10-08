package entity

import "time"

// Lock is ...
type Lock struct {
	ID       string
	Path     string
	Owner    User
	LockedAt time.Time
}
