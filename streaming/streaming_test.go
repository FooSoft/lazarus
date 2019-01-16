package streaming

import (
	"bytes"
	"fmt"
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

	reader := NewReader(bytes.NewReader(data))
	fmt.Println(reader.ReadBits(2))
}
