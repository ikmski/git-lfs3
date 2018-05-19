package main

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	lfsServer        *httptest.Server
	testMetaStore    *MetaStore
	testContentStore *ContentStore
)

const (
	testUser1         = "bilbo1"
	testPass1         = "baggins1"
	testUser2         = "bilbo2"
	testPass2         = "baggins2"
	testRepo          = "repo"
	content           = "this is my content"
	contentSize       = int64(len(content))
	contentOid        = "f97e1b2936a56511b3b6efc99011758e4700d60fb1674d31445d1ee40b663f24"
	nonExistingOid    = "aec070645fe53ee3b3763059376134f058cc337247c978add178b6ccdfb0019f"
	lockId            = "3cfec93346f7ff337c60f2da50cd86740715e2f6"
	nonExistingLockId = "f310c1555a2485e2e5229ea015a94c9d590763d3"
	lockPath          = "this/is/lock/path"
)

func TestMain(m *testing.M) {

	os.Remove("lfs-test.db")

	var err error
	testMetaStore, err = NewMetaStore("lfs-test.db")
	if err != nil {
		fmt.Printf("Error creating meta store: %s", err)
		os.Exit(1)
	}

	testContentStore, err = NewMockedContentStore("lfs-test-bucket")
	if err != nil {
		fmt.Printf("Error creating content store: %s", err)
		os.Exit(1)
	}

	if err := seedContentStore(); err != nil {
		fmt.Printf("Error seeding content store: %s", err)
		os.Exit(1)
	}

	if err := seedMetaStore(); err != nil {
		fmt.Printf("Error seeding meta store: %s", err)
		os.Exit(1)
	}

	app := newApp(testMetaStore, testContentStore)
	lfsServer = httptest.NewServer(app)

	ret := m.Run()

	lfsServer.Close()
	testMetaStore.Close()

	os.Remove("lfs-test.db")

	os.Exit(ret)
}

func seedMetaStore() error {

	if err := testMetaStore.AddUser(testUser1, testPass1); err != nil {
		return err
	}
	if err := testMetaStore.AddUser(testUser2, testPass2); err != nil {
		return err
	}

	o := &Object{
		Oid:  contentOid,
		Size: contentSize,
	}

	_, err := testMetaStore.Put(o)
	if err != nil {
		return err
	}

	lock := NewTestLock(lockId, lockPath, testUser1)
	if err := testMetaStore.AddLocks(testRepo, lock); err != nil {
		return err
	}

	return nil
}

func seedContentStore() error {

	meta := &MetaObject{
		Oid:  contentOid,
		Size: contentSize,
	}

	buf := bytes.NewBuffer([]byte(content))

	err := testContentStore.Put(meta, buf)
	if err != nil {
		return err
	}

	return nil
}
