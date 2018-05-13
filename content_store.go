package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

var (
	errHashMismatch = errors.New("Content hash does not match OID")
	errSizeMismatch = errors.New("Content size does not match")
)

type ContentStore struct {
	s3         s3iface.S3API
	downloader s3manageriface.DownloaderAPI
	uploader   s3manageriface.UploaderAPI
	bucket     string
}

func NewContentStore(credentials *credentials.Credentials, region string, bucket string) (*ContentStore, error) {

	c := new(ContentStore)

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials,
		Region:      aws.String(region),
	}))

	c.s3 = s3.New(sess)
	c.downloader = s3manager.NewDownloader(sess)
	c.uploader = s3manager.NewUploader(sess)
	c.bucket = bucket

	return c, nil
}

func (s *ContentStore) Get(meta *MetaObject, w io.WriterAt) (int64, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(transformKey(meta.Oid)),
	}

	return s.downloader.Download(w, input)
}

func (s *ContentStore) Put(meta *MetaObject, r io.Reader) error {

	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(transformKey(meta.Oid)),
		Body:   r,
	}

	_, err := s.uploader.Upload(uploadInput)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return err
	}

	hash := sha256.New()
	_, err = io.Copy(hash, r)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	shaStr := hex.EncodeToString(hash.Sum(nil))
	if shaStr != meta.Oid {
		return errHashMismatch
	}

	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(transformKey(meta.Oid)),
	}

	result, err := s.s3.HeadObject(headInput)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return err
	}

	if *result.ContentLength != meta.Size {
		return errSizeMismatch
	}

	return nil
}

func (s *ContentStore) Exists(meta *MetaObject) bool {

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(transformKey(meta.Oid)),
	}

	_, err := s.s3.HeadObject(input)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return false
	}

	return true
}

func transformKey(key string) string {

	if len(key) < 5 {
		return key
	}

	return filepath.Join(key[0:2], key[2:4], key[4:len(key)])
}
