package main

import (
	"fmt"
	"net/http"
	"testing"
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

	/*
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("expected response to contain content, got error: %s", err)
		}

		var list LockList
		if err := json.Unmarshal(body, &list); err != nil {
			t.Fatalf("expected response body to be LockList, got error: %s", err)
		}
		if len(list.Locks) != 1 {
			t.Errorf("expected returned lock count to match, got: %d", len(list.Locks))
		}
		if list.Locks[0].Id != lockId {
			t.Errorf("expected lockId to match, got: %s", list.Locks[0].Id)
		}
	*/
}
