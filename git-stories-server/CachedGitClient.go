package main

import (
	"encoding/json"

	git_stories_api "github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
	bolt "go.etcd.io/bbolt"
)

const BUCKET_NAME_LOG_ENTRY_ROWS = "LogEntryRows"

var BUCKET_NAME_LOG_ENTRY_ROWS_BYTES = []byte(BUCKET_NAME_LOG_ENTRY_ROWS)

type CachedGitClientReportProgressFunction func(done int, total int)

type CachedGitClient struct {
	storage        *bolt.DB
	gitClient      *GitClient
	reportProgress CachedGitClientReportProgressFunction
}

func (client *CachedGitClient) Create(storage *bolt.DB, directory string) *CachedGitClient {
	client.storage = storage
	client.gitClient = CreateGitClient(directory)
	return client
}

func (client *CachedGitClient) SetReportProgress(function CachedGitClientReportProgressFunction) {
	client.reportProgress = function
}

func (client *CachedGitClient) ReadDetailedLog(lengthLimit int) ([]git_stories_api.DetailedLogEntryRow, error) {
	var logEntries, readError = client.gitClient.ReadLog(lengthLimit)
	if nil != readError {
		return nil, readError
	}
	var logReader cachedGitClient_DetailedLogReader
	logReader.Create(client.storage, client.gitClient, client.reportProgress)
	var error = logReader.Load(logEntries)
	return logReader.GetRows(), error
}

type cachedGitClient_DetailedLogReader struct {
	storage        *bolt.DB
	gitClient      *GitClient
	reportProgress CachedGitClientReportProgressFunction
	logEntries     []LogEntryRow
	rows           []git_stories_api.DetailedLogEntryRow
	newRows        []git_stories_api.DetailedLogEntryRow
}

func (reader *cachedGitClient_DetailedLogReader) Create(storage *bolt.DB, gitClient *GitClient,
	reportProgress CachedGitClientReportProgressFunction) *cachedGitClient_DetailedLogReader {
	reader.storage = storage
	reader.gitClient = gitClient
	reader.reportProgress = reportProgress
	return reader
}

func (reader *cachedGitClient_DetailedLogReader) Load(logEntries []LogEntryRow) error {
	reader.logEntries = logEntries
	var transactionError = reader.storage.View(reader.loadRows)
	if nil != transactionError {
		return transactionError
	}
	if len(reader.newRows) > 0 {
		transactionError = reader.storage.Update(reader.storeCachedRows)
	}
	return transactionError
}

func (reader *cachedGitClient_DetailedLogReader) loadRows(transaction *bolt.Tx) error {
	var bucket = transaction.Bucket(BUCKET_NAME_LOG_ENTRY_ROWS_BYTES)
	for entryIndex, entry := range reader.logEntries {
		var row git_stories_api.DetailedLogEntryRow
		var cachedRowBytes []byte
		if bucket != nil {
			cachedRowBytes = bucket.Get([]byte(entry.CommitHash))
		}
		if cachedRowBytes == nil { // new row
			var e error
			row, e = reader.gitClient.ReadDetailedLogEntryRow(entry)
			if nil != e {
				return e
			}
			reader.newRows = append(reader.newRows, row)
		} else { // cached row
			var jsonError = json.Unmarshal(cachedRowBytes, &row)
			if nil != jsonError {
				return nil
			}
		}
		reader.rows = append(reader.rows, row)
		if reader.reportProgress != nil {
			reader.reportProgress(entryIndex, len(reader.logEntries))
		}
	}
	return nil
}

func (reader *cachedGitClient_DetailedLogReader) storeCachedRows(transaction *bolt.Tx) error {
	var bucket, bucketError = transaction.CreateBucketIfNotExists(BUCKET_NAME_LOG_ENTRY_ROWS_BYTES)
	if nil != bucketError {
		return common.CreateException("Unable to obtain bucket", bucketError)
	}
	for _, row := range reader.newRows {
		var rowBytes, jsonError = json.Marshal(row)
		if nil != jsonError {
			return jsonError
		}
		bucketError = bucket.Put([]byte(row.CommitHash), rowBytes)
		if nil != bucketError {
			return common.CreateException("Unable to write bucket", bucketError)
		}
	}
	return nil
}

func (reader *cachedGitClient_DetailedLogReader) GetRows() []git_stories_api.DetailedLogEntryRow {
	return reader.rows
}
