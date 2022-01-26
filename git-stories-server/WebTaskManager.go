package main

import (
	"sync"

	git_stories_api "github.com/hinst/git-stories-api"
)

type WebTask struct {
	Total        int                          `json:"total"`
	Done         int                          `json:"done"`
	Error        string                       `json:"error"`
	StoryEntries []git_stories_api.StoryEntry `json:"storyEntries"`
}

func (task *WebTask) IsDone() bool {
	return task.StoryEntries != nil || len(task.Error) > 0
}

type WebTaskManager struct {
	counter uint
	tasks   map[uint]*WebTask
	locker  sync.Mutex
}

func (manager *WebTaskManager) Create() *WebTaskManager {
	manager.counter = 0
	manager.tasks = make(map[uint]*WebTask)
	return manager
}

func (manager *WebTaskManager) Add(task *WebTask) uint {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	manager.counter += 1
	manager.tasks[manager.counter] = task
	return manager.counter
}

func (manager *WebTaskManager) Update(id uint, update func(task *WebTask)) {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	update(manager.tasks[id])
}

func (manager *WebTaskManager) Get(id uint) *WebTask {
	manager.locker.Lock()
	defer manager.locker.Unlock()

	var task = manager.tasks[id]
	if task != nil && task.IsDone() {
		delete(manager.tasks, id)
	}
	return task
}
