package usecase

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/ikmski/git-lfs3/entity"
)

// LockService is ...
type LockService interface {
	Lock(req *LockRequest) (*LockResult, error)
	Unlock(req *UnlockRequest) (*LockResult, error)
	List(req *LockListRequest) (*LockListResult, error)
	Verify(req *LockVerifyRequest) (*LockVerifyResult, error)
}

type lockService struct {
	LockRepository LockRepository
}

func NewLockService(lockRepo LockRepository) LockService {
	return &lockService{
		LockRepository: lockRepo,
	}
}

func (s *lockService) Lock(req *LockRequest) (*LockResult, error) {

	locks, _, err := s.LockRepository.FilteredFetch(req.Repo, req.Path, "", "1")
	if err != nil {
		return nil, err
	}

	if len(locks) > 0 {
		lock := locks[0]
		result := &LockResult{
			ID:           lock.ID,
			Path:         lock.Path,
			Owner:        lock.Owner.Name,
			LockedAt:     lock.LockedAt,
			AlreadyExist: true,
		}

		return result, nil
	}

	lock := entity.Lock{
		ID:       randomLockId(),
		Path:     req.Path,
		Owner:    entity.User{Name: req.User},
		LockedAt: time.Now().Unix(),
	}

	err = s.LockRepository.Add(req.Repo, lock)
	if err != nil {
		return nil, err
	}

	result := &LockResult{
		ID:           lock.ID,
		Path:         lock.Path,
		Owner:        lock.Owner.Name,
		LockedAt:     lock.LockedAt,
		AlreadyExist: false,
	}

	return result, nil
}

func (s *lockService) Unlock(req *UnlockRequest) (*LockResult, error) {

	lock, err := s.LockRepository.Delete(req.Repo, req.User, req.ID, req.Force)
	if err != nil {
		return nil, err
	}
	if lock == nil {
		// TODO
		return nil, nil
	}

	result := &LockResult{
		ID:           lock.ID,
		Path:         lock.Path,
		Owner:        lock.Owner.Name,
		LockedAt:     lock.LockedAt,
		AlreadyExist: true,
	}

	return result, nil
}

func (s *lockService) List(req *LockListRequest) (*LockListResult, error) {

	limit := ""
	if req.Limit > 0 {
		limit = strconv.Itoa(req.Limit)
	}

	locks, next, err := s.LockRepository.FilteredFetch(req.Repo, req.Path, req.Cursor, limit)
	if err != nil {
		return nil, err
	}

	listResult := &LockListResult{
		NextCursor: next,
	}

	for _, lock := range locks {
		l := &LockResult{
			ID:           lock.ID,
			Path:         lock.Path,
			Owner:        lock.Owner.Name,
			LockedAt:     lock.LockedAt,
			AlreadyExist: true,
		}
		listResult.Locks = append(listResult.Locks, l)
	}

	return listResult, nil
}

func (s *lockService) Verify(req *LockVerifyRequest) (*LockVerifyResult, error) {

	locks, next, err := s.LockRepository.FilteredFetch(req.Repo, req.Path, req.Cursor, strconv.Itoa(req.Limit))
	if err != nil {
		return nil, err
	}

	result := &LockVerifyResult{
		NextCursor: next,
	}

	for _, lock := range locks {
		l := &LockResult{
			ID:           lock.ID,
			Path:         lock.Path,
			Owner:        lock.Owner.Name,
			LockedAt:     lock.LockedAt,
			AlreadyExist: true,
		}

		if lock.Owner.Name == req.User {
			result.Ours = append(result.Ours, l)
		} else {
			result.Theirs = append(result.Theirs, l)
		}
	}

	return result, nil
}

func randomLockId() string {
	var id [20]byte
	rand.Read(id[:])
	return fmt.Sprintf("%x", id[:])
}
