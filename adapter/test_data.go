package adapter

import (
	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/usecase"
)

type TestData struct {
	databaseFile string

	userName1 string
	userPass1 string
	userName2 string
	userPass2 string

	repoName string

	content            string
	contentOid         string
	contentSize        int64
	bogusContent       string
	bogusContentSize   int64
	nonExistContentOid string
	nonExitContentSize int64

	lockID         string
	nonExistLockID string
	lockPath       string

	database           *bolt.DB
	metaDataRepository usecase.MetaDataRepository
	lockRepository     usecase.LockRepository
	userRepository     usecase.UserRepository
}

func newTestData() *TestData {
	return &TestData{

		databaseFile: "test-database.db",
		userName1:    "bilbo1",
		userPass1:    "baggins1",
		userName2:    "bilbo2",
		userPass2:    "baggins2",

		repoName: "repo",

		content:            "this is my content",
		contentOid:         "f97e1b2936a56511b3b6efc99011758e4700d60fb1674d31445d1ee40b663f24",
		contentSize:        18,
		bogusContent:       "this is bogus content",
		bogusContentSize:   21,
		nonExistContentOid: "aec070645fe53ee3b3763059376134f058cc337247c978add178b6ccdfb0019f",
		nonExitContentSize: 42,

		lockID:         "3cfec93346f7ff337c60f2da50cd86740715e2f6",
		nonExistLockID: "f310c1555a2485e2e5229ea015a94c9d590763d3",
		lockPath:       "this/is/lock/path",
	}
}
