package streaming

import (
	"bytes"
	"testing"
)

func TestBitReader(t *testing.T) {
	data := []byte{
		0x69, // 01101001
		0x96, // 10010110
		0xf0, // 11110000
		0xaa, // 10101010
		0x00, // 00000000
		0xff, // 11111111
	}

	r := NewBitReader(bytes.NewReader(data))

	readPass := func(c int, v uint64) {
		if value, err := r.ReadBitsUnsigned(c); value != v || err != nil {
			t.Fail()
		}
	}

	readFail := func(c int) {
		if value, err := r.ReadBitsUnsigned(c); value != 0 || err == nil {
			t.Fail()
		}
	}

	readPass(0, 0x00)
	readFail(65)
	readPass(2, 0x01)
	readPass(2, 0x02)
	readPass(3, 0x04)
	readPass(1, 0x01)
	readPass(12, 0x096f)
	readPass(8, 0x000a)
	readPass(20, 0x0a00ff)
	readFail(1)
}
