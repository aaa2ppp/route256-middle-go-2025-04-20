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
	statements []statement
}

type results struct {
	names []string
	X     string
}

type statement struct {
	A  string
	B  string
	is bool
	X  string
}

func parseStatement(s string) statement {

	// «A: I am x!» прибавляет подозреваемому A два очка;
	// «A: I am not x!» отнимает у подозреваемого A одно очко;
	// «A: B is x!» прибавляет подозреваемому B одно очко;
	// «A: B is not x!» отнимает у подозреваемого B одно очко.

	words := splitToWords(s)
	n := len(words)

	return statement{
		A: words[0],
		B: func() string {
			if words[2] == "am" {
				return ""
			} else {
				return words[1]
			}
		}(),

		// Засада!
		// Y: D is not!
		is: words[n-2] != "not",

		X: words[n-1],
	}
}

func splitToWords(s string) []string {
	const punctuation = ".,;:!?"

	words := strings.Split(strings.TrimSpace(s), " ")
	n := 0
	for i, w := range words {
		w = strings.TrimRight(words[i], punctuation)
		if w == "" {
			continue
		}
		words[n] = w
		n++
	}
	return words[:n]
}

func readParams(br *bufio.Reader) params {
	var n int
	if _, err := fmt.Fscanln(br, &n); err != nil {
		panic(err)
	}

	statements := make([]statement, 0, n)
	for i := 0; i < n; i++ {
		s, err := br.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s = strings.TrimSpace(s)
		statements = append(statements, parseStatement(s))
	}

	if debugEnable {
		for _, st := range statements {
			log.Printf("%+v\n", st)
		}
	}

	return params{statements}
}

func writeResults(bw *bufio.Writer, results results) {
	slices.Sort(results.names)
	for _, name := range results.names {
		fmt.Fprintf(bw, "%s is %s.\n", name, results.X)
	}
}

func solve(params params) results {
	scores := make(map[string]int)

	var x string
	for _, st := range params.statements {
		x = st.X

		if st.B == "" { // I am ...
			if st.is {
				scores[st.A] += 2
			} else {
				scores[st.A]--
			}
		} else {
			scores[st.A] += 0
			if st.is {
				scores[st.B]++
			} else {
				scores[st.B]--
			}
		}
	}

	maximum := math.MinInt
	for _, v := range scores {
		maximum = max(maximum, v)
	}

	var names []string
	for k, v := range scores {
		if v == maximum {
			names = append(names, k)
		}

	}
	return results{
		names: names,
		X:     x,
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
