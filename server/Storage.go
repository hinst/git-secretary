package main

import (
	"github.com/hinst/go-common"
	bolt "go.etcd.io/bbolt"
)

type Storage struct {
}

func (me *Storage) open() {
	var dbOptions = *bolt.DefaultOptions
	dbOptions.Timeout = 1
	dbOptions.ReadOnly = false
	var storage, e = bolt.Open("./storage.bolt", FILE_PERMISSION_OWNER_READ_WRITE, &dbOptions)
	common.AssertError(e)
	common.Use(storage)
}
