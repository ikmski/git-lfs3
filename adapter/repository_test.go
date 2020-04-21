package adapter

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func setupRepository(d *TestData) {

	db, err := bolt.Open(d.databaseFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Printf("error initializing test meta store: %s\n", err)
		os.Exit(1)
	}

	db.Update(func(tx *bolt.Tx) error {

		if _, err := tx.CreateBucketIfNotExists(usersBucket); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists(metaBucket); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists(locksBucket); err != nil {
			return err
		}

		return nil
	})

	d.metaDataRepository = NewMetaDataRepository(db)
	d.lockRepository = NewLockRepository(db)
	d.userRepository = NewUserRepository(db)
	d.database = db

	if err := d.userRepository.AddUser(d.userName1, d.userPass1); err != nil {
		teardownRepository(d)
		fmt.Printf("error adding test user to meta store: %s\n", err)
		os.Exit(1)
	}

	if _, err := d.metaDataRepository.Put(d.contentOid, d.contentSize); err != nil {
		teardownRepository(d)
		fmt.Printf("error seeding test meta store: %s\n", err)
		os.Exit(1)
	}
}

func teardownRepository(d *TestData) {
	if d.database != nil {
		d.database.Close()
	}
	os.RemoveAll(d.databaseFile)
}
