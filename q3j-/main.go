package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type params struct {
	cityMap []string
	x, y    int
	power   int32
}

type results struct {
	count int
}

func readParams(br *bufio.Reader) params {
	var n, m int
	if _, err := fmt.Fscanln(br, &n, &m); err != nil {
		panic(err)
	}

	cityMap := make([]string, n)
	for i := 0; i < n; i++ {
		row, err := br.ReadString('\n')
		if err != nil {
			panic(err)
		}
		cityMap[i] = row[:m]
	}

	var (
		x, y  int
		power int32
	)
	if _, err := fmt.Fscanln(br, &x, &y, &power); err != nil {
		panic(err)
	}

	return params{cityMap, x - 1, y - 1, power} // to 1-indexing
}

func writeResults(bw *bufio.Writer, results results) {
	fmt.Fprintln(bw, results.count)
}

func createBuildingsMap(cityMap []string) [][]int32 {
	n, m := len(cityMap), len(cityMap[0])

	buildings := makeMatrix[int32](n, m)
	var id int32
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if cityMap[i][j] != '0' && buildings[i][j] == 0 {
				id++
				markBuilding(buildings, cityMap, id, i, j)
			}
		}
	}

	return buildings
}

type point struct {
	i, j int
}

func markBuilding(buildings [][]int32, cityMap []string, id int32, i, j int) {
	n, m := len(buildings), len(buildings[0])

	var queue []point // XXX мне лень реализовывать очередь

	buildings[i][j] = id
	queue = append(queue, point{i, j})

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for _, offset := range []point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			neig := point{node.i + offset.i, node.j + offset.j}
			if !(0 <= neig.i && neig.i < n && 0 <= neig.j && neig.j < m) {
				continue
			}
			if cityMap[neig.i][neig.j] == '0' || buildings[neig.i][neig.j] != 0 {
				continue
			}
			buildings[neig.i][neig.j] = id
			queue = append(queue, neig)
		}
	}
}

func earthquake(buildings [][]int32, cityMap []string, i, j int, power int32) int {
	n, m := len(buildings), len(buildings[0])

	if power <= 1 {
		// никого не сможет разрушить, т.к. минимальня сейсмостойкость 1
		return 0
	}

	damaged := make(map[int32]struct{})
	var queue []point // XXX мне лень реализовывать очередь

	// обрабатываем до помещения фронта волны в очередь, т.к. храним силу волны
	// и id здания в одном и том же месте
	id := buildings[i][j]
	if id != 0 {
		// здесь есть строение, если оно не выстоит, запомним его id
		stability := int32(cityMap[i][j] - '0')
		if power > stability {
			damaged[id] = struct{}{}
		}
	}

	buildings[i][j] = -power // сохраняем отрицательным числом, чтобы отличать от id
	queue = append(queue, point{i, j})

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		// сила с которой волна дойдет до соседних узлов
		power := -buildings[node.i][node.j] - 1
		if power <= 1 {
			// уже никого не сможет разрушить
			continue
		}

		for _, offset := range []point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
			neig := point{node.i + offset.i, node.j + offset.j}

			if !(0 <= neig.i && neig.i < n && 0 <= neig.j && neig.j < m) {
				continue
			}

			id := buildings[neig.i][neig.j]
			if id < 0 {
				// сюда уже дошла волна
				continue
			}
			if id != 0 {
				// здесь есть строение, если оно не выстоит, запомним его id
				stability := int32(cityMap[neig.i][neig.j] - '0')
				if power > stability {
					damaged[id] = struct{}{}
				}
			}

			// гоним волну дальше
			buildings[neig.i][neig.j] = -power
			queue = append(queue, neig)
		}
	}

	return len(damaged)
}

func solve(params params) results {
	if params.power <= 1 {
		// никого не сможет разрушить, т.к. минимальная сейсмостойкость 1
		return results{0}
	}

	// отметим каждый дом уникальным номером > 0
	buildings := createBuildingsMap(params.cityMap)
	if debugEnable {
		log.Println("== buildings")
		for _, row := range buildings {
			log.Printf("%2d\n", row)
		}
	}

	// имитируем землетрясение
	count := earthquake(buildings, params.cityMap, params.x, params.y, params.power)
	if debugEnable {
		log.Println("== earthquake")
		for _, row := range buildings {
			log.Printf("%2d\n", row)
		}
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
