package set

import "unsafe"

// Uint8 into slice
func Uint8(a []byte, v uint8) {
	a[0] = v
}

// Uint16 into slice
func Uint16(a []byte, v uint16) {
	_ = a[1]
	a[0] = byte(v >> 8)
	a[1] = byte(v)
}

// Uint16Swapped into slice
func Uint16Swapped(a []byte, v uint16) {
	_ = a[1]
	a[0] = byte(v)
	a[1] = byte(v >> 8)
}

// Uint32 into slice
func Uint32(a []byte, v uint32) {
	_ = a[3]
	Uint16(a[2:], uint16(v))
	Uint16(a[0:], uint16(v>>16))
}

// Uint32Swapped into slice
func Uint32Swapped(a []byte, v uint32) {
	_ = a[3]
	Uint16Swapped(a[0:], uint16(v))
	Uint16Swapped(a[2:], uint16(v>>16))
}

// Float32 into slice
func Float32(a []byte, v float32) {
	_ = a[3]
	Uint32(a, *(*uint32)(unsafe.Pointer(&v)))
}

// Float32Swapped into slice
func Float32Swapped(a []byte, v float32) {
	_ = a[3]
	Uint32Swapped(a, *(*uint32)(unsafe.Pointer(&v)))
}
