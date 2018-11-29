package adapter

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/usecase"
)

const (
	testDBFile            = "test-database.db"
	testUser1             = "bilbo1"
	testPass1             = "baggins1"
	testUser2             = "bilbo2"
	testPass2             = "baggins2"
	testRepo              = "repo"
	testContent           = "this is my content"
	testContentSize       = int64(len(testContent))
	testContentOid        = "f97e1b2936a56511b3b6efc99011758e4700d60fb1674d31445d1ee40b663f24"
	testNonExistingOid    = "aec070645fe53ee3b3763059376134f058cc337247c978add178b6ccdfb0019f"
	testLockID            = "3cfec93346f7ff337c60f2da50cd86740715e2f6"
	testNonExistingLockID = "f310c1555a2485e2e5229ea015a94c9d590763d3"
	testLockPath          = "this/is/lock/path"
)

var (
	testDB                 *bolt.DB
	testMetaDataRepository usecase.MetaDataRepository
	testLockRepository     usecase.LockRepository
	testUserRepository     usecase.UserRepository
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func setupRepository() {

	testDB, err := bolt.Open(testDBFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Printf("error initializing test meta store: %s\n", err)
		os.Exit(1)
	}

	testDB.Update(func(tx *bolt.Tx) error {

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

	testMetaDataRepository = NewMetaDataRepository(testDB)
	testLockRepository = NewLockRepository(testDB)
	testUserRepository = NewUserRepository(testDB)

	if err := testUserRepository.AddUser(testUser1, testPass1); err != nil {
		teardownRepository()
		fmt.Printf("error adding test user to meta store: %s\n", err)
		os.Exit(1)
	}

	o := &usecase.ObjectRequest{Oid: testContentOid, Size: testContentSize}
	if _, err := testMetaDataRepository.Put(o); err != nil {
		teardownRepository()
		fmt.Printf("error seeding test meta store: %s\n", err)
		os.Exit(1)
	}
}

func teardownRepository() {
	if testDB != nil {
		testDB.Close()
	}
	os.RemoveAll(testDBFile)
}
