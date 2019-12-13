package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestBatchDownload(t *testing.T) {

	path := fmt.Sprintf("%s/%s/%s/objects/batch", lfsServer.URL, testUser1, testRepo)

	var objs []*ObjectRequest
	obj := &ObjectRequest{
		Oid:  testContentOid,
		Size: testContentSize,
	}
	objs = append(objs, obj)

	requestData := &BatchRequest{
		Operation: "download",
		Objects:   objs,
	}

	requestBody, _ := json.Marshal(requestData)

	req, err := http.NewRequest("POST", path, bytes.NewReader(requestBody))
	if err != nil {
		t.Fatalf("request error: %s", err)
	}
	req.Header.Set("Accept", metaMediaType)
	req.Header.Set("Content-Type", metaMediaType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected response to contain content, got error: %s", err)
	}

	var responseData BatchResponse
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if responseData.Transfer != "basic" {
		t.Errorf("got %v\nwant %v", responseData.Transfer, "basic")
	}

	if responseData.Objects[0].Oid != testContentOid {
		t.Errorf("got %v\nwant %v", responseData.Objects[0].Oid, testContentOid)
	}
	if responseData.Objects[0].Size != testContentSize {
		t.Errorf("got %v\nwant %v", responseData.Objects[0].Size, testContentSize)
	}
	_, ok := responseData.Objects[0].Actions["download"]
	if !ok {
		t.Errorf("got %v\nwant %v", ok, true)
	}

}

func TestBatchUpload(t *testing.T) {

	path := fmt.Sprintf("%s/%s/%s/objects/batch", lfsServer.URL, testUser1, testRepo)

	var objs []*ObjectRequest
	obj := &ObjectRequest{
		Oid:  testContentOid,
		Size: testContentSize,
	}
	objs = append(objs, obj)

	requestData := &BatchRequest{
		Operation: "upload",
		Objects:   objs,
	}

	requestBody, _ := json.Marshal(requestData)

	req, err := http.NewRequest("POST", path, bytes.NewReader(requestBody))
	if err != nil {
		t.Fatalf("request error: %s", err)
	}
	req.Header.Set("Accept", metaMediaType)
	req.Header.Set("Content-Type", metaMediaType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("expected response to contain content, got error: %s", err)
	}

	var responseData BatchResponse
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if responseData.Transfer != "basic" {
		t.Errorf("got %v\nwant %v", responseData.Transfer, "basic")
	}

	if responseData.Objects[0].Oid != testContentOid {
		t.Errorf("got %v\nwant %v", responseData.Objects[0].Oid, testContentOid)
	}
	if responseData.Objects[0].Size != testContentSize {
		t.Errorf("got %v\nwant %v", responseData.Objects[0].Size, testContentSize)
	}
	_, ok := responseData.Objects[0].Actions["upload"]
	if !ok {
		t.Errorf("got %v\nwant %v", ok, true)
	}

}
