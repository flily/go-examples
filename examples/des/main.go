package main

import (
	"flag"
	"io"
	"os"
)

func EncryptFile(in io.Reader, out io.Writer, key uint64) error {
	inBlock := make([]byte, 8)
	outBlock := make([]byte, 8)
	cipher := NewDES(key)

	for {
		n, err := in.Read(inBlock)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}

			return nil
		}

		pad := 8 - n
		for i := n; i < 8; i++ {
			inBlock[i] = byte(pad)
		}

		err = cipher.EncryptBlock(inBlock, outBlock, 0)
		if err != nil {
			panic(err)
		}

		_, err = out.Write(outBlock)
		if err != nil {
			panic(err)
		}
	}
}

func findPadding(in []byte) int {
	pad := int(in[7])
	if pad > 7 {
		return 8
	}

	correct := true
	for i := 0; i < pad; i++ {
		if in[7-i] != byte(pad) {
			correct = false
			break
		}
	}

	if correct {
		return 8 - pad
	}

	return 8
}

func DecryptFile(in io.Reader, out io.Writer, key uint64) error {
	inBlock := make([]byte, 8)
	outBlock := make([]byte, 8)
	cipher := NewDES(key)

	var lastBlock []byte
	for {
		n, err := in.Read(inBlock)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}

		if err == nil && n != 8 {
			panic("n != 8")
		}

		if lastBlock != nil {
			pad := 8
			if err != nil {
				pad = findPadding(lastBlock)
			}

			_, errWrite := out.Write(lastBlock[:pad])
			if errWrite != nil {
				panic(err)
			}

			if err != nil {
				return nil
			}
		}

		err = cipher.DecryptBlock(inBlock, outBlock, 0)
		if err != nil {
			panic(err)
		}

		if lastBlock == nil {
			lastBlock = make([]byte, 8)
		}

		copy(lastBlock, outBlock)
	}
}

func main() {
	key := flag.Uint64("key", 0x0011223344556677, "key")
	flag.Bool("encrypt", true, "encrypt")
	isDecrypt := flag.Bool("decrypt", false, "decrypt")
	input := flag.String("in", "", "input file")
	output := flag.String("out", "", "output file")
	flag.Parse()

	in := os.Stdin
	if *input != "" {
		if f, err := os.Open(*input); err != nil {
			panic(err)
		} else {
			in = f
		}
	}

	out := os.Stdout
	if *output != "" {
		if f, err := os.Create(*output); err != nil {
			panic(err)
		} else {
			out = f
		}
	}

	var err error
	if *isDecrypt {
		err = DecryptFile(in, out, *key)

	} else {
		err = EncryptFile(in, out, *key)
	}

	if err != nil {
		panic(err)
	}
}
