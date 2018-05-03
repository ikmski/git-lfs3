package main

import (
	"fmt"
	"time"
)

type RequestVars struct {
	Oid      string
	Size     int64
	User     string
	Password string
	Repo     string
}

type BatchVars struct {
	Transfers []string       `json:"transfers,omitempty"`
	Operation string         `json:"operation"`
	Objects   []*RequestVars `json:"objects"`
}

type MetaObject struct {
	Oid      string `json:"oid"`
	Size     int64  `json:"size"`
	Existing bool
}

type BatchResponse struct {
	Transfer string            `json:"transfer,omitempty"`
	Objects  []*Representation `json:"objects"`
}

type Representation struct {
	Oid     string           `json:"oid"`
	Size    int64            `json:"size"`
	Actions map[string]*link `json:"actions"`
	Error   *ObjectError     `json:"error,omitempty"`
}

type ObjectError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type User struct {
	Name string `json:"name"`
}

type Lock struct {
	Id       string    `json:"id"`
	Path     string    `json:"path"`
	Owner    User      `json:"owner"`
	LockedAt time.Time `json:"locked_at"`
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

type LockList struct {
	Locks      []Lock `json:"locks"`
	NextCursor string `json:"next_cursor,omitempty"`
	Message    string `json:"message,omitempty"`
}

type VerifiableLockRequest struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type VerifiableLockList struct {
	Ours       []Lock `json:"ours"`
	Theirs     []Lock `json:"theirs"`
	NextCursor string `json:"next_cursor,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (v *RequestVars) DownloadLink() string {

	return v.internalLink("objects")
}

func (v *RequestVars) UploadLink(useTus bool) string {

	return v.internalLink("objects")
}

func (v *RequestVars) internalLink(subpath string) string {

	path := ""

	if len(v.User) > 0 {
		path += fmt.Sprintf("/%s", v.User)
	}

	if len(v.Repo) > 0 {
		path += fmt.Sprintf("/%s", v.Repo)
	}

	path += fmt.Sprintf("/%s/%s", subpath, v.Oid)

	if config.Server.Tls {
		return fmt.Sprintf("https://%s%s", config.Server.Host, path)
	}

	return fmt.Sprintf("http://%s%s", config.Server.Host, path)
}

func (v *RequestVars) VerifyLink() string {

	path := fmt.Sprintf("/verify/%s", v.Oid)

	if config.Server.Tls {
		return fmt.Sprintf("https://%s%s", config.Server.Host, path)
	}

	return fmt.Sprintf("http://%s%s", config.Server.Host, path)
}

type link struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}
