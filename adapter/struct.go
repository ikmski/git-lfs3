package adapter

import (
	"time"

	"github.com/ikmski/git-lfs3/usecase"
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
	User     string
	Password string
	Repo     string
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

func convertBatchRequest(req *BatchRequest) *usecase.BatchRequest {

	var objs []*usecase.ObjectRequest

	for _, o := range req.Objects {

		item := &usecase.ObjectRequest{
			Oid:      o.Oid,
			Size:     o.Size,
			User:     o.User,
			Password: o.Password,
			Repo:     o.Repo,
		}

		objs = append(objs, item)
	}

	br := &usecase.BatchRequest{
		Operation: req.Operation,
		Transfers: req.Transfers,
		Ref:       req.Ref.Name,
		Objects:   objs,
	}

	return br
}

func convertBatchResponse(result *usecase.BatchResult) *BatchResponse {

	var objs []*ResponseObject

	res := &BatchResponse{
		Transfer: "basic",
		Objects:  objs,
	}

	for _, batchObj := range result.Objects {

		header := make(map[string]string)
		header["Accept"] = contentMediaType

		obj := &ResponseObject{}

		if batchObj.MetaExists {

		}

		if batchObj.ObjectExists {
			obj.Actions["download"] = &Link{
				Href:   "https://hoge",
				Header: header,
			}
		} else {
			obj.Actions["upload"] = &Link{
				Href:   "https://hoge",
				Header: header,
			}
		}

		objs = append(objs, obj)
	}

	return res
}

func convertObjectRequest(req *ObjectRequest) *usecase.ObjectRequest {

	or := &usecase.ObjectRequest{
		Oid:      req.Oid,
		Size:     req.Size,
		User:     req.User,
		Password: req.Password,
		Repo:     req.Repo,
	}

	return or
}
