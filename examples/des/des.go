package main

import (
	"encoding/binary"
	"errors"
)

var BIT32_TABLE = []uint64{
	0x00000000_00000000, // 0
	0x00000000_80000000, // 1
	0x00000000_40000000, // 2
	0x00000000_20000000, // 3
	0x00000000_10000000, // 4
	0x00000000_08000000, // 5
	0x00000000_04000000, // 6
	0x00000000_02000000, // 7
	0x00000000_01000000, // 8
	0x00000000_00800000, // 9
	0x00000000_00400000, // 10
	0x00000000_00200000, // 11
	0x00000000_00100000, // 12
	0x00000000_00080000, // 13
	0x00000000_00040000, // 14
	0x00000000_00020000, // 15
	0x00000000_00010000, // 16
	0x00000000_00008000, // 17
	0x00000000_00004000, // 18
	0x00000000_00002000, // 19
	0x00000000_00001000, // 20
	0x00000000_00000800, // 21
	0x00000000_00000400, // 22
	0x00000000_00000200, // 23
	0x00000000_00000100, // 24
	0x00000000_00000080, // 25
	0x00000000_00000040, // 26
	0x00000000_00000020, // 27
	0x00000000_00000010, // 28
	0x00000000_00000008, // 29
	0x00000000_00000004, // 30
	0x00000000_00000002, // 31
	0x00000000_00000001, // 32
}

var BIT64_TABLE = []uint64{
	0x0000000000000000,
	0x8000000000000000, // 1
	0x4000000000000000, // 2
	0x2000000000000000, // 3
	0x1000000000000000, // 4
	0x0800000000000000, // 5
	0x0400000000000000, // 6
	0x0200000000000000, // 7
	0x0100000000000000, // 8
	0x0080000000000000, // 9
	0x0040000000000000, // 10
	0x0020000000000000, // 11
	0x0010000000000000, // 12
	0x0008000000000000, // 13
	0x0004000000000000, // 14
	0x0002000000000000, // 15
	0x0001000000000000, // 16
	0x0000800000000000, // 17
	0x0000400000000000, // 18
	0x0000200000000000, // 19
	0x0000100000000000, // 20
	0x0000080000000000, // 21
	0x0000040000000000, // 22
	0x0000020000000000, // 23
	0x0000010000000000, // 24
	0x0000008000000000, // 25
	0x0000004000000000, // 26
	0x0000002000000000, // 27
	0x0000001000000000, // 28
	0x0000000800000000, // 29
	0x0000000400000000, // 30
	0x0000000200000000, // 31
	0x0000000100000000, // 32
	0x0000000080000000, // 33
	0x0000000040000000, // 34
	0x0000000020000000, // 35
	0x0000000010000000, // 36
	0x0000000008000000, // 37
	0x0000000004000000, // 38
	0x0000000002000000, // 39
	0x0000000001000000, // 40
	0x0000000000800000, // 41
	0x0000000000400000, // 42
	0x0000000000200000, // 43
	0x0000000000100000, // 44
	0x0000000000080000, // 45
	0x0000000000040000, // 46
	0x0000000000020000, // 47
	0x0000000000010000, // 48
	0x0000000000008000, // 49
	0x0000000000004000, // 50
	0x0000000000002000, // 51
	0x0000000000001000, // 52
	0x0000000000000800, // 53
	0x0000000000000400, // 54
	0x0000000000000200, // 55
	0x0000000000000100, // 56
	0x0000000000000080, // 57
	0x0000000000000040, // 58
	0x0000000000000020, // 59
	0x0000000000000010, // 60
	0x0000000000000008, // 61
	0x0000000000000004, // 62
	0x0000000000000002, // 63
	0x0000000000000001, // 64
}

var PADDING = []uint64{
	0x0000000000000000,
	0x0000000000000001,
	0x0000000000000202,
	0x0000000000030303,
	0x0000000004040404,
	0x0000000505050505,
	0x0000060606060606,
	0x0007070707070707,
}

var (
	IP_TABLE = []int{
		58, 50, 42, 34, 26, 18, 10, 2,
		60, 52, 44, 36, 28, 20, 12, 4,
		62, 54, 46, 38, 30, 22, 14, 6,
		64, 56, 48, 40, 32, 24, 16, 8,
		57, 49, 41, 33, 25, 17, 9, 1,
		59, 51, 43, 35, 27, 19, 11, 3,
		61, 53, 45, 37, 29, 21, 13, 5,
		63, 55, 47, 39, 31, 23, 15, 7,
	}

	IIP_TABLE = []int{
		40, 8, 48, 16, 56, 24, 64, 32,
		39, 7, 47, 15, 55, 23, 63, 31,
		38, 6, 46, 14, 54, 22, 62, 30,
		37, 5, 45, 13, 53, 21, 61, 29,
		36, 4, 44, 12, 52, 20, 60, 28,
		35, 3, 43, 11, 51, 19, 59, 27,
		34, 2, 42, 10, 50, 18, 58, 26,
		33, 1, 41, 9, 49, 17, 57, 25,
	}

	PC1 = []int{
		57, 49, 41, 33, 25, 17, 9,
		1, 58, 50, 42, 34, 26, 18,
		10, 2, 59, 51, 43, 35, 27,
		19, 11, 3, 60, 52, 44, 36,
		63, 55, 47, 39, 31, 23, 15,
		7, 62, 54, 46, 38, 30, 22,
		14, 6, 61, 53, 45, 37, 29,
		21, 13, 5, 28, 20, 12, 4,
	}

	PC2 = []int{
		14, 17, 11, 24, 1, 5,
		3, 28, 15, 6, 21, 10,
		23, 19, 12, 4, 26, 8,
		16, 7, 27, 20, 13, 2,
		41, 52, 31, 37, 47, 55,
		30, 40, 51, 45, 33, 48,
		44, 49, 39, 56, 34, 53,
		46, 42, 50, 36, 29, 32,
	}

	IterateShiftTable = []int{
		1, 1, 2, 2, 2, 2, 2, 2,
		1, 2, 2, 2, 2, 2, 2, 1,
	}

	E = []int{
		32, 1, 2, 3, 4, 5,
		4, 5, 6, 7, 8, 9,
		8, 9, 10, 11, 12, 13,
		12, 13, 14, 15, 16, 17,
		16, 17, 18, 19, 20, 21,
		20, 21, 22, 23, 24, 25,
		24, 25, 26, 27, 28, 29,
		28, 29, 30, 31, 32, 1,
	}

	P = []int{
		16, 7, 20, 21,
		29, 12, 28, 17,
		1, 15, 23, 26,
		5, 18, 31, 10,
		2, 8, 24, 14,
		32, 27, 3, 9,
		19, 13, 30, 6,
		22, 11, 4, 25,
	}

	SBox = [][]uint64{
		{
			14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7,
			0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8,
			4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0,
			15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13,
		},
		{
			15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10,
			3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5,
			0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15,
			13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9,
		},
		{
			10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8,
			13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1,
			13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7,
			1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 1,
		},
		{
			7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15,
			13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9,
			10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4,
			3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 1,
		},
		{
			2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9,
			14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6,
			4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14,
			11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3,
		},
		{
			12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11,
			10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8,
			9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6,
			4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13,
		},
		{
			4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1,
			13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6,
			1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2,
			6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12,
		},
		{
			13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7,
			1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2,
			7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8,
			2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11,
		},
	}
)

func permutation(data uint64, size int, n []int) uint64 {
	outSize := len(n)
	result := uint64(0)
	for i := 0; i < outSize; i++ {
		j := n[i]
		bit := (data >> (size - j)) & 1
		result |= bit << (outSize - i - 1)
	}

	return result
}

func permutationIf(data uint64, size int, n []int) uint64 {
	outSize := len(n)
	result := uint64(0)
	for i := 0; i < outSize; i++ {
		j := n[i]
		if data&(1<<(size-j)) != 0 {
			result |= 1 << (outSize - i - 1)
		}
	}

	return result
}

func leftShift28(data uint64, count int) uint64 {
	if count == 1 {
		return ((data & 0x800_0000) >> 27) | ((data & 0x7ff_ffff) << 1)
	}

	return ((data & 0xc00_0000) >> 26) | ((data & 0x3ff_ffff) << 2)
}

func makeKeys(key64 uint64, subKeys48 []uint64) {
	pcKey56 := permutation(key64, 64, PC1)
	c28, d28 := (pcKey56>>28)&0x0fffffff, (pcKey56>>0)&0x0fffffff
	for i := 0; i < 16; i++ {
		c28 = leftShift28(c28, IterateShiftTable[i])
		d28 = leftShift28(d28, IterateShiftTable[i])
		cd56 := (c28 << 28) | d28
		subKey48 := permutation(cd56, 56, PC2)
		subKeys48 = append(subKeys48, subKey48)
	}
}

func desIP(data uint64) uint64 {
	return permutation(data, 64, IP_TABLE)
}

func desIIP(data uint64) uint64 {
	return permutation(data, 64, IIP_TABLE)
}

func desS(n6 uint64, box int) uint64 {
	return SBox[box][n6]
}

func desSDataSplit(data48 uint64) []uint64 {
	data6 := make([]uint64, 8)
	for i := 0; i < 8; i++ {
		data6[i] = (data48 >> (42 - (i * 6))) & 0x3f
	}

	return data6
}

func desSDataCombine(data4 []uint64) uint64 {
	data32 := uint64(0)
	for i := 0; i < 8; i++ {
		data32 |= data4[i] << ((7 - i) * 4)
	}

	return data32
}

func desF(rData32 uint64, subKey48 uint64) uint64 {
	eData48 := permutation(rData32, 32, E)
	keyData48 := eData48 ^ subKey48

	data4 := make([]uint64, 8)
	data6 := desSDataSplit(keyData48)
	for i := 0; i < 8; i++ {
		data4[i] = desS(data6[i], i)
	}

	data32 := desSDataCombine(data4)
	pData32 := permutation(data32, 32, P)
	return pData32
}

func desEncryptBlockUint(data64 uint64, key64 uint64) uint64 {
	ipData64 := desIP(data64)
	subKeys48 := make([]uint64, 16)
	makeKeys(key64, subKeys48)

	dataL32, dataR32 := (ipData64>>32)&0xffffffff, (ipData64>>0)&0xffffffff
	for i := 0; i < 16; i++ {
		nextL32 := dataR32
		nextR32 := dataL32 ^ desF(dataR32, subKeys48[i])
		dataL32, dataR32 = nextL32, nextR32
	}

	finalData64 := (dataL32 << 32) | dataR32
	return desIIP(finalData64)
}

func desDecryptBlockUint(data64 uint64, key64 uint64) uint64 {
	ipData64 := desIP(data64)
	subKeys48 := make([]uint64, 16)
	makeKeys(key64, subKeys48)

	dataL32, dataR32 := (ipData64>>32)&0xffffffff, (ipData64>>0)&0xffffffff
	for i := 0; i < 16; i++ {
		nextR32 := dataL32
		nextL32 := dataR32 ^ desF(dataL32, subKeys48[15-i])
		dataL32, dataR32 = nextL32, nextR32
	}

	finalData64 := (dataL32 << 32) | dataR32
	return desIIP(finalData64)
}

type DES struct {
	subKeys48 []uint64
}

func NewDES(key64 uint64) *DES {
	des := &DES{}
	des.subKeys48 = make([]uint64, 16)
	makeKeys(key64, des.subKeys48)
	return des
}

func (d *DES) EncryptUint64(data64 uint64) uint64 {
	ipData64 := desIP(data64)

	dataL32, dataR32 := (ipData64>>32)&0xffffffff, (ipData64>>0)&0xffffffff
	for i := 0; i < 16; i++ {
		nextL32 := dataR32
		nextR32 := dataL32 ^ desF(dataR32, d.subKeys48[i])
		dataL32, dataR32 = nextL32, nextR32
	}

	finalData64 := (dataL32 << 32) | dataR32
	return desIIP(finalData64)
}

func (d *DES) EncryptBlock(in []byte, out []byte, offset int) error {
	if len(in) != len(out) {
		return errors.New("len(in) !=len(out)")
	}

	if offset+8 > len(in) {
		return errors.New("offset+8 > len(in)")
	}

	inData := binary.BigEndian.Uint64(in[offset : offset+8])
	outData := d.EncryptUint64(inData)
	binary.BigEndian.PutUint64(out[offset:offset+8], outData)
	return nil
}

func (d *DES) DecryptUint64(data64 uint64) uint64 {
	ipData64 := desIP(data64)

	dataL32, dataR32 := (ipData64>>32)&0xffffffff, (ipData64>>0)&0xffffffff
	for i := 0; i < 16; i++ {
		nextR32 := dataL32
		nextL32 := dataR32 ^ desF(dataL32, d.subKeys48[15-i])
		dataL32, dataR32 = nextL32, nextR32
	}

	finalData64 := (dataL32 << 32) | dataR32
	return desIIP(finalData64)
}

func (d *DES) DecryptBlock(in []byte, out []byte, offset int) error {
	if len(in) != len(out) {
		return errors.New("len(in) !=len(out)")
	}

	if offset+8 > len(in) {
		return errors.New("offset+8 > len(in)")
	}

	inData := binary.BigEndian.Uint64(in[offset : offset+8])
	outData := d.DecryptUint64(inData)
	binary.BigEndian.PutUint64(out[offset:offset+8], outData)
	return nil
}

func DesEncryptBlock(in []byte, out []byte, offset int, key uint64) error {
	if len(in) != len(out) {
		return errors.New("len(in) !=len(out)")
	}

	if offset+8 > len(in) {
		return errors.New("offset+8 > len(in)")
	}

	inData := binary.BigEndian.Uint64(in[offset : offset+8])
	outData := desEncryptBlockUint(inData, key)
	binary.BigEndian.PutUint64(out[offset:offset+8], outData)
	return nil
}

func DesDecryptBlock(in []byte, out []byte, offset int, key uint64) error {
	if len(in) != len(out) {
		return errors.New("len(in) !=len(out)")
	}

	if offset+8 > len(in) {
		return errors.New("offset+8 > len(in)")
	}

	inData := binary.BigEndian.Uint64(in[offset : offset+8])
	outData := desDecryptBlockUint(inData, key)
	binary.BigEndian.PutUint64(out[offset:offset+8], outData)
	return nil
}
