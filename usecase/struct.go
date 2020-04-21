package usecase

// BatchRequest is ...
type BatchRequest struct {
	Objects []*ObjectRequest
}

// ObjectRequest is ...
type ObjectRequest struct {
	Oid  string
	Size int64
	From int64
	To   int64
}

// BatchResult is ...
type BatchResult struct {
	Objects []*ObjectResult
}

// ObjectResult is ...
type ObjectResult struct {
	Oid          string
	Size         int64
	MetaExists   bool
	ObjectExists bool
}

type LockRequest struct {
	Repo    string
	User    string
	Path    string
	Refspec string
}

type UnlockRequest struct {
	Repo    string
	User    string
	ID      string
	Force   bool
	Refspec string
}

type LockListRequest struct {
	Repo    string
	Path    string
	Cursor  string
	Limit   int
	Refspec string
}

type LockVerifyRequest struct {
	Repo    string
	User    string
	Path    string
	Cursor  string
	Limit   int
	Refspec string
}

type LockResult struct {
	ID           string
	Path         string
	Owner        string
	LockedAt     int64 // UnixTime
	AlreadyExist bool
}

type LockListResult struct {
	Locks      []*LockResult
	NextCursor string
}

type LockVerifyResult struct {
	Ours       []*LockResult
	Theirs     []*LockResult
	NextCursor string
}
