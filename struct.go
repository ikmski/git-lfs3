package main

import (
	"fmt"
	"time"
)

type ObjectRequest struct {
	Oid      string `json:"oid"`
	Size     int64  `json:"size"`
	User     string
	Password string
	Repo     string
}

type BatchRequest struct {
	Operation string           `json:"operation"`
	Transfers []string         `json:"transfers,omitempty"`
	Ref       Ref              `json:"ref,omitempty"`
	Objects   []*ObjectRequest `json:"objects"`
}

type ObjectMetaData struct {
	Oid      string `json:"oid"`
	Size     int64  `json:"size"`
	Existing bool
}

type BatchResponse struct {
	Transfer string            `json:"transfer,omitempty"`
	Objects  []*ResponseObject `json:"objects"`
}

type ResponseObject struct {
	Oid     string           `json:"oid"`
	Size    int64            `json:"size"`
	Actions map[string]*Link `json:"actions"`
	Error   *ObjectError     `json:"error,omitempty"`
}

type Ref struct {
	Name string `json:"name"`
}

type Link struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
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

func (o *ObjectRequest) DownloadLink() string {

	return o.internalLink("objects")
}

func (o *ObjectRequest) UploadLink() string {

	return o.internalLink("objects")
}

func (o *ObjectRequest) internalLink(subpath string) string {

	path := ""

	if len(o.User) > 0 {
		path += fmt.Sprintf("/%s", o.User)
	}

	if len(o.Repo) > 0 {
		path += fmt.Sprintf("/%s", o.Repo)
	}

	path += fmt.Sprintf("/%s/%s", subpath, o.Oid)

	if config.Server.Tls {
		return fmt.Sprintf("https://%s%s", config.Server.Host, path)
	}

	return fmt.Sprintf("http://%s%s", config.Server.Host, path)
}

func (o *ObjectRequest) VerifyLink() string {

	path := fmt.Sprintf("/verify/%s", o.Oid)

	if config.Server.Tls {
		return fmt.Sprintf("https://%s%s", config.Server.Host, path)
	}

	return fmt.Sprintf("http://%s%s", config.Server.Host, path)
}
