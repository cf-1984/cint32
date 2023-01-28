package cint32

import (
	"fmt"
	"unsafe"
)

// Simple compression of 32-bit integers, which assumes mostly small values.
//
//	Values between -127 and 128 (small ints) compress into 1 byte
//	The excluded edge values of small ints are used as 1 byte prefixes
//	Medium: 128, values betweeen -0x8000 and 0x8000 compress into 2 bytes (3 bytes total)
//	Large:	129 (~ -127), all larger values compress into 4 bytes (5 bytes total)
func Compress[T ~int32](elems ...T) []byte {
	buf := []byte{}
	for _, e := range elems {
		if e > -127 && e < 128 {
			buf = append(buf, byte(e))
			continue
		}
		if e > -0x8000 && e < 0x8000 {
			buf = append(buf, 128, byte(e), byte(e>>8))
			continue
		}
		buf = append(buf, 129, byte(e), byte(e>>8), byte(e>>16), byte(e>>24))
	}
	return buf
}

// Decompresses simple 32-bit compressed bytes.
//
//	Interpretes first byte as prefix, where 128 -> medium integer and 129 -> large integer
//	Any other value is being treated as small integer
//	Small:	1 byte
//	Medium:	3 bytes
//	Large:	5 bytes
func Decompress[T ~uint8](elems ...T) ([]int32, error) {
	size := len(elems)
	if size < 1 {
		return nil, nil
	}

	buf := []int32{}
	for i := 0; i < size; i++ {
		switch e := elems[i:]; e[0] {
		case 129:
			if size-i < 5 {
				return nil, fmt.Errorf("insufficent amount of bytes to decompress: expected 5, received %d", size-i)
			}
			buf = append(buf, int32(e[1])|int32(e[2])<<8|int32(e[3])<<16|int32(utos(e[4]))<<24)
			i += 4
		case 128:
			if size-i < 3 {
				return nil, fmt.Errorf("insufficent amount of bytes to decompress: expected 3, received %d", size-i)
			}
			buf = append(buf, int32(e[1])|int32(utos(e[2]))<<8)
			i += 2
		default:
			buf = append(buf, int32(utos(e[0])))
		}
	}

	return buf, nil
}

// C-style type casting unsigned to signed 8-bit integers
func utos[T ~uint8](i T) int8 {
	return int8(*(*int8)(unsafe.Pointer(&i)))
}
