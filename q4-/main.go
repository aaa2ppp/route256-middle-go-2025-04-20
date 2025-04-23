package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type params struct {
	histograms [][]int32
}

type results struct {
	count int
}

func readParams(br *bufio.Reader) params {
	var n, m int
	if _, err := fmt.Fscanln(br, &n, &m); err != nil {
		panic(err)
	}

	histograms := makeMatrix[int32](n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var v int32
			if _, err := fmt.Fscan(br, &v); err != nil {
				panic(err)
			}
			histograms[i][j] = v
		}
		if _, err := br.ReadSlice('\n'); err != nil {
			panic(err)
		}
	}

	return params{histograms}
}

func writeResults(bw *bufio.Writer, results results) {
	fmt.Fprintln(bw, results.count)
}

// getBorder возвращает границу и границу сопрягаемой гистограммы,
// в удобном для использования в качестве ключа мапы виде. Мне показалось
// самым простым упаковать массив интов в строку.
func getBorder(histogram []int32) (string, string) {
	m := len(histogram)

	minimum := int32(math.MaxInt32)
	maximum := int32(0)
	for _, v := range histogram {
		minimum = min(minimum, v)
		maximum = max(maximum, v)
	}

	var positive strings.Builder
	positive.Grow(m * 4)
	for i := 0; i < m; i++ {
		v := histogram[i] - minimum
		for j := 0; j < 4; j++ {
			positive.WriteByte(byte(v & 0xff))
			v >>= 8
		}
	}

	var negative strings.Builder
	negative.Grow(m * 4)
	for i := m - 1; i >= 0; i-- {
		v := maximum - histogram[i]
		for j := 0; j < 4; j++ {
			negative.WriteByte(byte(v & 0xff))
			v >>= 8
		}
	}

	return positive.String(), negative.String()
}

func solve(params params) results {
	count := 0
	hashes := make(map[string]int)

	for _, histogram := range params.histograms {
		positive, negative := getBorder(histogram)
		count += hashes[negative]
		hashes[positive]++
	}

	return results{count}
}

// ----------------------------------------------------------------------------

type solveFunc func(params) results

func runTask(br *bufio.Reader, bw *bufio.Writer, solve solveFunc) {
	writeResults(bw, solve(readParams(br)))
}

func run(r io.Reader, w io.Writer, solve solveFunc) {
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	var t int
	if _, err := fmt.Fscanln(br, &t); err != nil {
		panic(err)
	}
	for i := 0; i < t; i++ {
		runTask(br, bw, solve)
	}
}

var debugEnable bool

func main() {
	_, debugEnable = os.LookupEnv("DEBUG")
	run(os.Stdin, os.Stdout, solve)
}

// ----------------------------------------------------------------------------

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
