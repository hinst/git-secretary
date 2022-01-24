package main

import "runtime"

type Configuration struct {
	Plugin     string
	WebPath    string
	PortNumber int
}

func (configuration *Configuration) SetDefault() *Configuration {
	configuration.Plugin = "plugins/story-girls-standard"
	configuration.WebPath = "/git-stories"
	configuration.PortNumber = 3003
	return configuration
}

const WINDOWS_EXECUTABLE_FILE_EXTENSION = ".exe"

func CheckWindows() bool {
	return runtime.GOOS == "windows"
}
