package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strings"
)

type params struct {
	scores    map[string]int
	statement string
}

type results struct {
	names     []string
	statement string
	maxScore  int
}

func trimLastChar(s string) string {
	if s != "" {
		return s[:len(s)-1]
	}
	return s
}

func readParams(br *bufio.Reader) params {
	scores := make(map[string]int)

	var n int
	if _, err := fmt.Fscanln(br, &n); err != nil {
		panic(err)
	}

	var statement string

	for i := 0; i < n; i++ {
		s, err := br.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s = strings.TrimSpace(s)
		words := strings.Split(s, " ")
		if debugEnable {
			log.Println(s, words)
		}

		// «A: I am x!» прибавляет подозреваемому A два очка;
		// «A: I am not x!» отнимает у подозреваемого A одно очко;
		// «A: B is x!» прибавляет подозреваемому B одно очко;
		// «A: B is not x!» отнимает у подозреваемого B одно очко.

		name := trimLastChar(words[0])
		if _, ok := scores[name]; !ok {
			scores[name] = 0
		}

		statement = trimLastChar(words[len(words)-1])

		if words[2] == "am" {
			// говорит про себя
			if words[3] == "not" {
				// отрицает
				scores[name]--
			} else {
				// подтверждает
				scores[name] += 2
			}
		} else {
			// говорит про другого
			other := words[1]
			if words[3] == "not" {
				// отрицает
				scores[other]--
			} else {
				// подтверждает
				scores[other]++
			}
		}
	}

	return params{
		scores:    scores,
		statement: statement,
	}
}

func writeResults(bw *bufio.Writer, results results) {
	slices.Sort(results.names)
	for _, name := range results.names {
		fmt.Fprintf(bw, "%s is %s.\n", name, results.statement)
	}
}

func solve(params params) results {
	maximum := math.MinInt
	for _, v := range params.scores {
		maximum = max(maximum, v)
	}
	var names []string
	for k, v := range params.scores {
		if v == maximum {
			names = append(names, k)
		}

	}
	return results{
		names:     names,
		statement: params.statement,
		maxScore:  maximum,
	}
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
