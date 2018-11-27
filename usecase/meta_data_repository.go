package usecase

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/entity"
)

var (
	errNoBucket       = errors.New("Bucket not found")
	errObjectNotFound = errors.New("Object not found")
	errNotOwner       = errors.New("Attempt to delete other user's lock")
)

var (
	usersBucket   = []byte("users")
	objectsBucket = []byte("objects")
	locksBucket   = []byte("locks")
)

// MetaDataRepository is ...
type MetaDataRepository interface {
	Close()

	Get(o *ObjectRequest) (*entity.MetaData, error)
	Put(o *ObjectRequest) (*entity.MetaData, error)
	Delete(o *ObjectRequest) error

	AddLocks(repo string, l ...entity.Lock) error
	Locks(repo string) ([]entity.Lock, error)
	FilteredLocks(repo, path, cursor, limit string) (locks []entity.Lock, next string, err error)
	DeleteLock(repo, user, id string, force bool) (*entity.Lock, error)

	AddUser(user, pass string) error
	DeleteUser(user string) error
	Users() ([]*MetaUser, error)

	Objects() ([]*entity.MetaData, error)
	AllLocks() ([]entity.Lock, error)
}

type metaDataRepository struct {
	db *bolt.DB
}

// NewMetaDataRepository is ...
func NewMetaDataRepository(dbFile string) (MetaDataRepository, error) {

	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {

		if _, err := tx.CreateBucketIfNotExists(usersBucket); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists(objectsBucket); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists(locksBucket); err != nil {
			return err
		}

		return nil
	})

	return &metaDataRepository{db: db}, nil
}

// Get retrieves the Meta information for an object given information in
// Object
func (s *metaDataRepository) Get(o *ObjectRequest) (*entity.MetaData, error) {

	meta, error := s.UnsafeGet(o)
	return meta, error
}

// Get retrieves the Meta information for an object given information in
// Object
// DO NOT CHECK authentication, as it is supposed to have been done before
func (s *metaDataRepository) UnsafeGet(o *ObjectRequest) (*entity.MetaData, error) {

	var meta entity.MetaData

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(objectsBucket)
		if bucket == nil {
			return errNoBucket
		}

		value := bucket.Get([]byte(o.Oid))
		if len(value) == 0 {
			return errObjectNotFound
		}

		dec := gob.NewDecoder(bytes.NewBuffer(value))
		return dec.Decode(&meta)
	})

	if err != nil {
		return nil, err
	}

	return &meta, nil
}

// Put writes meta information from Object to the store.
func (s *metaDataRepository) Put(o *ObjectRequest) (*entity.MetaData, error) {

	// Check if it exists first
	if meta, err := s.Get(o); err == nil {
		meta.Existing = true
		return meta, nil
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	meta := entity.MetaData{Oid: o.Oid, Size: o.Size}
	err := enc.Encode(meta)
	if err != nil {
		return nil, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(objectsBucket)
		if bucket == nil {
			return errNoBucket
		}

		err = bucket.Put([]byte(o.Oid), buf.Bytes())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &meta, nil
}

// Delete removes the meta information from Object to the store.
func (s *metaDataRepository) Delete(o *ObjectRequest) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(objectsBucket)
		if bucket == nil {
			return errNoBucket
		}

		err := bucket.Delete([]byte(o.Oid))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// AddLocks write locks to the store for the repo.
func (s *metaDataRepository) AddLocks(repo string, l ...entity.Lock) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(locksBucket)
		if bucket == nil {
			return errNoBucket
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
func (s *metaDataRepository) Locks(repo string) ([]entity.Lock, error) {

	var locks []entity.Lock
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(locksBucket)
		if bucket == nil {
			return errNoBucket
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
func (s *metaDataRepository) FilteredLocks(repo, path, cursor, limit string) (locks []entity.Lock, next string, err error) {

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
func (s *metaDataRepository) DeleteLock(repo, user, id string, force bool) (*entity.Lock, error) {

	var deleted *entity.Lock
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(locksBucket)
		if bucket == nil {
			return errNoBucket
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

// Close closes the underlying boltdb.
func (s *metaDataRepository) Close() {
	s.db.Close()
}

// AddUser adds user credentials to the meta store.
func (s *metaDataRepository) AddUser(user, pass string) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errNoBucket
		}

		err := bucket.Put([]byte(user), []byte(pass))
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// DeleteUser removes user credentials from the meta store.
func (s *metaDataRepository) DeleteUser(user string) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errNoBucket
		}

		err := bucket.Delete([]byte(user))
		return err
	})

	return err
}

// MetaUser encapsulates information about a meta store user
type MetaUser struct {
	Name string
}

// Users returns all MetaUsers in the meta store
func (s *metaDataRepository) Users() ([]*MetaUser, error) {

	var users []*MetaUser

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errNoBucket
		}

		bucket.ForEach(func(k, v []byte) error {
			users = append(users, &MetaUser{string(k)})
			return nil
		})
		return nil
	})

	return users, err
}

// Objects returns all MetaObjects in the meta store
func (s *metaDataRepository) Objects() ([]*entity.MetaData, error) {

	var objects []*entity.MetaData

	err := s.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(objectsBucket)
		if bucket == nil {
			return errNoBucket
		}

		bucket.ForEach(func(k, v []byte) error {
			var meta entity.MetaData
			dec := gob.NewDecoder(bytes.NewBuffer(v))
			err := dec.Decode(&meta)
			if err != nil {
				return err
			}

			objects = append(objects, &meta)
			return nil
		})
		return nil
	})

	return objects, err
}

// AllLocks return all locks in the store, lock path is prepended with repo
func (s *metaDataRepository) AllLocks() ([]entity.Lock, error) {

	var locks []entity.Lock
	err := s.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(locksBucket)
		if bucket == nil {
			return errNoBucket
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
