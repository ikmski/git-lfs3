package adapter

import (
	"time"
)

const (
	contentMediaType = "application/vnd.git-lfs"
	metaMediaType    = "application/vnd.git-lfs+json"
)

// BatchRequest is ...
type BatchRequest struct {
	Operation string           `json:"operation"`
	Transfers []string         `json:"transfers,omitempty"`
	Ref       Ref              `json:"ref,omitempty"`
	Objects   []*ObjectRequest `json:"objects"`
}

// ObjectRequest is ...
type ObjectRequest struct {
	Oid      string `json:"oid"`
	Size     int64  `json:"size"`
	User     string `json:"user"`
	Password string `json:"password"`
	Repo     string `json:"repo"`
}

// Ref is ...
type Ref struct {
	Name string `json:"name"`
}

// BatchResponse is ...
type BatchResponse struct {
	Transfer string            `json:"transfer,omitempty"`
	Objects  []*ResponseObject `json:"objects"`
}

// ResponseObject is ...
type ResponseObject struct {
	Oid     string           `json:"oid"`
	Size    int64            `json:"size"`
	Actions map[string]*Link `json:"actions"`
	Error   *ObjectError     `json:"error,omitempty"`
}

type LockRequest struct {
	Path string `json:"path"`
}

type LockResponse struct {
	Lock    *Lock  `json:"lock"`
	Message string `json:"message,omitempty"`
}

type UnlockRequest struct {
	Force bool `json:"force"`
}

type UnlockResponse struct {
	Lock    *Lock  `json:"lock"`
	Message string `json:"message,omitempty"`
}

type LockListResponse struct {
	Locks      []Lock `json:"locks"`
	NextCursor string `json:"next_cursor,omitempty"`
	Message    string `json:"message,omitempty"`
}

type LockVerifyRequest struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type LockVerifyResponse struct {
	Ours       []Lock `json:"ours"`
	Theirs     []Lock `json:"theirs"`
	NextCursor string `json:"next_cursor,omitempty"`
	Message    string `json:"message,omitempty"`
}

type Lock struct {
	ID       string    `json:"id"`
	Path     string    `json:"path"`
	Owner    User      `json:"owner"`
	LockedAt time.Time `json:"locked_at"`
}

type User struct {
	Name string `json:"name"`
}

func newResponseObject() *ResponseObject {
	r := new(ResponseObject)
	r.Actions = make(map[string]*Link)
	return r
}

// Link is ...
type Link struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}

// ObjectError is ...
type ObjectError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
