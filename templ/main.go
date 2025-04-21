package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type params struct {
	// todo
}

type results struct {
	// todo
}

func readParams(br *bufio.Reader) params {
	// todo
	return params{ /*todo*/ }
}

func writeResults(bw *bufio.Writer, results results) {
	// todo
}

func solve(params params) results {
	// todo
	return results{ /*todo*/ }
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
