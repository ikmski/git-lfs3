package adapter

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/entity"
	"github.com/ikmski/git-lfs3/usecase"
)

var (
	usersBucket = []byte("users")
)

type userRepository struct {
	db *bolt.DB
}

// NewUserRepository is ...
func NewUserRepository(db *bolt.DB) usecase.UserRepository {
	return &userRepository{db: db}
}

// AddUser adds user credentials to the meta store.
func (s *userRepository) AddUser(user, pass string) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
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
func (s *userRepository) DeleteUser(user string) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
		}

		err := bucket.Delete([]byte(user))
		return err
	})

	return err
}

// Users returns all MetaUsers in the meta store
func (s *userRepository) Users() ([]*entity.User, error) {

	var users []*entity.User

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errNoBucket
		}

		bucket.ForEach(func(k, v []byte) error {
			users = append(users, &entity.User{string(k)})
			return nil
		})
		return nil
	})

	return users, err
}
