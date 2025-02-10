package wc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Number of characters (-m) = number of runes --> utf8.RuneCount(byte[])
// Number of bytes --> we are reading the file into bytes
// Number of lines --> new line characters \n
// Words: count sequences of non-whitespace characters

type Options struct {
	line, char, word, byte bool
}

type Counts struct {
	Lines, Chars, Words, Bytes int
}

func Count(r io.Reader) (Counts, error) {
	var c Counts
	inWord := false
	buf := bufio.NewReader(r)
	leftover := []byte{}

	for {
		p := make([]byte, 4096)
		n, err := buf.Read(p)
		c.Bytes += n
		c.Lines += bytes.Count(p[:n], []byte{'\n'})

		data := append(leftover, p[:n]...)
		leftover = data[:0]
		for len(data) > 0 {
			r, size := utf8.DecodeRune(data)
			fmt.Print("size: ", size)
			if r == utf8.RuneError && size == 0 {
				leftover = append(leftover, data...)
				fmt.Println("leftover: ", leftover, "data: ", data)
				break
			}

			if unicode.IsSpace(r) {
				inWord = false
			} else if !inWord {
				c.Words++
				inWord = true
			}
			data = data[size:]
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return c, err
		}

	}
	return c, nil
}

func ProcessFile(filename string) (Counts, error) {
	if filename == "-" {
		return Count(io.NopCloser(os.Stdin))
	}

	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return Counts{}, err
	}
	return Count(f)
}

func ParseArgs(args []string) (Options, []string) {
	var opts Options
	var filenames []string

	for _, arg := range args {
		if arg == "--" {
			break
		}

		if len(arg) > 0 && arg[0] == '-' {
			for _, opt := range arg[1:] {
				switch opt {
				case 'l':
					opts.line = true
				case 'c':
					opts.byte = true
				case 'w':
					opts.word = true
				case 'm':
					opts.char = true
				}
			}
		} else {
			filenames = append(filenames, arg)
		}
	}
	return opts, filenames
}

func getOrder(opts Options) []string {
	order := make([]string, 0)
	if opts.line {
		order = append(order, "line")
	}
	if opts.byte {
		order = append(order, "byte")
	}
	if opts.word {
		order = append(order, "word")
	}
	if opts.char {
		order = append(order, "char")
	}
	if len(order) == 0 {
		return []string{"line", "word", "byte"}
	}
	return order
}

func FormatOutput(c Counts, opts Options, filename string) string {
	fields := getOrder(opts)
	var parts []string

	for _, field := range fields {
		switch field {
		case "line":
			parts = append(parts, fmt.Sprintf("%8d", c.Lines))
		case "char":
			parts = append(parts, fmt.Sprintf("%7d", c.Chars))
		case "word":
			parts = append(parts, fmt.Sprintf("%7d", c.Words))
		case "byte":
			parts = append(parts, fmt.Sprintf("%7d", c.Bytes))
		}
	}

	if filename != "" {
		parts = append(parts, filename)
	}
	return strings.Join(parts, " ")
}
