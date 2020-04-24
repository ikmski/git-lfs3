package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikmski/git-lfs3/adapter"
)

func TestLocksList(t *testing.T) {

	path := fmt.Sprintf("%s/%s/%s/locks", lfsServer.URL, testUser1, testRepo)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}
	req.Header.Set("Accept", "application/vnd.git-lfs+json")
	req.Header.Set("Content-Type", "application/vnd.git-lfs+json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected response to contain content, got error: %s", err)
	}

	var list adapter.LockListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatalf("expected response body to be LockList, got error: %s", err)
	}
	if len(list.Locks) != 1 {
		t.Errorf("expected returned lock count to match, got: %d", len(list.Locks))
	}
	if list.Locks[0].ID != testLockId {
		t.Errorf("expected lockId to match, got: %s", list.Locks[0].ID)
	}
}

func TestLocksVerify(t *testing.T) {

	path := fmt.Sprintf("%s/%s/%s/locks/verify", lfsServer.URL, testUser1, testRepo)
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}
	req.Header.Set("Accept", "application/vnd.git-lfs+json")
	req.Header.Set("Content-Type", "application/vnd.git-lfs+json")

	buf := bytes.NewBufferString(fmt.Sprintf(`{"cursor": "", "limit": 0}`))
	req.Body = ioutil.NopCloser(buf)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("response error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected response to contain content, got error: %s", err)
	}

	var list adapter.LockVerifyResponse
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatalf("expected response body to be VerifiableLockList, got error: %s", err)
	}
}

func TestLock(t *testing.T) {

	path := "TestLock"
	lock, err := addLock(testUser1, path)

	if err != nil {
		t.Fatalf("create lock error: %s", err)
	}
	if lock == nil {
		t.Errorf("expected lock to be created, got: %s", lock)
	}
	if lock.Owner.Name != testUser1 {
		t.Errorf("expected lock owner to be match, got: %s", lock.Owner.Name)
	}
	if lock.Path != path {
		t.Errorf("expected lock path to be match, got: %s", lock.Path)
	}
}

func TestUnlock(t *testing.T) {

	l, err := addLock(testUser1, "TestUnlock")
	if err != nil {
		t.Fatalf("create lock error: %s", err)
	}

	url := fmt.Sprintf("%s/%s/%s/locks/%s/unlock", lfsServer.URL, testUser1, testRepo, l.ID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}
	req.Header.Set("Accept", "application/vnd.git-lfs+json")
	req.Header.Set("Content-Type", "application/vnd.git-lfs+json")

	buf := bytes.NewBufferString(fmt.Sprintf(`{"force": %t}`, false))
	req.Body = ioutil.NopCloser(buf)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("response error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected response to contain content, got error: %s", err)
	}

	var unlockResponse adapter.UnlockResponse
	if err := json.Unmarshal(body, &unlockResponse); err != nil {
		t.Fatalf("expected response body to be UnlockResponse, got error: %s", err)
	}

	lock := unlockResponse.Lock
	if lock == nil || lock.ID != l.ID {
		t.Errorf("expected deleted lock to be returned, got: %s", lock)
	}
}

func addLock(username string, path string) (*adapter.Lock, error) {

	url := fmt.Sprintf("%s/%s/%s/locks", lfsServer.URL, username, testRepo)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.git-lfs+json")
	req.Header.Set("Content-Type", "application/vnd.git-lfs+json")

	buf := bytes.NewBufferString(fmt.Sprintf(`{"path":"%s"}`, path))
	req.Body = ioutil.NopCloser(buf)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected status 201, got %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("expected response to contain content, got error: %s", err)
	}

	var lockResponse adapter.LockResponse
	if err := json.Unmarshal(body, &lockResponse); err != nil {
		return nil, fmt.Errorf("expected response body to be LockResponse, got error: %s", err)
	}

	return lockResponse.Lock, nil
}
