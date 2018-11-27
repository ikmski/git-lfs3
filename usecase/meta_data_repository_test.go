package usecase

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ikmski/git-lfs3/entity"
)

var (
	testMetaDataRepository MetaDataRepository
)

func TestGetMeta(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	meta, err := testMetaDataRepository.Get(&ObjectRequest{Oid: testContentOid})
	if err != nil {
		t.Fatalf("Error retreiving meta: %s", err)
	}

	if meta.Oid != testContentOid {
		t.Errorf("expected to get content oid, got: %s", meta.Oid)
	}

	if meta.Size != testContentSize {
		t.Errorf("expected to get content size, got: %d", meta.Size)
	}
}

func TestPutMeta(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	meta, err := testMetaDataRepository.Put(&ObjectRequest{Oid: testNonExistingOid, Size: 42})
	if err != nil {
		t.Errorf("expected put to succeed, got : %s", err)
	}

	if meta.Existing {
		t.Errorf("expected meta to not have existed")
	}

	meta, err = testMetaDataRepository.Get(&ObjectRequest{Oid: testNonExistingOid})
	if err != nil {
		t.Errorf("expected to be able to retreive new put, got : %s", err)
	}

	if meta.Oid != testNonExistingOid {
		t.Errorf("expected oids to match, got: %s", meta.Oid)
	}

	if meta.Size != 42 {
		t.Errorf("expected sizes to match, got: %d", meta.Size)
	}

	meta, err = testMetaDataRepository.Put(&ObjectRequest{Oid: testNonExistingOid, Size: 42})
	if err != nil {
		t.Errorf("expected put to succeed, got : %s", err)
	}

	if !meta.Existing {
		t.Errorf("expected meta to now exist")
	}
}

func TestLocks(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	for i := 0; i < 5; i++ {
		lock := NewTestLock(randomLockId(), fmt.Sprintf("path-%d", i), fmt.Sprintf("user-%d", i))
		if err := testMetaDataRepository.AddLocks(testRepo, lock); err != nil {
			t.Errorf("expected AddLocks to succeed, got : %s", err)
		}
	}

	locks, err := testMetaDataRepository.Locks(testRepo)
	if err != nil {
		t.Errorf("expected Locks to succeed, got : %s", err)
	}
	if len(locks) != 5 {
		t.Errorf("expected returned lock count to match, got: %d", len(locks))
	}
}

func TestFilteredLocks(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	testLocks := make([]entity.Lock, 0, 5)
	for i := 0; i < 5; i++ {
		lock := NewTestLock(randomLockId(), fmt.Sprintf("path-%d", i), fmt.Sprintf("user-%d", i))
		testLocks = append(testLocks, lock)
	}
	if err := testMetaDataRepository.AddLocks(testRepo, testLocks...); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	locks, next, err := testMetaDataRepository.FilteredLocks(testRepo, "", "", "3")
	if err != nil {
		t.Errorf("expected FilteredLocks to succeed, got : %s", err)
	}
	if len(locks) != 3 {
		t.Errorf("expected locks count to match limit, got: %d", len(locks))
	}
	if next == "" {
		t.Errorf("expected next to exist")
	}

	locks, next, err = testMetaDataRepository.FilteredLocks(testRepo, "", next, "2")
	if err != nil {
		t.Errorf("expected FilteredLocks to succeed, got : %s", err)
	}
	if len(locks) != 2 {
		t.Errorf("expected locks count to match limit, got: %d", len(locks))
	}
	if next != "" {
		t.Errorf("expected next to not exist, got: %s", next)
	}
}

func TestAddLocks(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	lock := NewTestLock(testLockId, testLockPath, testUser1)
	if err := testMetaDataRepository.AddLocks(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	locks, _, err := testMetaDataRepository.FilteredLocks(testRepo, lock.Path, "", "1")
	if err != nil {
		t.Errorf("expected FilteredLocks to succeed, got : %s", err)
	}
	if len(locks) != 1 {
		t.Errorf("expected lock to be existed")
	}
	if locks[0].ID != testLockId {
		t.Errorf("expected lockId to match, got: %s", locks[0])
	}
}

func TestDeleteLock(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	lock := NewTestLock(testLockId, testLockPath, testUser1)
	if err := testMetaDataRepository.AddLocks(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	deleted, err := testMetaDataRepository.DeleteLock(testRepo, testUser1, lock.ID, false)
	if err != nil {
		t.Errorf("expected DeleteLock to succeed, got : %s", err)
	}
	if deleted == nil || deleted.ID != lock.ID {
		t.Errorf("expected deleted lock to be returned, got : %s", deleted)
	}
}

func TestDeleteLockNotOwner(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	lock := NewTestLock(testLockId, testLockPath, testUser1)
	if err := testMetaDataRepository.AddLocks(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	deleted, err := testMetaDataRepository.DeleteLock(testRepo, testUser2, lock.ID, false)
	if err == nil || deleted != nil {
		t.Errorf("expected DeleteLock to failed")
	}

	if err != errNotOwner {
		t.Errorf("expected DeleteLock error match, got: %s", err)
	}
}

func TestDeleteLockNotOwnerForce(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	lock := NewTestLock(testLockId, testLockPath, testUser1)
	if err := testMetaDataRepository.AddLocks(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	deleted, err := testMetaDataRepository.DeleteLock(testRepo, testUser2, lock.ID, true)
	if err != nil {
		t.Errorf("expected DeleteLock(force) to succeed, got : %s", err)
	}
	if deleted == nil || deleted.ID != lock.ID {
		t.Errorf("expected deleted lock to be returned, got : %s", deleted)
	}
}

func TestDeleteLockNonExisting(t *testing.T) {

	setupMeta()
	defer teardownMeta()

	lock := NewTestLock(testLockId, testLockPath, testUser1)
	if err := testMetaDataRepository.AddLocks(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	deleted, err := testMetaDataRepository.DeleteLock(testRepo, testUser1, testNonExistingLockId, false)
	if err != nil {
		t.Errorf("expected DeleteLock to succeed, got : %s", err)
	}
	if deleted != nil {
		t.Errorf("expected nil returned, got : %s", deleted)
	}
}

func NewTestLock(id, path, user string) entity.Lock {

	return entity.Lock{
		ID:   id,
		Path: path,
		Owner: entity.User{
			Name: user,
		},
		LockedAt: time.Now(),
	}
}

func setupMeta() {

	store, err := NewMetaDataRepository("test-meta-store.db")
	if err != nil {
		fmt.Printf("error initializing test meta store: %s\n", err)
		os.Exit(1)
	}

	testMetaDataRepository = store
	if err := testMetaDataRepository.AddUser(testUser1, testPass1); err != nil {
		teardownMeta()
		fmt.Printf("error adding test user to meta store: %s\n", err)
		os.Exit(1)
	}

	o := &ObjectRequest{Oid: testContentOid, Size: testContentSize}
	if _, err := testMetaDataRepository.Put(o); err != nil {
		teardownMeta()
		fmt.Printf("error seeding test meta store: %s\n", err)
		os.Exit(1)
	}
}

func teardownMeta() {
	testMetaDataRepository.Close()
	os.RemoveAll("test-meta-store.db")
}

func randomLockId() string {

	var id [20]byte
	rand.Read(id[:])
	return fmt.Sprintf("%x", id[:])
}
