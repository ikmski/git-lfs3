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

	b := bytes.NewBuffer([]byte("test content"))
	var contentSize int64 = 12

	testContentRepository = &contentRepository{
		s3: TestS3{
			headResult: s3.HeadObjectOutput{
				ContentLength: &contentSize,
			},
		},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     "test_bucket",
	}

	m := &entity.MetaData{
		Oid:  "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72",
		Size: contentSize,
	}

	err := testContentRepository.Put(m, b)
	if err != nil {
		t.Fatalf("expected put to succeed, got: %s", err)
	}
}

func TestContentStorePutHashMismatch(t *testing.T) {

	var contentSize int64 = 12

	testContentRepository = &contentRepository{
		s3: TestS3{
			headResult: s3.HeadObjectOutput{
				ContentLength: &contentSize,
			},
		},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     "test_bucket",
	}

	m := &entity.MetaData{
		Oid:  "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72",
		Size: contentSize,
	}

	b := bytes.NewBuffer([]byte("bogus content"))

	err := testContentRepository.Put(m, b)
	if err == nil {
		t.Fatal("expected put with bogus content to fail")
	}
}

func TestContentStorePutSizeMismatch(t *testing.T) {

	var contentSize int64 = 12

	testContentRepository = &contentRepository{
		s3: TestS3{
			headResult: s3.HeadObjectOutput{
				ContentLength: &contentSize,
			},
		},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     "test_bucket",
	}

	m := &entity.MetaData{
		Oid:  "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72",
		Size: 14,
	}

	b := bytes.NewBuffer([]byte("test content"))

	err := testContentRepository.Put(m, b)
	if err == nil {
		t.Fatal("expected put with bogus size to fail")
	}
}

func TestContentStoreGet(t *testing.T) {

	b := bytes.NewBuffer([]byte("test content"))

	testContentRepository = &contentRepository{
		s3: TestS3{},
		downloader: TestDownloader{
			content: b,
		},
		uploader: TestUploader{},
		bucket:   "test_bucket",
	}

	m := &entity.MetaData{
		Oid:  "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72",
		Size: 12,
	}

	fileName := "tmp_object"
	f, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("%s", err)
		return
	}
	defer f.Close()
	defer os.Remove(fileName)

	_, err = testContentRepository.Get(m, f, "")
	if err != nil {
		t.Fatalf("expected get to succeed, got: %s", err)
	}

	by, _ := ioutil.ReadAll(f)
	if string(by) != "test content" {
		t.Fatalf("expected to read content, got: %s", string(by))
	}
}

func TestContenStoreNotExists(t *testing.T) {

	testContentRepository = &contentRepository{
		s3: TestS3{
			err: errors.New("error"),
		},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     "test_bucket",
	}

	m := &entity.MetaData{
		Oid:  "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72",
		Size: 12,
	}

	if testContentRepository.Exists(m) {
		t.Fatalf("expected to get an error, but content existed")
	}
}

func TestContentStoreExists(t *testing.T) {

	testContentRepository = &contentRepository{
		s3:         TestS3{},
		downloader: TestDownloader{},
		uploader:   TestUploader{},
		bucket:     "test_bucket",
	}

	m := &entity.MetaData{
		Oid:  "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72",
		Size: 12,
	}

	if !testContentRepository.Exists(m) {
		t.Fatalf("expected content to exist")
	}
}
