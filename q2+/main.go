package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type params struct {
	target     string
	statements []string
}

type results struct {
	targetAge int
}

func trimLastChar(s string) string {
	if s != "" {
		return s[:len(s)-1]
	}
	return s
}

func readParams(br *bufio.Reader) params {
	s, err := br.ReadString('\n')
	if err != nil {
		panic(err)
	}
	s = strings.TrimSpace(s)
	words := strings.Split(s, " ")

	// Вопрос, который вы задали себе, выглядит следующим образом:
	// «How old is Ai?» — Cколько лет другу по имени Ai?
	target := trimLastChar(words[3])

	statements := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		s, err := br.ReadString('\n')
		if err != nil {
			panic(err)
		}
		statements = append(statements, strings.TrimSpace(s))
	}

	return params{
		target:     target,
		statements: statements,
	}
}

func writeResults(bw *bufio.Writer, results results) {
	fmt.Fprintln(bw, results.targetAge)
}

func solve(params params) results {
	statements := params.statements

	removeStatement := func(i int) {
		n := len(statements)
		statements[i] = statements[n-1]
		statements = statements[:n-1]
	}

	ages := make(map[string]int)

	for cnt := len(statements); len(statements) > 0 && cnt > 0; cnt-- {
		for i := len(statements) - 1; i >= 0; i-- {
			s := statements[i]
			words := strings.Split(s, " ")

			// Факты могут быть следующих типов:
			// 1. «Aj is X years old» — возраст друга по имени Aj составляет X лет.
			// 2. «Aj is the same age as Ak» — другу по имени Aj столько же лет, сколько другу по имени Ak.
			// 3. «Aj is X years younger than Ak» — друг по имени Aj младше друга по имени Ak на X лет.
			// 4. «Aj is X years older than Ak» — друг по имени Aj старше друга по имени Ak на X лет.	return params{ /*todo*/ }

			name := words[0]

			if words[4] == "old" {
				age, err := strconv.Atoi(words[2])
				if err != nil {
					panic(err)
				}
				ages[name] = age
				if debugEnable {
					log.Println(name, ages[name])
				}
				removeStatement(i)
				continue
			}

			other := words[len(words)-1]

			if _, ok := ages[name]; !ok {
				if _, ok := ages[other]; !ok {
					continue
				}

				if words[3] == "same" {
					ages[name] = ages[other]
					if debugEnable {
						log.Println(name, ages[name])
					}
					removeStatement(i)
					continue
				}

				diff, err := strconv.Atoi(words[2])
				if err != nil {
					panic(err)
				}

				if words[4] == "younger" {
					ages[name] = ages[other] - diff
					if debugEnable {
						log.Println(name, ages[name])
					}
					removeStatement(i)
					continue
				}

				if words[4] == "older" {
					ages[name] = ages[other] + diff
					if debugEnable {
						log.Println(name, ages[name])
					}
					removeStatement(i)
					continue
				}
			} else {
				if _, ok := ages[other]; ok {
					removeStatement(i)
					continue
				}

				if words[3] == "same" {
					ages[other] = ages[name]
					if debugEnable {
						log.Println(other, ages[other])
					}
					removeStatement(i)
					continue
				}

				diff, err := strconv.Atoi(words[2])
				if err != nil {
					panic(err)
				}

				if words[4] == "younger" {
					ages[other] = ages[name] + diff
					if debugEnable {
						log.Println(other, ages[other])
					}
					removeStatement(i)
					continue
				}

				if words[4] == "older" {
					ages[other] = ages[name] - diff
					if debugEnable {
						log.Println(other, ages[other])
					}
					removeStatement(i)
					continue
				}
			}
		}
	}

	return results{ages[params.target]}
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
