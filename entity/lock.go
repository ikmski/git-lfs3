package entity

// Lock is ...
type Lock struct {
	ID       string
	Path     string
	Owner    User
	LockedAt int64 // UnixTime
}
