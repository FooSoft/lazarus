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

func testRead(test *testing.T, reader *BitReader, count int, expected uint64) {
	value, err := reader.ReadUint64(count)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		test.Fail()
	} else if value != expected {
		log.Printf("value: %.16x, expected: %.16x\n", value, expected)
		test.Fail()
	}
}

func TestUnsigned(t *testing.T) {
	r := NewBitReader(bytes.NewReader(data))

	testRead(t, r, 0, 0x00)
	testRead(t, r, 8, 0x01)
	testRead(t, r, 16, 0x4523)
	testRead(t, r, 3, 0x67&0x07)
	testRead(t, r, 13, 0x8967>>3)
	testRead(t, r, 13, 0xcdab&0x1fff)
	testRead(t, r, 2, (0xcdab>>13)&3)
	testRead(t, r, 9, 0xefcdab>>15)
}

func TestUnsignedSmall(t *testing.T) {
	r := NewBitReader(bytes.NewReader(data))

	testRead(t, r, 0, 0x00)
	testRead(t, r, 8, 0x01)
	testRead(t, r, 8, 0x23)
	testRead(t, r, 8, 0x45)
	testRead(t, r, 16, 0x8967)
}

func TestUnsignedTiny(t *testing.T) {
	r := NewBitReader(bytes.NewReader(data))

	testRead(t, r, 0, 0x00)
	testRead(t, r, 8, 0x01)
	testRead(t, r, 16, 0x4523)
	testRead(t, r, 3, 0x67&0x07)
	testRead(t, r, 13, 0x8967>>3)
	testRead(t, r, 13, 0xCDAB&0x1fff)
	testRead(t, r, 2, (0xcdab>>13)&0x03)
	testRead(t, r, 5, (0xefcdab>>15)&0x1f)
	testRead(t, r, 4, (0xefcdab>>20)&0x0f)
}
