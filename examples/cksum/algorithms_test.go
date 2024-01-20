package main

import (
	"testing"
)

func TestInternetChecksum(t *testing.T) {
	// An IPv4 header example presents in https://en.wikipedia.org/wiki/Internet_checksum
	// Which checksum field is 0xb861
	data := []byte{
		0x45, 0x00, 0x00, 0x73,
		0x00, 0x00, 0x40, 0x00,
		0x40, 0x11, 0xb8, 0x61,
		0xc0, 0xa8, 0x00, 0x01,
		0xc0, 0xa8, 0x00, 0xc7,
	}

	c := NewInternetChecksum()
	_ = c.Update(data)
	got := c.Checksum()
	exp := uint64(0)
	if got != exp {
		t.Errorf("got %04x; expected %04x", got, exp)
	}
}
