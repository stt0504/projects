//go:build !solution

package otp

import (
	"errors"
	"io"
)

type CipherReader struct {
	r    io.Reader
	prng io.Reader
}

type CipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func (c CipherReader) Read(p []byte) (int, error) {
	var count1, count2 int
	var err1, err2 error
	count1, err1 = c.r.Read(p)
	if err1 != nil && count1 == 0 {
		return 0, err1
	}
	var randomNumber = make([]byte, count1)
	count2, err2 = c.prng.Read(randomNumber)
	if err2 != nil && count2 == 0 {
		return 0, err2
	}
	for i := 0; i < min(count1, count2); i++ {
		p[i] = p[i] ^ randomNumber[i]
	}
	return min(count1, count2), nil
}

func (c *CipherWriter) Write(p []byte) (int, error) {
	var err error
	p2 := make([]byte, len(p))
	copy(p2, p)
	var randomNumber = make([]byte, len(p2))
	count1, err1 := c.prng.Read(randomNumber)
	if err1 != nil {
		if count1 == 0 {
			return 0, err1
		}
		err = errors.Join(err1)
	}
	for i := 0; i < min(count1, len(p2)); i++ {
		p2[i] = p2[i] ^ randomNumber[i]
	}
	count2, err2 := c.w.Write(p2)
	if err2 != nil {
		err = errors.Join(err2)
	}
	return count2, err
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &CipherReader{r, prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &CipherWriter{w, prng}
}
