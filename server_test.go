package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var (
	lfsServer        *httptest.Server
	testMetaStore    *MetaStore
	testContentStore *ContentStore
)

const (
	testUser1         = "bilbo1"
	testPass1         = "baggins1"
	testUser2         = "bilbo2"
	testPass2         = "baggins2"
	testRepo          = "repo"
	content           = "this is my content"
	contentSize       = int64(len(content))
	contentOid        = "f97e1b2936a56511b3b6efc99011758e4700d60fb1674d31445d1ee40b663f24"
	nonExistingOid    = "aec070645fe53ee3b3763059376134f058cc337247c978add178b6ccdfb0019f"
	lockId            = "3cfec93346f7ff337c60f2da50cd86740715e2f6"
	nonExistingLockId = "f310c1555a2485e2e5229ea015a94c9d590763d3"
	lockPath          = "this/is/lock/path"
)

func TestMain(m *testing.M) {

	os.Remove("lfs-test.db")

	var err error
	testMetaStore, err = NewMetaStore("lfs-test.db")
	if err != nil {
		fmt.Printf("Error creating meta store: %s", err)
		os.Exit(1)
	}

	cleanup, addr := prepareS3Container()
	defer cleanup()

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("dummy", "dummy", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.ApNortheast1RegionID),
		Endpoint:         aws.String(addr),
	}))

	testContentStore, err = NewContentStore(sess, "lfs-test-bucket")
	if err != nil {
		fmt.Printf("Error creating content store: %s", err)
		os.Exit(1)
	}

	err = createBucket("lfs-test-bucket")
	if err != nil {
		fmt.Printf("Error creating bucket: %s", err)
		os.Exit(1)
	}

	if err := seedContentStore(); err != nil {
		fmt.Printf("Error seeding content store: %s", err)
		os.Exit(1)
	}

	if err := seedMetaStore(); err != nil {
		fmt.Printf("Error seeding meta store: %s", err)
		os.Exit(1)
	}

	app := newApp(testMetaStore, testContentStore)
	lfsServer = httptest.NewServer(app)

	ret := m.Run()

	lfsServer.Close()
	testMetaStore.Close()

	os.Remove("lfs-test.db")

	os.Exit(ret)
}

func seedMetaStore() error {

	if err := testMetaStore.AddUser(testUser1, testPass1); err != nil {
		return err
	}
	if err := testMetaStore.AddUser(testUser2, testPass2); err != nil {
		return err
	}

	o := &Object{
		Oid:  contentOid,
		Size: contentSize,
	}

	_, err := testMetaStore.Put(o)
	if err != nil {
		return err
	}

	lock := NewTestLock(lockId, lockPath, testUser1)
	if err := testMetaStore.AddLocks(testRepo, lock); err != nil {
		return err
	}

	return nil
}

func createBucket(bucket string) error {

	createInput := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}
	_, err := testContentStore.s3.CreateBucket(createInput)
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

	return nil
}

func seedContentStore() error {

	meta := &MetaObject{
		Oid:  contentOid,
		Size: contentSize,
	}

	buf := bytes.NewBuffer([]byte(content))

	err := testContentStore.Put(meta, buf)
	if err != nil {
		return err
	}

	return nil
}

func prepareS3Container() (func(), string) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		fmt.Printf("couldn't not connect docker host: %s", err.Error())
		os.Exit(1)
	}

	resource, err := pool.Run("atlassianlabs/localstack", "latest", []string{})
	if err != nil {
		fmt.Printf("couldn't start S3 container: %s", err.Error())
		os.Exit(1)
	}

	addr := fmt.Sprintf("http://localhost:%s", resource.GetPort("4572/tcp"))

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			fmt.Printf("couldn't cleanup S3 container: %s", err.Error())
			os.Exit(1)
		}
	}

	err = pool.Retry(func() error {
		resp, err := http.Get(addr)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("didn't return status code 200: %s", resp.Status)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("couldn't prepare S3 container: %s", err.Error())
		os.Exit(1)
	}

	return cleanup, addr
}
