package adapter

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/ikmski/git-lfs3/entity"
	"github.com/ikmski/git-lfs3/usecase"
)

var (
	errHashMismatch = errors.New("Content hash does not match OID")
	errSizeMismatch = errors.New("Content size does not match")
)

type contentRepository struct {
	s3         s3iface.S3API
	downloader s3manageriface.DownloaderAPI
	uploader   s3manageriface.UploaderAPI
	bucket     string
}

func NewContentRepository(bucket string) (usecase.ContentRepository, error) {

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	r := &contentRepository{
		s3:         s3.New(sess),
		downloader: s3manager.NewDownloader(sess),
		uploader:   s3manager.NewUploader(sess),
		bucket:     bucket,
	}

	return r, nil
}

func (r *contentRepository) Get(meta *entity.MetaData, w io.WriterAt, from int64, to int64) (int64, error) {

	rangeHeader := ""
	if from > 0 && to > from {
		rangeHeader = fmt.Sprintf("bytes=%s-%s", strconv.FormatInt(from, 10), strconv.FormatInt(to, 10))
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(transformKey(meta.Oid)),
		Range:  &rangeHeader,
	}

	return r.downloader.Download(w, input)
}

func (r *contentRepository) Put(meta *entity.MetaData, reader io.Reader) error {

	hash := sha256.New()
	tee := io.TeeReader(reader, hash)

	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(transformKey(meta.Oid)),
		Body:   tee,
	}

	_, err := r.uploader.Upload(uploadInput)
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

	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(transformKey(meta.Oid)),
	}

	result, err := r.s3.HeadObject(headInput)
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

	shaStr := hex.EncodeToString(hash.Sum(nil))
	if shaStr != meta.Oid {
		return errHashMismatch
	}

	return nil
}

func (r *contentRepository) Exists(meta *entity.MetaData) bool {

	input := &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(transformKey(meta.Oid)),
	}

	_, err := r.s3.HeadObject(input)
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
