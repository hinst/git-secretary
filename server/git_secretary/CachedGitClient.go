package git_secretary

import (
	"encoding/json"

	git_stories_api "github.com/hinst/git-stories-api"
	"github.com/hinst/go-common"
	bolt "go.etcd.io/bbolt"
)

const BUCKET_NAME_LOG_ENTRY_ROWS = "LogEntryRows"
const CACHED_GIT_CLIENT_PAGE_SIZE = 1000

var BUCKET_NAME_LOG_ENTRY_ROWS_BYTES = []byte(BUCKET_NAME_LOG_ENTRY_ROWS)

type CachedGitClientReceiveProgressFunction func(total int, done int)

type CachedGitClient struct {
	storage         *bolt.DB
	gitClient       *GitClient
	receiveProgress CachedGitClientReceiveProgressFunction
}

func (client *CachedGitClient) Create(storage *bolt.DB, directory string) *CachedGitClient {
	client.storage = storage
	client.gitClient = CreateGitClient(directory)
	return client
}

func (client *CachedGitClient) SetProgressReceiver(function CachedGitClientReceiveProgressFunction) {
	client.receiveProgress = function
}

func (client *CachedGitClient) ReadDetailedLog(lengthLimit int) ([]git_stories_api.RepositoryLogEntry, error) {
	var logEntries, readError = client.gitClient.ReadLog(lengthLimit)
	if nil != readError {
		return nil, readError
	}
	var logReader cachedGitClient_DetailedLogReader
	logReader.Create(client.storage, client.gitClient, client.receiveProgress)
	var error = logReader.Load(logEntries)
	return logReader.GetRows(), error
}

type cachedGitClient_DetailedLogReader struct {
	storage         *bolt.DB
	gitClient       *GitClient
	receiveProgress CachedGitClientReceiveProgressFunction
	allLogEntries   LogEntryRows
	rows            []git_stories_api.RepositoryLogEntry
	newRows         []git_stories_api.RepositoryLogEntry
}

func (reader *cachedGitClient_DetailedLogReader) Create(storage *bolt.DB, gitClient *GitClient,
	reportProgress CachedGitClientReceiveProgressFunction) *cachedGitClient_DetailedLogReader {
	reader.storage = storage
	reader.gitClient = gitClient
	reader.receiveProgress = reportProgress
	reader.rows = nil
	reader.newRows = nil
	return reader
}

func (reader *cachedGitClient_DetailedLogReader) Load(logEntries LogEntryRows) error {
	reader.allLogEntries = logEntries
	var logEntryGroups = logEntries.GetPortions(CACHED_GIT_CLIENT_PAGE_SIZE)
	for groupIndex, logEntries := range logEntryGroups {
		var transactionError = reader.storage.View(func(transaction *bolt.Tx) error {
			return reader.loadRows(groupIndex, logEntries, transaction)
		})
		if nil != transactionError {
			return transactionError
		}
		if len(reader.newRows) > 0 {
			transactionError = reader.storage.Update(reader.storeNewRows)
		}
		if nil != transactionError {
			return transactionError
		}
		reader.newRows = nil
	}
	return nil
}

func (reader *cachedGitClient_DetailedLogReader) loadRows(groupIndex int, logEntries LogEntryRows, transaction *bolt.Tx) error {
	var bucket = transaction.Bucket(BUCKET_NAME_LOG_ENTRY_ROWS_BYTES)
	for entryIndex, entry := range logEntries {
		var row git_stories_api.RepositoryLogEntry
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
		if reader.receiveProgress != nil {
			var overallEntryIndex = (CACHED_GIT_CLIENT_PAGE_SIZE * groupIndex) + entryIndex
			reader.receiveProgress(len(reader.allLogEntries), overallEntryIndex)
		}
	}
	return nil
}

func (reader *cachedGitClient_DetailedLogReader) storeNewRows(transaction *bolt.Tx) error {
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

func (reader *cachedGitClient_DetailedLogReader) GetRows() []git_stories_api.RepositoryLogEntry {
	return reader.rows
}
