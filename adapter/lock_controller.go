package adapter

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ikmski/git-lfs3/usecase"
)

// LockController is ...
type LockController interface {
	Lock(ctx Context)
	Unlock(ctx Context)
	List(ctx Context)
	Verify(ctx Context)
}

type lockController struct {
	LockService usecase.LockService
}

func NewLockController(s usecase.LockService) LockController {
	return &lockController{
		LockService: s,
	}
}

func (c *lockController) Lock(ctx Context) {

}

func (c *lockController) Unlock(ctx Context) {

}

func (c *lockController) List(ctx Context) {

	req, err := parseListRequest(ctx)
	if err != nil {

	}

	result, err := c.LockService.List(req)
	if err != nil {

	}

	res := convertListResponse(result)
	json, err := json.Marshal(res)
	if err != nil {

	}

	ctx.SetHeader("Content-Type", metaMediaType)
	ctx.SetStatus(200)
	ctx.GetResponseWriter().Write(json)
}

func (c *lockController) Verify(ctx Context) {

}

func parseListRequest(ctx Context) (*usecase.LockListRequest, error) {

	repo := ctx.GetParam("repo")
	path := ctx.GetParam("path")
	cursor := ctx.GetParam("cursor")
	limitValue := ctx.GetParam("limit")
	limit := 0
	if limitValue != "" {
		limit, _ = strconv.Atoi(limitValue)
	}

	req := &usecase.LockListRequest{
		Repo:   repo,
		Path:   path,
		Cursor: cursor,
		Limit:  limit,
	}

	return req, nil
}

func convertListResponse(result *usecase.LockListResult) *LockListResponse {

	var locks []Lock

	for _, l := range result.Locks {

		lock := Lock{
			ID:       l.ID,
			Path:     l.Path,
			Owner:    User{Name: l.Owner},
			LockedAt: time.Unix(l.LockedAt, 0),
		}

		locks = append(locks, lock)
	}

	res := &LockListResponse{
		Locks:      locks,
		NextCursor: result.NextCursor,
	}

	return res
}
