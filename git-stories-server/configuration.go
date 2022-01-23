package main

import "runtime"

type Configuration struct {
	Plugin     string
	PortNumber int
}

const WINDOWS_EXECUTABLE_FILE_EXTENSION = ".exe"

func CheckWindows() bool {
	return runtime.GOOS == "windows"
}
