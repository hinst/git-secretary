package main

import "unsafe"

var SizeOfInt = unsafe.Sizeof(int(0))
var BitSizeOfInt = SizeOfInt * 8

var SizeOfUint = unsafe.Sizeof(uint(0))
var BitSizeOfUInt = SizeOfUint * 8
