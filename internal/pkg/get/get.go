package get

import "unsafe"

// Uint8 from slice
func Uint8(a []byte) uint8 {
	return a[0]
}

// Uint16 from slice
func Uint16(a []byte) uint16 {
	_ = a[1]
	return uint16(a[1]) | uint16(a[0])<<8
}

// Uint16Swapped from slice
func Uint16Swapped(a []byte) uint16 {
	_ = a[1]
	return uint16(a[0]) | uint16(a[1])<<8
}

// Uint32 from slice
func Uint32(a []byte) uint32 {
	_ = a[3]
	return uint32(Uint16(a[2:])) | uint32(Uint16(a))<<16
}

// Uint32Swapped from slice
func Uint32Swapped(a []byte) uint32 {
	_ = a[3]
	return uint32(Uint16Swapped(a)) | uint32(Uint16Swapped(a[2:]))<<16
}

// Float32 from slice
func Float32(a []byte) float32 {
	x := Uint32(a)
	return *(*float32)(unsafe.Pointer(&x))
}

// Float32Swapped from slice
func Float32Swapped(a []byte) float32 {
	x := Uint32Swapped(a)
	return *(*float32)(unsafe.Pointer(&x))
}
