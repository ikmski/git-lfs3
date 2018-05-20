package main

import (
	"bytes"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type MockedS3Data struct {
	data *bytes.Buffer
}

func NewMockedS3Data() *MockedS3Data {
	return &MockedS3Data{data: new(bytes.Buffer)}
}

func (d *MockedS3Data) Read(p []byte) (int, error) {
	return d.data.Read(p)
}

func (d *MockedS3Data) Write(p []byte) (int, error) {
	d.data.Reset()
	return d.data.Write(p)
}

func (d *MockedS3Data) Close() error {
	return nil
}

func (d *MockedS3Data) Bytes() []byte {
	return d.data.Bytes()
}

func (d *MockedS3Data) Len() int64 {
	return int64(d.data.Len())
}

var mockedDataStore map[string]*MockedS3Data

type MockedS3 struct {
	s3iface.S3API
	getResult  s3.GetObjectOutput
	headResult s3.HeadObjectOutput
	err        error
}

type MockedDownloader struct {
	s3manageriface.DownloaderAPI
	content *bytes.Buffer
	err     error
}

type MockedUploader struct {
	s3manageriface.UploaderAPI
	output s3manager.UploadOutput
	err    error
}

func (ms MockedS3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {

	d, ok := mockedDataStore[*input.Key]
	if ok {
		result := s3.GetObjectOutput{
			Body: d,
		}

		return &result, nil
	}

	return &s3.GetObjectOutput{}, errors.New("")
}

func (ms MockedS3) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {

	d, ok := mockedDataStore[*input.Key]
	if ok {
		len := d.Len()
		result := s3.HeadObjectOutput{
			ContentLength: &len,
		}

		return &result, nil
	}

	return &s3.HeadObjectOutput{}, errors.New("")
}

func (md MockedDownloader) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {

	d, ok := mockedDataStore[*input.Key]
	if ok {
		w.WriteAt(d.Bytes(), 0)
		return d.Len(), nil
	}

	return int64(0), errors.New("")
}

func (mu MockedUploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {

	d, ok := mockedDataStore[*input.Key]
	if ok {
		_, err := io.Copy(d, input.Body)
		if err != nil {
			return &s3manager.UploadOutput{}, err
		}
		return &s3manager.UploadOutput{}, nil
	}

	d = NewMockedS3Data()
	_, err := io.Copy(d, input.Body)
	if err != nil {
		return &s3manager.UploadOutput{}, err
	}
	mockedDataStore[*input.Key] = d
	return &s3manager.UploadOutput{}, nil
}

func NewMockedContentStore(bucket string) (*ContentStore, error) {

	mockedDataStore = make(map[string]*MockedS3Data)
	contentStore := &ContentStore{
		s3:         MockedS3{},
		downloader: MockedDownloader{},
		uploader:   MockedUploader{},
		bucket:     bucket,
	}

	return contentStore, nil
}
