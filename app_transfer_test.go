package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/ikmski/git-lfs3/entity"
)

func TestDownload(t *testing.T) {

	path := fmt.Sprintf("%s/%s/%s/objects/%s", lfsServer.URL, testUser1, testRepo, testContentOid)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}
	req.Header.Set("Accept", "application/vnd.git-lfs")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	by, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected response to contain content, got error: %s", err)
	}

	if string(by) != testContent {
		t.Fatalf("expected content to be `content`, got: %s", string(by))
	}
}

func TestDownloadWithRange(t *testing.T) {

	path := fmt.Sprintf("%s/%s/%s/objects/%s", lfsServer.URL, testUser1, testRepo, testContentOid)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}
	req.Header.Set("Accept", "application/vnd.git-lfs")

	fromByte := 5
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", fromByte))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}

	if res.StatusCode != 206 {
		t.Fatalf("expected status 206, got %d", res.StatusCode)
	}

	cr := res.Header.Get("Content-Range")
	if len(cr) > 0 {
		expected := fmt.Sprintf("bytes %d-%d/%d", fromByte, len(testContent)-1, len(testContent)-fromByte)
		if cr != expected {
			t.Fatalf("expected Content-Range header of %q, got %q", expected, cr)
		}
	} else {
		t.Fatalf("missing Content-Range header in response")
	}

	by, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected response to contain content, got error: %s", err)
	}

	if string(by) != testContent[fromByte:] {
		t.Fatalf("expected content to be `content`, got: %s", string(by))
	}

}

func TestUpload(t *testing.T) {

	path := fmt.Sprintf("%s/%s/%s/objects/%s", lfsServer.URL, testUser1, testRepo, testContentOid)
	req, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}
	req.Header.Set("Accept", "application/vnd.git-lfs")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(testContent)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("response error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	fileName := "tmp_object"
	f, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("%s", err)
		return
	}
	defer f.Close()
	defer os.Remove(fileName)

	m := &entity.MetaData{
		Oid:  testContentOid,
		Size: testContentSize,
	}
	_, err = testContentRepo.Get(m, f, 0, 0)
	if err != nil {
		t.Fatalf("error retreiving from content store: %s", err)
	}

	c, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("error reading content: %s", err)
	}

	if string(c) != testContent {
		t.Fatalf("expected content, got `%s`", string(c))
	}
}
