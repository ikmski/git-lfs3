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
	Path    string
	Refspec string
}

type LockListRequest struct {
	Path    string
	ID      string
	Cursor  string
	Limit   int32
	Refspec string
}

type LockVerifyRequest struct {
	Cursor  string
	Limit   int32
	Refspec string
}

type UnlockRequest struct {
	Force   bool
	Refspec string
}

type LockResult struct {
	ID       string
	Path     string
	Owner    string
	LockedAt int64 // UnixTime
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
