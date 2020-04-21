package adapter

import (
	"testing"
)

func TestGetMeta(t *testing.T) {

	setupRepository()
	defer teardownRepository()

	meta, err := testMetaDataRepository.Get(testContentOid)
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

	setupRepository()
	defer teardownRepository()

	meta, err := testMetaDataRepository.Put(testNonExistingOid, 42)
	if err != nil {
		t.Errorf("expected put to succeed, got : %s", err)
	}

	/*
		if meta.Existing {
			t.Errorf("expected meta to not have existed")
		}
	*/

	meta, err = testMetaDataRepository.Get(testNonExistingOid)
	if err != nil {
		t.Errorf("expected to be able to retreive new put, got : %s", err)
	}

	if meta.Oid != testNonExistingOid {
		t.Errorf("expected oids to match, got: %s", meta.Oid)
	}

	if meta.Size != 42 {
		t.Errorf("expected sizes to match, got: %d", meta.Size)
	}

	meta, err = testMetaDataRepository.Put(testNonExistingOid, 42)
	if err != nil {
		t.Errorf("expected put to succeed, got : %s", err)
	}

	/*
		if !meta.Existing {
			t.Errorf("expected meta to now exist")
		}
	*/
}
