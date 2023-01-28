package cint32

import (
	"fmt"
	"unsafe"
)

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

func utos[T ~uint8](i T) int8 {
	// C-style type casting unsigned to signed
	return int8(*(*int8)(unsafe.Pointer(&i)))
}
