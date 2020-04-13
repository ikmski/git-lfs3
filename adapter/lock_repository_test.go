package adapter

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/ikmski/git-lfs3/entity"
)

func TestLocks(t *testing.T) {

	setupRepository()
	defer teardownRepository()

	for i := 0; i < 5; i++ {
		lock := NewTestLock(randomLockId(), fmt.Sprintf("path-%d", i), fmt.Sprintf("user-%d", i))
		if err := testLockRepository.Add(testRepo, lock); err != nil {
			t.Errorf("expected AddLocks to succeed, got : %s", err)
		}
	}

	locks, err := testLockRepository.Fetch(testRepo)
	if err != nil {
		t.Errorf("expected Locks to succeed, got : %s", err)
	}
	if len(locks) != 5 {
		t.Errorf("expected returned lock count to match, got: %d", len(locks))
	}
}

func TestFilteredLocks(t *testing.T) {

	setupRepository()
	defer teardownRepository()

	testLocks := make([]entity.Lock, 0, 5)
	for i := 0; i < 5; i++ {
		lock := NewTestLock(randomLockId(), fmt.Sprintf("path-%d", i), fmt.Sprintf("user-%d", i))
		testLocks = append(testLocks, lock)
	}
	if err := testLockRepository.Add(testRepo, testLocks...); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	locks, next, err := testLockRepository.FilteredFetch(testRepo, "", "", "3")
	if err != nil {
		t.Errorf("expected FilteredLocks to succeed, got : %s", err)
	}
	if len(locks) != 3 {
		t.Errorf("expected locks count to match limit, got: %d", len(locks))
	}
	if next == "" {
		t.Errorf("expected next to exist")
	}

	locks, next, err = testLockRepository.FilteredFetch(testRepo, "", next, "2")
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

	setupRepository()
	defer teardownRepository()

	lock := NewTestLock(testLockID, testLockPath, testUser1)
	if err := testLockRepository.Add(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	locks, _, err := testLockRepository.FilteredFetch(testRepo, lock.Path, "", "1")
	if err != nil {
		t.Errorf("expected FilteredLocks to succeed, got : %s", err)
	}
	if len(locks) != 1 {
		t.Errorf("expected lock to be existed")
	}
	if locks[0].ID != testLockID {
		t.Errorf("expected lockId to match, got: %v", locks[0])
	}
}

func TestDeleteLock(t *testing.T) {

	setupRepository()
	defer teardownRepository()

	lock := NewTestLock(testLockID, testLockPath, testUser1)
	if err := testLockRepository.Add(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	deleted, err := testLockRepository.Delete(testRepo, testUser1, lock.ID, false)
	if err != nil {
		t.Errorf("expected DeleteLock to succeed, got : %s", err)
	}
	if deleted == nil || deleted.ID != lock.ID {
		t.Errorf("expected deleted lock to be returned, got : %v", deleted)
	}
}

func TestDeleteLockNotOwner(t *testing.T) {

	setupRepository()
	defer teardownRepository()

	lock := NewTestLock(testLockID, testLockPath, testUser1)
	if err := testLockRepository.Add(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	deleted, err := testLockRepository.Delete(testRepo, testUser2, lock.ID, false)
	if err == nil || deleted != nil {
		t.Errorf("expected DeleteLock to failed")
	}

	if err != errNotOwner {
		t.Errorf("expected DeleteLock error match, got: %s", err)
	}
}

func TestDeleteLockNotOwnerForce(t *testing.T) {

	setupRepository()
	defer teardownRepository()

	lock := NewTestLock(testLockID, testLockPath, testUser1)
	if err := testLockRepository.Add(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	deleted, err := testLockRepository.Delete(testRepo, testUser2, lock.ID, true)
	if err != nil {
		t.Errorf("expected DeleteLock(force) to succeed, got : %s", err)
	}
	if deleted == nil || deleted.ID != lock.ID {
		t.Errorf("expected deleted lock to be returned, got : %v", deleted)
	}
}

func TestDeleteLockNonExisting(t *testing.T) {

	setupRepository()
	defer teardownRepository()

	lock := NewTestLock(testLockID, testLockPath, testUser1)
	if err := testLockRepository.Add(testRepo, lock); err != nil {
		t.Errorf("expected AddLocks to succeed, got : %s", err)
	}

	deleted, err := testLockRepository.Delete(testRepo, testUser1, testNonExistingLockID, false)
	if err != nil {
		t.Errorf("expected DeleteLock to succeed, got : %s", err)
	}
	if deleted != nil {
		t.Errorf("expected nil returned, got : %v", deleted)
	}
}

func NewTestLock(id, path, user string) entity.Lock {

	return entity.Lock{
		ID:   id,
		Path: path,
		Owner: entity.User{
			Name: user,
		},
		LockedAt: time.Now().Unix(),
	}
}

func randomLockId() string {

	var id [20]byte
	rand.Read(id[:])
	return fmt.Sprintf("%x", id[:])
}
