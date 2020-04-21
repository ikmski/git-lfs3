package adapter

import (
	"testing"
)

func TestGetMeta(t *testing.T) {

	d := newTestData()
	setupRepository(d)
	defer teardownRepository(d)

	meta, err := d.metaDataRepository.Get(d.contentOid)
	if err != nil {
		t.Fatalf("Error retreiving meta: %s", err)
	}

	if meta.Oid != d.contentOid {
		t.Errorf("expected to get content oid, got: %s", meta.Oid)
	}

	if meta.Size != d.contentSize {
		t.Errorf("expected to get content size, got: %d", meta.Size)
	}
}

func TestPutMeta(t *testing.T) {

	d := newTestData()
	setupRepository(d)
	defer teardownRepository(d)

	meta, err := d.metaDataRepository.Put(d.nonExistContentOid, d.nonExitContentSize)
	if err != nil {
		t.Errorf("expected put to succeed, got : %s", err)
	}

	/*
		if meta.Existing {
			t.Errorf("expected meta to not have existed")
		}
	*/

	meta, err = d.metaDataRepository.Get(d.nonExistContentOid)
	if err != nil {
		t.Errorf("expected to be able to retreive new put, got : %s", err)
	}

	if meta.Oid != d.nonExistContentOid {
		t.Errorf("expected oids to match, got: %s", meta.Oid)
	}

	if meta.Size != d.nonExitContentSize {
		t.Errorf("expected sizes to match, got: %d", meta.Size)
	}

	meta, err = d.metaDataRepository.Put(d.nonExistContentOid, d.nonExitContentSize)
	if err != nil {
		t.Errorf("expected put to succeed, got : %s", err)
	}

	/*
		if !meta.Existing {
			t.Errorf("expected meta to now exist")
		}
	*/
}
