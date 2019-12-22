package adapter

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/entity"
	"github.com/ikmski/git-lfs3/usecase"
)

var (
	objectsBucket = []byte("objects")
)

type metaDataRepository struct {
	db *bolt.DB
}

// NewMetaDataRepository is ...
func NewMetaDataRepository(db *bolt.DB) usecase.MetaDataRepository {

	return &metaDataRepository{db: db}
}

// Get retrieves the Meta information for an object given information in
// Object
func (s *metaDataRepository) Get(o *usecase.ObjectRequest) (*entity.MetaData, error) {

	meta, error := s.UnsafeGet(o)
	return meta, error
}

// Get retrieves the Meta information for an object given information in
// Object
// DO NOT CHECK authentication, as it is supposed to have been done before
func (s *metaDataRepository) UnsafeGet(o *usecase.ObjectRequest) (*entity.MetaData, error) {

	var meta entity.MetaData

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(objectsBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
		}

		value := bucket.Get([]byte(o.Oid))
		if len(value) == 0 {
			return errors.New("Object not found")
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
func (s *metaDataRepository) Put(o *usecase.ObjectRequest) (*entity.MetaData, error) {

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
			return errors.New("Bucket not found")
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
func (s *metaDataRepository) Delete(o *usecase.ObjectRequest) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(objectsBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
		}

		err := bucket.Delete([]byte(o.Oid))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// Objects returns all MetaObjects in the meta store
func (s *metaDataRepository) Objects() ([]*entity.MetaData, error) {

	var objects []*entity.MetaData

	err := s.db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(objectsBucket)
		if bucket == nil {
			return errors.New("Bucket not found")
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
