package main

import (
	"bufio"
	"bytes"
	"fmt"
	"iter"
	"os"
	"slices"
)

type result struct {
	min, max, sum, n int
}

func (r *result) update(x int) {
	if r.n == 0 {
		r.min = x
		r.max = x
	} else {
		r.min = min(r.min, x)
		r.max = max(r.max, x)
	}
	r.sum += x
	r.n++
}

func (r *result) String() string {
	mean := float64(r.sum) / 10.0 / float64(r.n)
	return fmt.Sprintf("%.1f/%.1f/%.1f", float64(r.min)/10.0, mean, float64(r.max)/10.0)
}

type item struct {
	name []byte
	res  *result
}

// empirical
const N = 1 << 17

func hash(b []byte) uint64 {
	var h uint64
	for _, c := range b {
		h = h*31 + uint64(c)
	}
	return h & (N - 1)
}

func filter[T any](m []T, f func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, k := range m {
			if !f(k) {
				continue
			}
			if !yield(k) {
				return
			}
		}
	}
}

func read(data *bufio.Reader) []*item {
	m := make([]*item, N)

	var buf []byte
	var err error
	for {
		// station name
		buf, err = data.ReadSlice(';')
		if err != nil {
			// eof
			break
		}
		// remove the semicolon
		buf = buf[:len(buf)-1]

		idx := hash(buf)
		it := m[idx]
		if it == nil {
			name := make([]byte, len(buf))
			// we need to save the name as buf's underlying slice will change with next call of ReadSlice
			copy(name, buf)
			it = &item{name: name, res: &result{}}
			m[idx] = it
		}

		// temperature
		buf, _ = data.ReadSlice('\n')
		// remove the LF
		buf = buf[:len(buf)-1]

		p := 1
		x := 0
		for i := len(buf) - 1; i >= 0; i-- {
			k := buf[len(buf)-i-1]
			if k == '.' {
				continue
			}
			if k == '-' {
				x *= -1
				break
			}
			x += p * int(k-48)
			p *= 10
		}

		it.res.update(x)
	}

	return slices.SortedFunc(filter(m, func(it *item) bool { return it != nil }), func(a, b *item) int {
		return bytes.Compare(a.name, b.name)
	})
}

func print(results []*item) {
	fmt.Print("{")
	for i, it := range results {
		if i != len(results)-1 {
			fmt.Printf("%s=%v, ", it.name, it.res)
		} else {
			fmt.Printf("%s=%v", it.name, it.res)
		}
	}
	fmt.Println("}")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v path/to/1brc/input/file\n", os.Args[0])
		os.Exit(1)
	}
	path := os.Args[1]
	f, _ := os.Open(path)
	defer f.Close()

	print(read(bufio.NewReaderSize(f, 1<<20)))
}
