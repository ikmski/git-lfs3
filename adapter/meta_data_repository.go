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

	db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(objectsBucket)
		if err != nil {
			return err
		}
		return nil

	})

	return &metaDataRepository{db: db}
}

// Get retrieves the Meta information for an object given information in
// Object
func (r *metaDataRepository) Get(o *usecase.ObjectRequest) (*entity.MetaData, error) {

	meta, error := r.UnsafeGet(o)
	return meta, error
}

// Get retrieves the Meta information for an object given information in
// Object
// DO NOT CHECK authentication, as it is supposed to have been done before
func (r *metaDataRepository) UnsafeGet(o *usecase.ObjectRequest) (*entity.MetaData, error) {

	var meta entity.MetaData

	err := r.db.View(func(tx *bolt.Tx) error {

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
func (r *metaDataRepository) Put(o *usecase.ObjectRequest) (*entity.MetaData, error) {

	// Check if it exists first
	meta, err := r.Get(o)
	if err == nil {
		return meta, nil
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	meta = &entity.MetaData{
		Oid:  o.Oid,
		Size: o.Size,
	}

	err = enc.Encode(meta)
	if err != nil {
		return nil, err
	}

	err = r.db.Update(func(tx *bolt.Tx) error {

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

	return meta, nil
}

// Delete removes the meta information from Object to the store.
func (r *metaDataRepository) Delete(o *usecase.ObjectRequest) error {

	err := r.db.Update(func(tx *bolt.Tx) error {

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
func (r *metaDataRepository) Objects() ([]*entity.MetaData, error) {

	var objects []*entity.MetaData

	err := r.db.View(func(tx *bolt.Tx) error {

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
