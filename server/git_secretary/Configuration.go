package git_secretary

import "runtime"

type Configuration struct {
	Plugin          string
	WebPath         string
	PortNumber      int
	AutoOpenEnabled bool
}

func (configuration *Configuration) SetDefault() *Configuration {
	configuration.Plugin = "plugins/story-girls-standard"
	configuration.WebPath = "/git-stories"
	configuration.PortNumber = 3003
	configuration.AutoOpenEnabled = true
	return configuration
}

const WINDOWS_EXECUTABLE_FILE_EXTENSION = ".exe"

func CheckWindows() bool {
	return runtime.GOOS == "windows"
}
