//go:build !solution

package externalsort

import (
	"container/heap"
	"io"
	"os"
	"sort"
	"strings"
)

type Reader struct {
	ir io.Reader
}

func (r Reader) ReadLine() (l string, err error) {
	p := make([]byte, 1)
	l = ""
	var num int
	for p[0] != '\n' {
		num, err = r.ir.Read(p)
		if err == nil || err == io.EOF && num == 1 {
			if p[0] != '\n' {
				l += string(p[0])
			}
		} else {
			return l, err
		}
	}
	return l, nil
}

type Writer struct {
	iw io.Writer
}

func (w Writer) Write(l string) error {
	_, err := w.iw.Write([]byte(l + "\n"))
	return err
}

func NewReader(r io.Reader) LineReader {
	return &Reader{ir: r}
}

func NewWriter(w io.Writer) LineWriter {
	return &Writer{iw: w}
}

type Item struct {
	readerIndex int
	line        string
}

type Heap []Item

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	return h[i].line < h[j].line
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(Item))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := &Heap{}
	heap.Init(h)

	var err error
	var line string
	for i, reader := range readers {
		line, err = reader.ReadLine()
		if err == nil || err == io.EOF {
			if !(err == io.EOF && line == "") {
				heap.Push(h, Item{i, line})
			}
		} else {
			return err
		}
	}

	for h.Len() > 0 {
		popped := heap.Pop(h).(Item)
		err = w.Write(popped.line)
		if err != nil {
			return err
		}
		line, err = readers[popped.readerIndex].ReadLine()
		if err == nil || (err == io.EOF) {
			if !(err == io.EOF && line == "") {
				heap.Push(h, Item{popped.readerIndex, line})
			}
		} else {
			return err
		}
	}
	return nil
}

func Sort(w io.Writer, in ...string) error {
	readers := make([]LineReader, 0, len(in))
	var err error

	for _, fileName := range in {
		var file *os.File
		file, err = os.Open(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		var lines []string
		r := NewReader(file)

		for {
			var line string
			line, err = r.ReadLine()
			if err == io.EOF {
				if line != "" {
					lines = append(lines, line)
				}
				break
			} else if err != nil {
				return err
			}
			lines = append(lines, line)
		}

		sort.Strings(lines)

		var str string
		for _, line := range lines {
			str += line + "\n"
		}
		readers = append(readers, NewReader(strings.NewReader(str)))
	}

	return Merge(NewWriter(w), readers...)
}
