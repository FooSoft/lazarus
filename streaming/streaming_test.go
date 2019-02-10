package streaming

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBitReader(t *testing.T) {
	data := []byte{
		0x01, // 00000001
		0x23, // 00100011
		0x45, // 01000101
		0x67, // 01100111
		0x89, // 01100111
		0xAB, // 10101011
		0xCD, // 11001101
		0xEF, // 11101111
	}

	r := NewBitReader(bytes.NewReader(data))

	readPass := func(c int, v uint64) {
		if value, err := r.ReadUint64(c); value != v || err != nil {
			fmt.Printf("%.16x (expected %.16x)\n", value, v)
			t.Fail()
		}
	}

	readPass(0, 0x00)
	readPass(8, 0x01)
	readPass(16, 0x4523)
	readPass(3, 0x67&0x07)
	readPass(13, 0x8967>>3)
	readPass(13, 0xcdab&0x1fff)
	readPass(2, (0xcdab>>13)&3)
	readPass(9, 0xefcdab>>15)
}
