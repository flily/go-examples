package main

import "strings"

type Checksum interface {
	Reset()
	Update([]byte) error
	Size() int
	BlockSize() int
	Checksum() uint64
}

func NewChecksum(algo string) Checksum {
	switch strings.ToLower(algo) {
	case "net":
		return NewInternetChecksum()

	case "bsd":
		return NewBSDChecksum()

	default:
		return nil
	}
}

type InternetChecksum struct {
	checksum uint64
}

func NewInternetChecksum() *InternetChecksum {
	c := &InternetChecksum{}
	c.Reset()

	return c
}

func (c *InternetChecksum) Reset() {
	c.checksum = 0
}

func (c *InternetChecksum) BlockSize() int {
	return 1024
}

func (c *InternetChecksum) Size() int {
	return 2
}

func (c *InternetChecksum) Update(data []byte) error {
	n := [2]uint64{0, 0}
	for i := 0; i < len(data); i++ {
		n[i%2] += uint64(data[i])
	}

	s := (n[0] << 8) + n[1]
	c.checksum += s
	for c.checksum > 0xffff {
		h := c.checksum >> 16
		l := c.checksum & 0xffff
		c.checksum = h + l
	}

	return nil
}

func (c *InternetChecksum) Checksum() uint64 {
	return (^c.checksum) & 0x0000ffff
}

type BSDChecksum struct {
	checksum uint64
}

func NewBSDChecksum() *BSDChecksum {
	c := &BSDChecksum{}
	c.Reset()

	return c
}

func (c *BSDChecksum) Reset() {
	c.checksum = 0
}

func (c *BSDChecksum) BlockSize() int {
	return 1024
}

func (c *BSDChecksum) Size() int {
	return 2
}

func (c *BSDChecksum) Update(data []byte) error {
	s := c.checksum

	for i := 0; i < len(data); i++ {
		s = (s >> 1) + ((s & 1) << 15)
		s += uint64(data[i])
		s &= 0xffff
	}

	c.checksum = s
	return nil
}

func (c *BSDChecksum) Checksum() uint64 {
	return c.checksum
}
