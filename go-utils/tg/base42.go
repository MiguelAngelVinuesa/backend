package tg

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

func RandomBase42(size int) string {
	if size > 32 {
		size = 32
	}
	buf := make([]byte, 0, 32)[:size]
	for _, err := rand.Read(buf); err != nil; {
		time.Sleep(100 * time.Microsecond)
	}
	return ToBase42(buf)
}

func RandomPassword() string {
	buf := make([]byte, 8)
	for _, err := rand.Read(buf); err != nil; {
		time.Sleep(100 * time.Microsecond)
	}

	out := ToBase42(buf)

	for ix := 0; ix < 2+rand.Intn(2); ix++ {
		l1 := len(out) - 4
		ix1 := 2 + rand.Intn(l1)
		ix2 := rand.Intn(chars1L)
		out = out[:ix1] + chars1[ix2:ix2+1] + out[ix1:]
	}

	for ix := 0; ix < 2+rand.Intn(2); ix++ {
		l1 := len(out) - 4
		ix1 := 2 + rand.Intn(l1)
		ix2 := rand.Intn(chars2L)
		out = out[:ix1] + chars2[ix2:ix2+1] + out[ix1:]
	}

	return out
}

const (
	chars1  = "#$%*()-+=<>[]"
	chars1L = len(chars1)
	chars2  = "0123456789"
	chars2L = len(chars2)
)

// ToBase42 convert a slice of bytes into a base42 string.
func ToBase42(in []byte) string {
	buf := bytes.Buffer{}
	max := len(in) - 1
	for ix := 0; ix <= max; ix++ {
		single := true
		b := uint16(in[ix])
		if ix < max {
			ix++
			b |= uint16(in[ix]) << 8
			single = false
		}
		buf.WriteByte(to42[b%42])
		if n := (b / 42) % 42; !single || n > 0 {
			buf.WriteByte(to42[n])
		}
		if !single {
			buf.WriteByte(to42[b/42/42])
		}
	}
	return string(buf.Bytes())
}

// FromBase42 converts a base42 string to a slice of bytes.
// It will return an error if the input contains invalid characters.
func FromBase42(in string) ([]byte, error) {
	if in == "" {
		return []byte{}, nil
	}

	buf := bytes.Buffer{}
	max := len(in)

	buf.Grow((max + 2) * 2 / 3)

	max--
	for ix := 0; ix <= max; ix++ {
		single := true
		i := uint16(from42[in[ix]])
		if i > 42 {
			goto error
		}

		if ix < max {
			ix++
			n := uint16(from42[in[ix]])
			if n > 42 {
				goto error
			}
			i += n * 42
		}

		if ix < max {
			ix++
			n := uint16(from42[in[ix]])
			if n > 42 {
				goto error
			}
			i += n * 42 * 42
			single = false

		}

		buf.WriteByte(uint8(i & 0xff))
		if !single {
			buf.WriteByte(uint8(i >> 8))
		}
	}

	return buf.Bytes(), nil

error:
	return nil, fmt.Errorf("FormBase42: invalid input [%s]", in)
}

var (
	to42 = []byte("346789ABCDEFGHJKLMNPRTWXYabcdefghjknprstxy")
)

var (
	from42 = [256]byte{
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 0, 1, 255, 2, 3, 4, 5, 255, 255, 255, 255, 255, 255,
		255, 6, 7, 8, 9, 10, 11, 12, 13, 255, 14, 15, 16, 17, 18, 255,
		19, 255, 20, 255, 21, 255, 255, 22, 23, 24, 255, 255, 255, 255, 255, 255,
		255, 25, 26, 27, 28, 29, 30, 31, 32, 255, 33, 34, 255, 255, 35, 255,
		36, 255, 37, 38, 39, 255, 255, 255, 40, 41, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	}
)
