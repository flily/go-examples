package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	MODE_ENCODE = 0
	MODE_DECODE = 1
)

var BASE64_ENCODE_STANDARD_MAP = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', // 0-7
	'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', // 8-15
	'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', // 16-23
	'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', // 24-31
	'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', // 32-39
	'o', 'p', 'q', 'r', 's', 't', 'u', 'v', // 40-47
	'w', 'x', 'y', 'z', '0', '1', '2', '3', // 48-55
	'4', '5', '6', '7', '8', '9', '+', '/', // 56-63
}

var BASE64_ENCODE_URLSAFE_MAP = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', // 0-7
	'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', // 8-15
	'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', // 16-23
	'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', // 24-31
	'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', // 32-39
	'o', 'p', 'q', 'r', 's', 't', 'u', 'v', // 40-47
	'w', 'x', 'y', 'z', '0', '1', '2', '3', // 48-55
	'4', '5', '6', '7', '8', '9', '-', '_', // 56-63
}

var BASE64_DECODE_MAP = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// SP    !     "     #     $     %     &     '     (     )     *     +     ,     -     .     /
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x3f, 0x3e, 0x00, 0x3f,
	// 0     1     2     3     4     5     6     7     8     9     :     ;     <     =     >     ?
	0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// @     A     B     C     D     E     F     G     H     I     J     K     L     M     N     O
	0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e,
	// P     Q     R     S     T     U     V     W     X     Y     Z     [     \     ]     ^     _
	0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x00, 0x00, 0x00, 0x00, 0x3f,
	// `     a     b     c     d     e     f     g     h     i     j     k     l     m     n     o
	0x00, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
	// p     q     r     s     t     u     v     w     x     y     z     {     |     }     ~   DEL
	0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

type LineBreakWriter struct {
	writer *bufio.Writer
	count  int
	width  int
}

func NewLineBreakWriter(writer io.Writer, width int) *LineBreakWriter {
	w := &LineBreakWriter{
		writer: bufio.NewWriter(writer),
		count:  0,
		width:  width,
	}

	return w
}

func (w *LineBreakWriter) ToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	w.writer = bufio.NewWriter(file)
	return nil
}

func (w *LineBreakWriter) Write(data []byte) (int, error) {
	c := 0
	var err error

	for _, b := range data {
		err = w.writer.WriteByte(b)
		if err != nil {
			break
		}

		c += 1
		w.count += 1
		if w.width > 0 && w.count >= w.width {
			err = w.writer.WriteByte('\n')
			if err != nil {
				break
			}
			w.count = 0
		}
	}

	return c, err
}

func (w *LineBreakWriter) Flush() {
	w.writer.Flush()
}

func Base64EncodeFile(in io.Reader, out io.Writer, charmap []byte) error {
	buf := make([]byte, 3)
	encoded := make([]byte, 4)
	defer func() {
		_, _ = out.Write([]byte("\n"))
	}()

	for {
		buf[1] = 0
		buf[2] = 0
		n, err := in.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		encoded[0] = charmap[(buf[0]&0xfc)>>2]
		encoded[1] = charmap[((buf[0]&0x03)<<4)|((buf[1]&0xf0)>>4)]
		encoded[2] = charmap[((buf[1]&0x0f)<<2)|((buf[2]&0xc0)>>6)]
		encoded[3] = charmap[(buf[2]&0x3f)>>0]

		if n < 3 {
			encoded[3] = '='
			if n < 2 {
				encoded[2] = '='
			}
		}

		_, err = out.Write(encoded)
		if err != nil {
			return err
		}
	}
}

func Base64DecodeFile(file io.Reader, output io.Writer) error {
	in := make([]byte, 4)
	out := make([]byte, 3)
	i := 0
	reader := bufio.NewReader(file)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}

			out = make([]byte, 0, 3)
			if i > 0 {
				out = append(out, ((in[0] & 0x3f) << 2))
			}

			if i > 1 {
				out[1] |= ((in[1] & 0x30) >> 4)
				out = append(out, ((in[1] & 0x0f) << 4))
			}

			if i > 2 {
				out[2] |= ((in[2] & 0x3c) >> 2)
				out = append(out, ((in[2] & 0x03) << 6))
			}

			_, err = output.Write(out)
			return err
		}

		v := BASE64_DECODE_MAP[b]
		switch {
		case b == '\n', b == '\r', b == ' ', b == '\t':
			continue

		case b == '=', b == 'A':
		case v == 0:
			return fmt.Errorf("invalid base64 character '%c'", b)
		}

		in[i] = v
		i++

		if i < 4 {
			continue
		}

		out[0] = ((in[0] & 0x3f) << 2) | ((in[1] & 0x30) >> 4)
		out[1] = ((in[1] & 0x0f) << 4) | ((in[2] & 0x3c) >> 2)
		out[2] = ((in[2] & 0x03) << 6) | ((in[3] & 0x3f) >> 0)
		_, err = output.Write(out)
		if err != nil {
			return err
		}
		i = 0
	}
}

func openFile(name string) (io.ReadCloser, error) {
	if name == "-" {
		return os.Stdin, nil
	}

	return os.Open(name)
}

func usage() {
	name := os.Args[0]
	fmt.Printf("Usage: %s [-e | -d] file1 [file2 ...]\n", name)
}

type FileEncodeHandler func(io.Reader) error

func main() {
	modeDecode := flag.Bool("d", false, "decode mode")
	width := flag.Int("b", 0, "width of encoded line, 0 means no line break, usually 64 or 76")
	urlsafe := flag.Bool("u", false, "use URL safe encoding")
	output := flag.String("o", "", "output to file")
	flag.Usage = usage
	flag.Parse()

	charmap := BASE64_ENCODE_STANDARD_MAP
	if *urlsafe {
		charmap = BASE64_ENCODE_URLSAFE_MAP
	}

	out := NewLineBreakWriter(os.Stdout, *width)
	defer out.Flush()
	if *output != "" {
		err := out.ToFile(*output)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}
	}

	fileList := []string{"-"}
	if flag.NArg() > 0 {
		fileList = flag.Args()
	}

	for _, filename := range fileList {
		file, err := openFile(filename)
		if err != nil {
			fmt.Printf("Open file '%s' failed: %s\n", filename, err)
			continue
		}

		defer file.Close()

		if *modeDecode {
			err = Base64DecodeFile(file, out)
		} else {
			err = Base64EncodeFile(file, out, charmap)
		}

		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
	}
}
