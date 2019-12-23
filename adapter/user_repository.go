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

	db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(usersBucket)
		if err != nil {
			return err
		}
		return nil

	})

	return &userRepository{db: db}
}

// AddUser adds user credentials to the meta store.
func (r *userRepository) AddUser(user, pass string) error {

	err := r.db.Update(func(tx *bolt.Tx) error {

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
func (r *userRepository) DeleteUser(user string) error {

	err := r.db.Update(func(tx *bolt.Tx) error {

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
func (r *userRepository) Users() ([]*entity.User, error) {

	var users []*entity.User

	err := r.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
		}

		bucket.ForEach(func(k, v []byte) error {
			u := &entity.User{
				Name: string(k),
			}
			users = append(users, u)
			return nil
		})
		return nil
	})

	return users, err
}
