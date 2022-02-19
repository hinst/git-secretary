package git_secretary

import (
	"encoding/json"

	git_stories_api "github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
	bolt "go.etcd.io/bbolt"
)

type Storage struct {
	RepositoryLogEntriesBucketName      string
	RepositoryLogEntriesBucketNameBytes []byte
	FilePath                            string
	db                                  *bolt.DB
}

func (me *Storage) Create() *Storage {
	me.FilePath = common.ExecutableFileDirectory + "/storage.bolt"
	me.RepositoryLogEntriesBucketName = "RepositoryLogEntries"
	me.RepositoryLogEntriesBucketNameBytes = []byte(me.RepositoryLogEntriesBucketName)

	var dbOptions = *bolt.DefaultOptions
	dbOptions.Timeout = 1
	dbOptions.ReadOnly = false
	var db, e = bolt.Open(me.FilePath, FILE_PERMISSION_OWNER_READ_WRITE, &dbOptions)
	common.AssertError(common.CreateExceptionIf("Unable to open db file "+me.FilePath, e))
	me.db = db
	return me
}

func (me *Storage) ReadRepositoryLogEntry(commitHash string) (result *git_stories_api.RepositoryLogEntry) {
	var e = me.db.View(func(transaction *bolt.Tx) error {
		var bucket = transaction.Bucket(me.RepositoryLogEntriesBucketNameBytes)
		if bucket != nil {
			var cachedRowBytes = bucket.Get([]byte(commitHash))
			if cachedRowBytes != nil {
				var jsonError = json.Unmarshal(cachedRowBytes, result)
				if nil != jsonError {
					return jsonError
				}
			}
		}
		return nil
	})
	if e != nil {
		panic(common.CreateException("Unable to read repository log entry", e))
	}
	return
}
