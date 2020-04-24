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

	req, err := parseLockRequest(ctx)
	if err != nil {

	}

	result, err := c.LockService.Lock(req)
	if err != nil {

	}

	res := convertLockResponse(result)
	json, err := json.Marshal(res)
	if err != nil {

	}

	ctx.SetHeader("Content-Type", metaMediaType)
	ctx.SetStatus(200)
	ctx.GetResponseWriter().Write(json)
}

func (c *lockController) Unlock(ctx Context) {

	req, err := parseUnlockRequest(ctx)
	if err != nil {

	}

	result, err := c.LockService.Unlock(req)
	if err != nil {

	}

	res := convertUnlockResponce(result)
	json, err := json.Marshal(res)
	if err != nil {

	}

	ctx.SetHeader("Content-Type", metaMediaType)
	ctx.SetStatus(200)
	ctx.GetResponseWriter().Write(json)
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

	req, err := parseVerifyRequest(ctx)
	if err != nil {

	}

	result, err := c.LockService.Verify(req)
	if err != nil {

	}

	res := convertVerifyResponse(result)
	json, err := json.Marshal(res)
	if err != nil {

	}

	ctx.SetHeader("Content-Type", metaMediaType)
	ctx.SetStatus(200)
	ctx.GetResponseWriter().Write(json)
}

func parseLockRequest(ctx Context) (*usecase.LockRequest, error) {

	data, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}

	var req usecase.LockRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	req.Repo = ctx.GetParam("repo")
	req.User = ctx.GetParam("user")

	return &req, nil
}

func convertLockResponse(result *usecase.LockResult) *LockResponse {

	lock := &Lock{
		ID:       result.ID,
		Path:     result.Path,
		Owner:    User{Name: result.Owner},
		LockedAt: time.Unix(result.LockedAt, 0),
	}

	res := &LockResponse{
		Lock: lock,
	}

	return res
}

func parseUnlockRequest(ctx Context) (*usecase.UnlockRequest, error) {

	data, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}

	var req usecase.UnlockRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	req.Repo = ctx.GetParam("repo")
	req.User = ctx.GetParam("user")
	req.ID = ctx.GetParam("id")

	return &req, nil
}

func convertUnlockResponce(result *usecase.LockResult) *UnlockResponse {

	lock := &Lock{
		ID:       result.ID,
		Path:     result.Path,
		Owner:    User{Name: result.Owner},
		LockedAt: time.Unix(result.LockedAt, 0),
	}

	res := &UnlockResponse{
		Lock: lock,
	}

	return res
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

func parseVerifyRequest(ctx Context) (*usecase.LockVerifyRequest, error) {

	data, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}

	var req usecase.LockVerifyRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	req.Repo = ctx.GetParam("repo")
	req.User = ctx.GetParam("user")
	req.Path = ctx.GetParam("path")

	return &req, nil
}

func convertVerifyResponse(result *usecase.LockVerifyResult) *LockVerifyResponse {

	var ours []Lock
	var theirs []Lock

	for _, l := range result.Ours {

		lock := Lock{
			ID:       l.ID,
			Path:     l.Path,
			Owner:    User{Name: l.Owner},
			LockedAt: time.Unix(l.LockedAt, 0),
		}

		ours = append(ours, lock)
	}

	for _, l := range result.Theirs {

		lock := Lock{
			ID:       l.ID,
			Path:     l.Path,
			Owner:    User{Name: l.Owner},
			LockedAt: time.Unix(l.LockedAt, 0),
		}

		theirs = append(theirs, lock)
	}

	res := &LockVerifyResponse{
		Ours:       ours,
		Theirs:     theirs,
		NextCursor: result.NextCursor,
	}

	return res
}
