package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/entity"
	"github.com/ikmski/git-lfs3/usecase"
)

var (
	locksBucket = []byte("locks")
	errNotOwner = errors.New("Attempt to delete other user's lock")
)

type lockRepository struct {
	db *bolt.DB
}

// NewLockRepository is ...
func NewLockRepository(db *bolt.DB) usecase.LockRepository {
	return &lockRepository{db: db}
}

// AddLocks write locks to the store for the repo.
func (s *lockRepository) AddLocks(repo string, l ...entity.Lock) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(locksBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
		}

		var locks []entity.Lock
		data := bucket.Get([]byte(repo))
		if data != nil {
			if err := json.Unmarshal(data, &locks); err != nil {
				return err
			}
		}
		locks = append(locks, l...)
		sort.Sort(LocksByCreatedAt(locks))
		data, err := json.Marshal(&locks)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(repo), data)
	})
	return err
}

// Locks retrieves locks for the repo from the store
func (s *lockRepository) Locks(repo string) ([]entity.Lock, error) {

	var locks []entity.Lock
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(locksBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
		}

		data := bucket.Get([]byte(repo))
		if data != nil {
			if err := json.Unmarshal(data, &locks); err != nil {
				return err
			}
		}
		return nil
	})
	return locks, err
}

// FilteredLocks return filtered locks for the repo
func (s *lockRepository) FilteredLocks(repo, path, cursor, limit string) (locks []entity.Lock, next string, err error) {

	locks, err = s.Locks(repo)
	if err != nil {
		return
	}

	if cursor != "" {
		lastSeen := -1
		for i, l := range locks {
			if l.ID == cursor {
				lastSeen = i
				break
			}
		}

		if lastSeen > -1 {
			locks = locks[lastSeen:]
		} else {
			err = fmt.Errorf("cursor (%s) not found", cursor)
			return
		}
	}

	if path != "" {
		var filtered []entity.Lock
		for _, l := range locks {
			if l.Path == path {
				filtered = append(filtered, l)
			}
		}

		locks = filtered
	}

	if limit != "" {
		var size int
		size, err = strconv.Atoi(limit)
		if err != nil || size < 0 {
			locks = make([]entity.Lock, 0)
			err = fmt.Errorf("Invalid limit amount: %s", limit)
			return
		}

		size = int(math.Min(float64(size), float64(len(locks))))
		if size+1 < len(locks) {
			next = locks[size].ID
		}
		locks = locks[:size]
	}

	return locks, next, nil
}

// DeleteLock removes lock for the repo by id from the store
func (s *lockRepository) DeleteLock(repo, user, id string, force bool) (*entity.Lock, error) {

	var deleted *entity.Lock
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(locksBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
		}

		var locks []entity.Lock
		data := bucket.Get([]byte(repo))
		if data != nil {
			if err := json.Unmarshal(data, &locks); err != nil {
				return err
			}
		}
		newLocks := make([]entity.Lock, 0, len(locks))

		var lock entity.Lock
		for _, l := range locks {
			if l.ID == id {
				if l.Owner.Name != user && !force {
					return errNotOwner
				}
				lock = l
			} else if len(l.ID) > 0 {
				newLocks = append(newLocks, l)
			}
		}
		if lock.ID == "" {
			return nil
		}
		deleted = &lock

		if len(newLocks) == 0 {
			return bucket.Delete([]byte(repo))
		}

		data, err := json.Marshal(&newLocks)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(repo), data)
	})
	return deleted, err
}

// LocksByCreatedAt is ...
type LocksByCreatedAt []entity.Lock

func (c LocksByCreatedAt) Len() int {
	return len(c)
}

func (c LocksByCreatedAt) Less(i, j int) bool {
	return c[i].LockedAt.Before(c[j].LockedAt)
}

func (c LocksByCreatedAt) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// AllLocks return all locks in the store, lock path is prepended with repo
func (s *lockRepository) AllLocks() ([]entity.Lock, error) {

	var locks []entity.Lock
	err := s.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(locksBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
		}

		bucket.ForEach(func(k, v []byte) error {

			var l []entity.Lock
			if err := json.Unmarshal(v, &l); err != nil {
				return err
			}

			for _, lv := range l {
				lv.Path = fmt.Sprintf("%s:%s", k, lv.Path)
				locks = append(locks, lv)
			}

			return nil
		})

		return nil
	})
	return locks, err
}
