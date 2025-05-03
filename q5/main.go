package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
)

type params struct {
	desk [][]byte
	axes [][2]int
}

type results struct {
	desks iter.Seq[[][]byte]
}

func readParams(br *bufio.Reader) params {
	var n, m, k int
	if _, err := fmt.Fscanln(br, &n, &m, &k); err != nil {
		panic(err)
	}

	desk := make([][]byte, 0, n)
	for i := 0; i < n; i++ {
		row, err := br.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		desk = append(desk, row[:m])
	}

	axes := make([][2]int, 0, k)
	for i := 0; i < k; i++ {
		var q1, q2 int
		if _, err := fmt.Fscanln(br, &q1, &q2); err != nil {
			panic(err)
		}
		axes = append(axes, [2]int{q1 - 1, q2 - 1}) // to 0-indexing
	}

	return params{desk, axes}
}

func writeResults(bw *bufio.Writer, results results) {
	for desk := range results.desks {
		for _, row := range desk {
			bw.Write(row)
			bw.WriteByte('\n')
		}
		bw.WriteByte('\n')
	}
}

func solve(params params) results {
	desk := params.desk
	axes := params.axes
	desks := func(yield func([][]byte) bool) {
		for _, axis := range axes {
			desk = solveStep(desk, axis[0], axis[1])
			if !yield(desk) {
				break
			}
		}
	}
	return results{desks}
}

func solveStep(desk [][]byte, q1, q2 int) [][]byte {
	n, m := len(desk), len(desk[0])

	a1, ok := getQPoint(n, m, q1)
	if !ok {
		return nil
	}
	a2, ok := getQPoint(n, m, q2)
	if !ok {
		return nil
	}

	if debugEnable {
		log.Println("as:", a1, a2)
	}

	var mirror func(point) (point, bool)
	switch {
	case a1.i == a2.i:
		dir := sign(a2.j - a1.j)
		if debugEnable {
			log.Println("dir:", dir)
		}
		mirror = func(p point) (point, bool) {
			return mirrorHorizontal(a1.i, dir, p)
		}

	case a1.j == a2.j:
		dir := sign(a2.i - a1.i)
		if debugEnable {
			log.Println("dir:", dir)
		}
		mirror = func(p point) (point, bool) {
			return mirrorVertical(a1.j, dir, p)
		}

	default:
		// todo
		return desk
	}

	c1, _ := mirror(point{0, 0})
	c2, _ := mirror(point{0, m - 1})
	c3, _ := mirror(point{n - 1, m - 1})
	c4, _ := mirror(point{n - 1, 0})

	if debugEnable {
		log.Println("cs:", c1, c2, c3, c4)
	}

	i_min := min(c1.i, c2.i, c3.i, c4.i)
	i_max := max(c1.i, c2.i, c3.i, c4.i)
	j_min := min(c1.j, c2.j, c3.j, c4.j)
	j_max := max(c1.j, c2.j, c3.j, c4.j)

	n_new := i_max - i_min + 1
	m_new := j_max - j_min + 1

	desk_new := makeMatrix[byte](n_new, m_new)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			p, ok := mirror(point{i, j})
			if debugEnable {
				log.Printf("%v -> %v, %v", point{i, j}, p, ok)
			}
			_ = ok
			i_new := p.i - i_min
			j_new := p.j - j_min
			if desk[i][j] == '#' || desk_new[i_new][j_new] == '#' {
				desk_new[i_new][j_new] = '#'
			} else {
				desk_new[i_new][j_new] = '.'
			}
		}
	}

	return desk_new
}

func sign(a int) int {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}

func mirrorHorizontal(a_i, dir int, p point) (point, bool) {
	a_i *= 2
	i := p.i*2 + 1

	if sign(i-a_i)*dir != 1 {
		// not is right
		return p, false
	}

	i = (a_i - (i - a_i) - 1) / 2
	return point{i, p.j}, true
}

func mirrorVertical(a_j, dir int, p point) (point, bool) {
	a_j *= 2
	j := p.j*2 + 1

	if sign(j-a_j)*dir != -1 {
		// not is right
		return p, false
	}

	j = (a_j - (j - a_j) - 1) / 2
	return p, true
}

func mirrorDiagonal(a point, dir point, p point) (point, bool) {
	return p, false // todo
}

type point struct {
	i, j int
}

func getQPoint(n, m, q int) (point, bool) {
	if q < m {
		return point{0, q}, true
	}

	q -= m
	if q < n {
		return point{q, m}, true
	}

	q -= n
	if q < m {
		return point{n, m - q}, true
	}

	q -= m
	if q < n {
		return point{n - q, 0}, true
	}

	return point{}, false
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
