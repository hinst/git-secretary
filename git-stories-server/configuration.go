package main

import "runtime"

type Configuration struct {
	Plugin     string
	PortNumber int
}

func CheckWindows() bool {
	return runtime.GOOS == "windows"
}
