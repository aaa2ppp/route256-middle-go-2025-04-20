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
	desk [][]cell
	axes [][2]int
}

type results struct {
	desks iter.Seq[[][]cell]
}

func readParams(br *bufio.Reader) params {
	var n, m, k int
	if _, err := fmt.Fscanln(br, &n, &m, &k); err != nil {
		panic(err)
	}

	desk := makeMatrix[cell](n, m)
	for i := 0; i < n; i++ {
		row, _, err := br.ReadLine()
		if err != nil {
			panic(err)
		}
		for j, c := range row {
			if c == '#' {
				desk[i][j] = filledCell
			}
		}
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
			for _, c := range row {
				bw.WriteByte(cellToChar[c])
			}
			bw.WriteByte('\n')
		}
		bw.WriteByte('\n')
	}
}

type dir struct {
	i, j int
}

type cell byte

var filledCell = cell(0b1111)

var cellToChar = []byte{
	0b0000: '.',
	0b1111: '#',
	0b0100: '^',
	0b1011: '^',
	0b1000: '>',
	0b0111: '>',
	0b0001: 'v',
	0b1110: 'v',
	0b0010: '<',
	0b1101: '<',
	0b1001: '/',
	0b0110: '/',
	0b0011: '\\',
	0b1100: '\\',
	0b1010: 'x',
	0b0101: 'x',
}

func (c cell) mirror(d dir) cell {
	//    0
	//  3   1  3210 (&0b1111)
	//    2
	switch d {
	case dir{0, 1}, dir{0, -1}:
		//   2
		// 3   1  3012
		//   0
		return (c&0b0001)<<2 | (c&0b0100)>>2 | (c & 0b1010)
	case dir{1, 0}, dir{-1, 0}:
		//   0
		// 1   3  1230
		//   2
		return (c&0b0010)<<2 | (c&0b1000)>>2 | (c & 0b0101)
	case dir{1, 1}, dir{-1, -1}:
		//   3
		// 0   2  0123
		//   1
		return (c&0b0001)<<3 | (c&0b1000)>>3 | (c&0b0010)<<1 | (c&0b0100)>>1
	case dir{1, -1}, dir{-1, 1}:
		//   1
		// 2   0  2301
		//   3
		return (c&0b0001)<<1 | (c&0b0010)>>1 | (c&0b0100)<<1 | (c&0b1000)>>1
	}
	panic(fmt.Errorf("unknown dir %v", d))
}

func (c cell) flex(d dir) cell {
	//    0
	//  3   1  3210 (&0b1111)
	//    2
	switch d {
	case dir{1, 1}:
		//   0
		// ^/
		// 3   1  __10
		//   ^/
		//   2
		return (c&0b1000)>>3 | (c&0b0100)>>1 | (c & 0b0011)
	case dir{-1, -1}:
		//   0
		//  /v
		// 3   1  32__
		//    /v
		//   2
		return (c&0b0001)<<3 | (c&0b0010)<<1 | (c & 0b1100)
	case dir{-1, 1}:
		//   0
		//    \^
		// 3   1  3__0
		//  \^
		//   2
		return (c&0b0010)>>1 | (c&0b0100)<<1 | (c & 0b1001)
	case dir{1, -1}:
		//   0
		//   v\
		// 3   1  _21_
		// v\
		//   2
		return (c&0b0001)<<1 | (c&0b1000)>>1 | (c & 0b0110)
	}
	panic(fmt.Errorf("unknown dir %v", d))
}

type point struct {
	i, j int
}

func flex(p point, c cell, a point, d dir) (point, cell) {
	// здесь мы переходим в координаты x2, чтобы точки осей всегда были четными,
	// а точки клеток нечетными со смещением вправо-вниз (фактически это центры клеток).
	// это позволяет легче определять справа или слева от оси находится клетка

	switch d {
	case dir{0, 1}, dir{0, -1}:
		a_i := a.i * 2
		p_i := p.i*2 + 1

		if sign(p_i-a_i)*d.j == -1 {
			// not right
			return p, c
		}

		p_i = (a_i - (p_i - a_i) - 1) / 2
		return point{p_i, p.j}, c.mirror(d)

	case dir{-1, 0}, dir{1, 0}:
		a_j := a.j * 2
		p_j := p.j*2 + 1

		if sign(p_j-a_j)*d.i == 1 {
			// not right
			return p, c
		}

		p_j = (a_j - (p_j - a_j) - 1) / 2
		return point{p.i, p_j}, c.mirror(d)

	default:
		a_i := a.i * 2
		a_j := a.j * 2
		p_i := p.i*2 + 1
		p_j := p.j*2 + 1

		x_i := a_i + (p_j-a_j)*d.i*d.j // NOTE: всегда |d.i| = |d.j| = 1
		x_j := a_j + (p_i-a_i)*d.i*d.j

		if p_i == x_i {
			// мы на оси
			if debugEnable {
				log.Println("flex", p_i, p_j)
			}
			return p, c.flex(d)
		}

		// XXX просто "затолкал" условия
		if (d.i*d.j == 1) && sign(p_i-x_i)*d.i == -1 {
			// not right
			return p, c
		} else if (d.i*d.j == -1) && sign(x_j-p_j)*d.i == -1 {
			// not right
			return p, c
		}

		// здесь возвращаемся в оригинальные координаты
		return point{(x_i - 1) / 2, (x_j - 1) / 2}, c.mirror(d)
	}
}

func solve(params params) results {
	desk := params.desk
	axes := params.axes
	desks := func(yield func([][]cell) bool) {
		for _, axis := range axes {
			desk = solveStep(desk, axis[0], axis[1])
			if desk == nil {
				return
			}
			if !yield(desk) {
				return
			}
		}
	}
	return results{desks}
}

func solveStep(desk [][]cell, q1, q2 int) [][]cell {
	if debugEnable {
		log.Println("---")
		for _, row := range desk {
			log.Printf("%2d\n", row)
		}
	}

	n, m := len(desk), len(desk[0])

	a1, ok := getAxisPoint(n, m, q1)
	if !ok {
		return nil
	}
	a2, ok := getAxisPoint(n, m, q2)
	if !ok {
		return nil
	}

	if debugEnable {
		log.Println("a:", a1, a2)
	}

	d := dir{
		sign(a2.i - a1.i),
		sign(a2.j - a1.j),
	}

	if debugEnable {
		log.Println("d:", d)
	}

	// загибаем углы, чтобы узнать новые размеры и смещение
	c1, _ := flex(point{0, 0}, 0, a1, d)
	c2, _ := flex(point{0, m - 1}, 0, a1, d)
	c3, _ := flex(point{n - 1, m - 1}, 0, a1, d)
	c4, _ := flex(point{n - 1, 0}, 0, a1, d)

	if debugEnable {
		log.Println("c:", c1, c2, c3, c4)
	}

	min_i := min(c1.i, c2.i, c3.i, c4.i, a1.i, a2.i)
	max_i := max(c1.i, c2.i, c3.i, c4.i, a1.i-1, a2.i-1)
	min_j := min(c1.j, c2.j, c3.j, c4.j, a1.j, a2.j)
	max_j := max(c1.j, c2.j, c3.j, c4.j, a1.j-1, a2.j-1)

	new_n := max_i - min_i + 1
	new_m := max_j - min_j + 1
	new_desk := makeMatrix[cell](new_n, new_m)

	offset_i, offset_j := min_i, min_j
	if debugEnable {
		log.Println("NxM   :", new_n, new_m)
		log.Println("offset:", offset_i, offset_j)
	}

	min_i, max_i = new_n-1, 0
	min_j, max_j = new_m-1, 0

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c := desk[i][j]
			if c == 0 {
				continue
			}

			p := point{i, j}
			p_new, c_new := flex(p, c, a1, d)

			if debugEnable {
				log.Printf("%v(%04b) -> %v(%04b)", p, c, p_new, c_new)
			}

			new_i := p_new.i - offset_i
			new_j := p_new.j - offset_j
			new_desk[new_i][new_j] |= c_new

			min_i = min(min_i, new_i)
			max_i = max(max_i, new_i)
			min_j = min(min_j, new_j)
			max_j = max(max_j, new_j)
		}
	}

	if debugEnable {
		for _, row := range new_desk {
			log.Printf("%2d\n", row)
		}
	}

	new_desk = new_desk[min_i : max_i+1]
	for i := range new_desk {
		new_desk[i] = new_desk[i][min_j : max_j+1]
	}

	return new_desk
}

func sign(a int) int {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}

func getAxisPoint(n, m, q int) (point, bool) {
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
