package main

import (
	"encoding/json"

	git_stories_api "github.com/hinst/git-stories-api"
	bolt "go.etcd.io/bbolt"
)

const BUCKET_NAME_LOG_ENTRY_ROWS = "LogEntryRows"

var BUCKET_NAME_LOG_ENTRY_ROWS_BYTES = []byte(BUCKET_NAME_LOG_ENTRY_ROWS)

type CachedGitClient struct {
	storage   *bolt.DB
	gitClient *GitClient
}

func (client *CachedGitClient) Create(directory string) *CachedGitClient {
	client.gitClient = CreateGitClient(directory)
	return client
}

func (client *CachedGitClient) ReadDetailedLog(lengthLimit int) ([]git_stories_api.DetailedLogEntryRow, error) {
	var logEntries, readError = client.gitClient.ReadLog(lengthLimit)
	if nil != readError {
		return nil, readError
	}
	var rows []git_stories_api.DetailedLogEntryRow
	var newRows []git_stories_api.DetailedLogEntryRow
	var transactionError = client.storage.View(func(transaction *bolt.Tx) error {
		var bucket, bucketError = transaction.CreateBucketIfNotExists(BUCKET_NAME_LOG_ENTRY_ROWS_BYTES)
		if nil != bucketError {
			return bucketError
		}
		for _, entry := range logEntries {
			var row git_stories_api.DetailedLogEntryRow
			var cachedRowBytes = bucket.Get([]byte(entry.CommitHash))
			if cachedRowBytes == nil { // new row
				var e error
				row, e = client.gitClient.ReadDetailedLogEntryRow(entry)
				if nil != e {
					return e
				}
				newRows = append(newRows, row)
			} else { // cached row
				var jsonError = json.Unmarshal(cachedRowBytes, &row)
				if nil != jsonError {
					return nil
				}
			}
			rows = append(rows, row)
		}
		return nil
	})
	if len(newRows) > 0 { // store newly loaded
		client.storage.Update(func(transaction *bolt.Tx) error {
			var bucket, bucketError = transaction.CreateBucketIfNotExists(BUCKET_NAME_LOG_ENTRY_ROWS_BYTES)
			if nil != bucketError {
				return bucketError
			}
			for _, row := range newRows {
				var rowBytes, jsonError = json.Marshal(row)
				if nil != jsonError {
					return jsonError
				}
				bucketError = bucket.Put([]byte(row.CommitHash), rowBytes)
				if nil != bucketError {
					return bucketError
				}
			}
			return nil
		})
	}
	if nil != transactionError {
		return nil, transactionError
	}
	return rows, nil
}

func (client *CachedGitClient) ReadDetailedLogRow(entry LogEntryRow, cachedRowBytes []byte) (
	row git_stories_api.DetailedLogEntryRow) {

}
