package main

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/adapter"
	"github.com/ikmski/git-lfs3/entity"
	"github.com/ikmski/git-lfs3/usecase"
)

var (
	lfsServer        *httptest.Server
	testMetaDataRepo usecase.MetaDataRepository
	testContentRepo  usecase.ContentRepository
)

const (
	testUser1             = "bilbo1"
	testPass1             = "baggins1"
	testUser2             = "bilbo2"
	testPass2             = "baggins2"
	testRepo              = "repo"
	testContent           = "this is my content"
	testContentSize       = int64(len(testContent))
	testContentOid        = "f97e1b2936a56511b3b6efc99011758e4700d60fb1674d31445d1ee40b663f24"
	testNonExistingOid    = "aec070645fe53ee3b3763059376134f058cc337247c978add178b6ccdfb0019f"
	testLockId            = "3cfec93346f7ff337c60f2da50cd86740715e2f6"
	testNonExistingLockId = "f310c1555a2485e2e5229ea015a94c9d590763d3"
	testLockPath          = "this/is/lock/path"
)

func TestMain(m *testing.M) {

	os.Remove("lfs-test.db")

	var err error

	db, err := bolt.Open("lfs-test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Printf("Error creating test db: %s", err)
		os.Exit(1)
	}

	usersBucket := []byte("users")
	objectsBucket := []byte("objects")
	locksBucket := []byte("locks")

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

	testMetaDataRepo = adapter.NewMetaDataRepository(db)
	if err != nil {
		fmt.Printf("Error creating meta store: %s", err)
		os.Exit(1)
	}

	testContentRepo, err = adapter.NewMockedContentRepository("lfs-test-bucket")
	if err != nil {
		fmt.Printf("Error creating content store: %s", err)
		os.Exit(1)
	}

	err = seedContentRepository()
	if err != nil {
		fmt.Printf("Error seeding content store: %s", err)
		os.Exit(1)
	}

	err = seedMetaDataRepository()
	if err != nil {
		fmt.Printf("Error seeding meta store: %s", err)
		os.Exit(1)
	}

	conf := serverConfig{}
	batchService := usecase.NewBatchService(testMetaDataRepo, testContentRepo)
	transferService := usecase.NewTransferService(testMetaDataRepo, testContentRepo)

	batchController := adapter.NewBatchController(batchService)
	transferController := adapter.NewTransferController(transferService)

	app := newApp(conf, batchController, transferController)
	lfsServer = httptest.NewServer(app)

	ret := m.Run()

	lfsServer.Close()
	db.Close()

	os.Remove("lfs-test.db")

	os.Exit(ret)
}

func seedMetaDataRepository() error {

	/*
		if err := testMetaDataRepo.AddUser(testUser1, testPass1); err != nil {
			return err
		}
		if err := testMetaDataRepo.AddUser(testUser2, testPass2); err != nil {
			return err
		}
	*/

	o := &usecase.ObjectRequest{
		Oid:  testContentOid,
		Size: testContentSize,
	}

	_, err := testMetaDataRepo.Put(o)
	if err != nil {
		return err
	}

	/*
		lock := NewTestLock(testLockId, testLockPath, testUser1)
		if err := testMetaStore.AddLocks(testRepo, lock); err != nil {
			return err
		}
	*/

	return nil
}

func seedContentRepository() error {

	meta := &entity.MetaData{
		Oid:  testContentOid,
		Size: testContentSize,
	}

	buf := bytes.NewBuffer([]byte(testContent))

	err := testContentRepo.Put(meta, buf)
	if err != nil {
		return err
	}

	return nil
}
