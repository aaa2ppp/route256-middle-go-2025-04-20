package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type params struct {
	target     string
	statements []statement
}

type results struct {
	targetAge int
}

func parseQuery(s string) string {

	// Вопрос, который вы задали себе, выглядит следующим образом:
	// «How old is Ai?» — Cколько лет другу по имени Ai?

	words := splitToWords(s)
	return words[3]
}

type statementType int

const (
	_ statementType = iota
	isOld
	sameAs
	lessThan
)

type statement struct {
	type_ statementType
	name  string
	other string
	age   int
}

func parseStatement(s string) statement {

	// Факты могут быть следующих типов:
	// 1. «Aj is X years old» — возраст друга по имени Aj составляет X лет.
	// 2. «Aj is the same age as Ak» — другу по имени Aj столько же лет, сколько другу по имени Ak.
	// 3. «Aj is X years younger than Ak» — друг по имени Aj младше друга по имени Ak на X лет.
	// 4. «Aj is X years older than Ak» — друг по имени Aj старше друга по имени Ak на X лет.

	words := splitToWords(s)
	n := len(words)

	var st statement

	st.name = words[0]

	if words[3] == "years" {
		v, err := strconv.Atoi(words[2])
		if err != nil {
			panic(err)
		}
		st.age = v
	}

	switch {
	case words[4] == "old":
		st.type_ = isOld
	case words[3] == "same":
		st.type_ = sameAs
	case words[4] == "younger":
		st.type_ = lessThan
	case words[4] == "older":
		st.type_ = lessThan
		st.age = -st.age
	default:
		panic(fmt.Errorf("unknown type of statement: %v %v", s, words))
	}

	if st.type_ != isOld {
		st.other = words[n-1]
	}

	return st
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
	s, err := br.ReadString('\n')
	if err != nil {
		panic(err)
	}
	s = strings.TrimSpace(s)
	target := parseQuery(s)

	statements := make([]statement, 0, 3)
	for i := 0; i < 3; i++ {
		s, err := br.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s = strings.TrimSpace(s)
		statements = append(statements, parseStatement(s))
	}

	if debugEnable {
		log.Println("target:", target)
		log.Println("statements:", statements)
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
	statements := slices.Clone(params.statements)
	ages := make(map[string]int, 3)

	removeStatement := func(i int) {
		n := len(statements)
		statements[i] = statements[n-1]
		statements = statements[:n-1]
	}

mainLoop:
	for len(statements) > 0 {
		for i, st := range statements {

			if st.type_ == isOld {
				ages[st.name] = st.age
				removeStatement(i)
				continue mainLoop
			}

			nameAge, nameOk := ages[st.name]
			otherAge, otherOk := ages[st.other]

			if !nameOk && !otherOk {
				continue
			}

			switch st.type_ {
			case sameAs:
				if !nameOk {
					ages[st.name] = otherAge
				} else {
					ages[st.other] = nameAge
				}
			case lessThan:
				if !nameOk {
					ages[st.name] = otherAge - st.age
				} else {
					ages[st.other] = nameAge + st.age
				}
			default:
				panic(fmt.Errorf("unknown statemtType(%d)", st.type_))
			}

			removeStatement(i)
			continue mainLoop
		}

		if debugEnable {
			log.Printf("%+v", params.statements)
			log.Printf("%+v", statements)
		}
		panic("no applied statements")
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
