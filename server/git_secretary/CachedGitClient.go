package git_secretary

import (
	git_stories_api "github.com/hinst/git-stories-api"
)

const BUCKET_NAME_LOG_ENTRY_ROWS = "LogEntryRows"
const CACHED_GIT_CLIENT_PAGE_SIZE = 1000

var BUCKET_NAME_LOG_ENTRY_ROWS_BYTES = []byte(BUCKET_NAME_LOG_ENTRY_ROWS)

type CachedGitClientReceiveProgressFunction func(total int, done int)

type CachedGitClient struct {
	storage         *Storage
	gitClient       *GitClient
	receiveProgress CachedGitClientReceiveProgressFunction
}

func (client *CachedGitClient) Create(storage *Storage, directory string) *CachedGitClient {
	client.storage = storage
	client.gitClient = CreateGitClient(directory)
	return client
}

func (client *CachedGitClient) SetProgressReceiver(function CachedGitClientReceiveProgressFunction) {
	client.receiveProgress = function
}

func (client *CachedGitClient) ReadDetailedLog(lengthLimit int) ([]*git_stories_api.RepositoryLogEntry, error) {
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
	storage         *Storage
	gitClient       *GitClient
	receiveProgress CachedGitClientReceiveProgressFunction
	allLogEntries   RepositoryLogEntryHeaders
	entries         []*git_stories_api.RepositoryLogEntry
	newEntries      []*git_stories_api.RepositoryLogEntry
}

func (reader *cachedGitClient_DetailedLogReader) Create(storage *Storage, gitClient *GitClient,
	reportProgress CachedGitClientReceiveProgressFunction) *cachedGitClient_DetailedLogReader {
	reader.storage = storage
	reader.gitClient = gitClient
	reader.receiveProgress = reportProgress
	reader.entries = nil
	reader.newEntries = nil
	return reader
}

func (reader *cachedGitClient_DetailedLogReader) Load(logEntries RepositoryLogEntryHeaders) error {
	reader.allLogEntries = logEntries
	var logEntryGroups = logEntries.GetPortions(CACHED_GIT_CLIENT_PAGE_SIZE)
	for groupIndex, logEntries := range logEntryGroups {
		if e := reader.loadEntries(groupIndex, logEntries); e != nil {
			return e
		}
		if len(reader.newEntries) > 0 {
			if e := reader.storage.WriteRepositoryLogEntries(reader.newEntries); e != nil {
				return e
			}
		}
		reader.newEntries = nil
	}
	return nil
}

func (reader *cachedGitClient_DetailedLogReader) loadEntries(groupIndex int, logHeaders RepositoryLogEntryHeaders) error {
	var cachedEntries, e = reader.loadCachedEntries(logHeaders)
	if e != nil {
		return e
	}
	for headerIndex, header := range logHeaders {
		var logEntry = cachedEntries[header.CommitHash]
		if logEntry == nil { // new row
			var newLogEntry, e = reader.gitClient.ReadDetailedLogEntryRow(header)
			if nil != e {
				return e
			}
			logEntry = &newLogEntry
			reader.newEntries = append(reader.newEntries, logEntry)
		}
		reader.entries = append(reader.entries, logEntry)
		if reader.receiveProgress != nil {
			var overallEntryIndex = (CACHED_GIT_CLIENT_PAGE_SIZE * groupIndex) + headerIndex
			reader.receiveProgress(len(reader.allLogEntries), overallEntryIndex)
		}
	}
	return nil
}

func (reader *cachedGitClient_DetailedLogReader) loadCachedEntries(logHeaders RepositoryLogEntryHeaders) (
	map[string]*git_stories_api.RepositoryLogEntry, error,
) {
	var commitHashes = make([]string, len(logHeaders))
	for i := range logHeaders {
		commitHashes[i] = logHeaders[i].CommitHash
	}
	var entryArray, e = reader.storage.ReadRepositoryLogEntries(commitHashes)
	if e != nil {
		return nil, e
	}
	var entryMap = make(map[string]*git_stories_api.RepositoryLogEntry)
	for _, entry := range entryArray {
		entryMap[entry.CommitHash] = entry
	}
	return entryMap, nil
}

func (reader *cachedGitClient_DetailedLogReader) GetRows() []*git_stories_api.RepositoryLogEntry {
	return reader.entries
}
