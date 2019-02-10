package streaming

import (
	"bytes"
	"log"
	"testing"
)

var data = []byte{
	0x01, // 00000001
	0x23, // 00100011
	0x45, // 01000101
	0x67, // 01100111
	0x89, // 01100111
	0xAB, // 10101011
	0xCD, // 11001101
	0xEF, // 11101111
}

func TestBitReader(t *testing.T) {
	r := NewBitReader(bytes.NewReader(data))

	test := func(count int, expected uint64) {
		if value := r.ReadUint(count); value != expected {
			log.Printf("value: %.16x, expected: %.16x\n", value, expected)
			t.Fail()
		}

		if err := r.Error(); err != nil {
			log.Printf("error: %s\n", err.Error())
			t.Fail()
		}
	}

	test(0, 0x00)
	test(8, 0x01)
	test(16, 0x4523)
	test(3, 0x67&0x07)
	test(13, 0x8967>>3)
	test(13, 0xcdab&0x1fff)
	test(2, (0xcdab>>13)&3)
	test(9, 0xefcdab>>15)
}
