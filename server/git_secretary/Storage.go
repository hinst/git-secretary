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

	FilePath string
	db       *bolt.DB
}

func (me *Storage) Create() *Storage {
	me.RepositoryLogEntriesBucketName = "RepositoryLogEntries"
	me.RepositoryLogEntriesBucketNameBytes = []byte(me.RepositoryLogEntriesBucketName)
	me.FilePath = common.ExecutableFileDirectory + "/storage.bolt"

	var dbOptions = *bolt.DefaultOptions
	dbOptions.Timeout = 1
	dbOptions.ReadOnly = false
	var db, e = bolt.Open(me.FilePath, FILE_PERMISSION_OWNER_READ_WRITE, &dbOptions)
	common.AssertError(common.CreateExceptionIf("Unable to open db file "+me.FilePath, e))
	me.db = db
	return me
}

func (me *Storage) ReadRepositoryLogEntry(commitHash string) (
	result *git_stories_api.RepositoryLogEntry, e error,
) {
	e = me.db.View(func(transaction *bolt.Tx) error {
		var bucket = transaction.Bucket(me.RepositoryLogEntriesBucketNameBytes)
		if bucket != nil {
			var cachedRowBytes = bucket.Get([]byte(commitHash))
			if cachedRowBytes != nil {
				result = &git_stories_api.RepositoryLogEntry{}
				var jsonError = json.Unmarshal(cachedRowBytes, result)
				if nil != jsonError {
					return jsonError
				}
			}
		}
		return nil
	})
	return
}

func (me *Storage) ReadRepositoryLogEntries(commitHashes []string) (
	result []*git_stories_api.RepositoryLogEntry, e error,
) {
	e = me.db.View(func(transaction *bolt.Tx) error {
		var bucket = transaction.Bucket(me.RepositoryLogEntriesBucketNameBytes)
		if bucket != nil {
			for _, commitHash := range commitHashes {
				var cachedRowBytes = bucket.Get([]byte(commitHash))
				if cachedRowBytes != nil {
					var entry = &git_stories_api.RepositoryLogEntry{}
					var jsonError = json.Unmarshal(cachedRowBytes, entry)
					if nil != jsonError {
						return jsonError
					}
					result = append(result, entry)
				}
			}
		}
		return nil
	})
	return
}

func (me *Storage) WriteRepositoryLogEntries(entries []*git_stories_api.RepositoryLogEntry) error {
	return me.db.Update(func(transaction *bolt.Tx) error {
		var bucket, dbError = transaction.CreateBucketIfNotExists(me.RepositoryLogEntriesBucketNameBytes)
		if nil != dbError {
			return dbError
		}
		for _, row := range entries {
			var rowBytes, jsonError = json.Marshal(row)
			if nil != jsonError {
				return jsonError
			}
			dbError = bucket.Put([]byte(row.CommitHash), rowBytes)
			if nil != dbError {
				return dbError
			}
		}
		return nil
	})
}
