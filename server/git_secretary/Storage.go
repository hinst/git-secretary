package git_secretary

import (
	"github.com/hinst/go-common"
	bolt "go.etcd.io/bbolt"
)

type Storage struct {
}

func (me *Storage) Create() *Storage {
	var dbOptions = *bolt.DefaultOptions
	dbOptions.Timeout = 1
	dbOptions.ReadOnly = false
	var storage, e = bolt.Open(me.GetFilePath(), FILE_PERMISSION_OWNER_READ_WRITE, &dbOptions)
	common.AssertError(e)
	common.Use(storage)
	return me
}

func (me *Storage) GetFilePath() string {
	return common.GetExecutableDir() + "/storage.bolt"
}
