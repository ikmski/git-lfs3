package adapter

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/ikmski/git-lfs3/entity"
	"github.com/ikmski/git-lfs3/usecase"
)

var testContentRepository usecase.ContentRepository

type TestS3 struct {
	s3iface.S3API
	getResult  s3.GetObjectOutput
	headResult s3.HeadObjectOutput
	err        error
}

type TestDownloader struct {
	s3manageriface.DownloaderAPI
	content *bytes.Buffer
	err     error
}

type TestUploader struct {
	s3manageriface.UploaderAPI
	output s3manager.UploadOutput
	err    error
}

const testS3BucketName = "test_bucket"

func (s TestS3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s.getResult, s.err
}

func (s TestS3) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	return &s.headResult, s.err
}

func (d TestDownloader) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {
	w.WriteAt(d.content.Bytes(), 0)
	return int64(d.content.Len()), d.err
}

func (u TestUploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	ioutil.ReadAll(input.Body)
	return &u.output, u.err
}

func TestContentStorePut(t *testing.T) {

	d := newTestData()
	testContentRepository = &contentRepository{
		s3: TestS3{
			headResult: s3.HeadObjectOutput{
				ContentLength: &d.contentSize,
			},
		},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     testS3BucketName,
	}

	m := &entity.MetaData{
		Oid:  d.contentOid,
		Size: d.contentSize,
	}

	b := bytes.NewBuffer([]byte(d.content))
	err := testContentRepository.Put(m, b)
	if err != nil {
		t.Fatalf("expected put to succeed, got: %s", err)
	}
}

func TestContentStorePutHashMismatch(t *testing.T) {

	d := newTestData()
	testContentRepository = &contentRepository{
		s3: TestS3{
			headResult: s3.HeadObjectOutput{
				ContentLength: &d.contentSize,
			},
		},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     testS3BucketName,
	}

	m := &entity.MetaData{
		Oid:  d.contentOid,
		Size: d.contentSize,
	}

	b := bytes.NewBuffer([]byte(d.bogusContent))

	err := testContentRepository.Put(m, b)
	if err == nil {
		t.Fatal("expected put with bogus content to fail")
	}
}

func TestContentStorePutSizeMismatch(t *testing.T) {

	d := newTestData()
	testContentRepository = &contentRepository{
		s3: TestS3{
			headResult: s3.HeadObjectOutput{
				ContentLength: &d.contentSize,
			},
		},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     testS3BucketName,
	}

	m := &entity.MetaData{
		Oid:  d.contentOid,
		Size: d.bogusContentSize,
	}

	b := bytes.NewBuffer([]byte(d.contentOid))

	err := testContentRepository.Put(m, b)
	if err == nil {
		t.Fatal("expected put with bogus size to fail")
	}
}

func TestContentStoreGet(t *testing.T) {

	d := newTestData()
	b := bytes.NewBuffer([]byte(d.content))

	testContentRepository = &contentRepository{
		s3: TestS3{},
		downloader: TestDownloader{
			content: b,
		},
		uploader: TestUploader{},
		bucket:   testS3BucketName,
	}

	m := &entity.MetaData{
		Oid:  d.contentOid,
		Size: d.contentSize,
	}

	fileName := "tmp_object"
	f, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("%s", err)
		return
	}
	defer f.Close()
	defer os.Remove(fileName)

	_, err = testContentRepository.Get(m, f, 0, 0)
	if err != nil {
		t.Fatalf("expected get to succeed, got: %s", err)
	}

	by, _ := ioutil.ReadFile(fileName)
	if string(by) != d.content {
		t.Fatalf("expected to read content, got: %s", string(by))
	}
}

func TestContenStoreNotExists(t *testing.T) {

	d := newTestData()

	testContentRepository = &contentRepository{
		s3: TestS3{
			err: errors.New("error"),
		},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     testS3BucketName,
	}

	m := &entity.MetaData{
		Oid:  d.contentOid,
		Size: d.contentSize,
	}

	if testContentRepository.Exists(m) {
		t.Fatalf("expected to get an error, but content existed")
	}
}

func TestContentStoreExists(t *testing.T) {

	d := newTestData()

	testContentRepository = &contentRepository{
		s3:         TestS3{},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     testS3BucketName,
	}

	m := &entity.MetaData{
		Oid:  d.contentOid,
		Size: d.contentSize,
	}

	if !testContentRepository.Exists(m) {
		t.Fatalf("expected content to exist")
	}
}
