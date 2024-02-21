package utils

import "unsafe"

var SizeOfFloat32 int = int(unsafe.Sizeof(float32(0)))
var SizeOfFloat64 int = int(unsafe.Sizeof(float64(0)))

var SizeOfInt int = int(unsafe.Sizeof(int(0)))
var SizeOfInt8 int = int(unsafe.Sizeof(int8(0)))
var SizeOfInt16 int = int(unsafe.Sizeof(int16(0)))
var SizeOfInt32 int = int(unsafe.Sizeof(int32(0)))
var SizeOfInt64 int = int(unsafe.Sizeof(int64(0)))

var SizeOfUint32 int = int(unsafe.Sizeof(uint32(0)))
var SizeOfUint64 int = int(unsafe.Sizeof(uint64(0)))

