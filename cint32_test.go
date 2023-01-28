package cint32_test

import (
	"bytes"
	"testing"

	"github.com/cf-1984/cint32"
)

func equal[T ~int32](slice []T, elems ...T) bool {
	if len(slice) != len(elems) {
		return false
	}

	for i := 0; i < len(slice); i++ {
		if slice[i] != elems[i] {
			return false
		}
	}

	return true
}

func Test_empty_intslice_yields_empty_byteslice(t *testing.T) {
	if yield := cint32.Compress([]int32{}...); len(yield) != 0 {
		t.Errorf("length = %d; want 0", len(yield))
	}
}

func Test_small_int_yields_one_byte(t *testing.T) {
	yield := cint32.Compress(int32(0), int32(-127+1), int32(128-1))
	want := []byte{0, 130, 127}

	if !bytes.Equal(yield, want) {
		t.Errorf("compressed = %v; want %v", yield, want)
	}
}

func Test_medium_int_yields_three_bytes(t *testing.T) {
	yield := cint32.Compress([]int32{
		-127,
		128,
		-0x8000 + 1,
		0x8000 - 1,
	}...)

	want := []byte{
		128, 129, 255,
		128, 128, 0,
		128, 1, 128,
		128, 255, 127,
	}

	if !bytes.Equal(yield, want) {
		t.Errorf("compressed = %v; want %v", yield, want)
	}
}

func Test_large_int_yields_five_bytes(t *testing.T) {
	yield := cint32.Compress([]int32{
		-0x8000,
		0x8000,
	}...)

	want := []byte{
		129, 0, 128, 255, 255,
		129, 0, 128, 0, 0,
	}

	if !bytes.Equal(yield, want) {
		t.Errorf("compressed = %v; want %v", yield, want)
	}
}

func Test_empty_byteslice_yields_empty_intslice(t *testing.T) {
	if yield, err := cint32.Decompress([]byte{}...); err != nil {
		t.Errorf("err != nil; want nil")
	} else if len(yield) != 0 {
		t.Errorf("length = %d; want 0", len(yield))
	}
}

func Test_one_byte_yields_small_int(t *testing.T) {
	want := []int32{0, -126, 127}
	if yield, err := cint32.Decompress(byte(0), 130, 127); err != nil {
		t.Errorf("err != nil; want nil")
	} else if !equal(yield, want...) {
		t.Errorf("decompressed = %v; want %v", yield, want)
	}

}

func Test_too_few_bytes_for_medium_yields_error(t *testing.T) {
	if _, err := cint32.Decompress(byte(128), 2); err == nil {
		t.Errorf("err == nil; want error")
	}
}

func Test_three_bytes_yields_medium_int(t *testing.T) {
	want := []int32{
		-127,
		128,
		-0x8000 + 1,
		0x8000 - 1,
	}

	yield, err := cint32.Decompress(
		uint8(128), 129, 255,
		128, 128, 0,
		128, 1, 128,
		128, 255, 127,
	)

	if err != nil {
		t.Errorf("err != nil; want nil")
	} else if !equal(yield, want...) {
		t.Errorf("decompressed = %v; want %v", yield, want)
	}
}

func Test_too_few_bytes_for_large_yields_error(t *testing.T) {
	if _, err := cint32.Decompress(byte(129), 2, 3, 4); err == nil {
		t.Errorf("err == nil; want error")
	}
}

func Test_five_bytes_yields_large_int(t *testing.T) {
	want := []int32{
		-0x8000,
		0x8000,
	}

	yield, err := cint32.Decompress(
		byte(129), 0, 128, 255, 255,
		129, 0, 128, 0, 0,
	)

	if err != nil {
		t.Errorf("err != nil; want nil")
	} else if !equal(yield, want...) {
		t.Errorf("decompressed = %v; want %v", yield, want)
	}
}

func Test_roundtrip_yields_input_values(t *testing.T) {
	want := []int32{0, 123, 45678, -567, -789456}
	yield, err := cint32.Decompress(cint32.Compress(want...)...)
	if err != nil {
		t.Errorf("err != nil; want nil")
	} else if !equal(yield, want...) {
		t.Errorf("roundtrip = %v; want %v", yield, want)
	}
}
