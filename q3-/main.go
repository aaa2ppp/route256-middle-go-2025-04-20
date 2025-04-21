package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type params struct {
	k    int
	nums []int
}

type results struct {
	ok bool
}

func readParams(br *bufio.Reader) params {
	var n, k int
	if _, err := fmt.Fscanln(br, &n, &k); err != nil {
		panic(err)
	}
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		var v int
		if _, err := fmt.Fscan(br, &v); err != nil {
			panic(err)
		}
		nums = append(nums, v)
	}
	_, err := br.ReadSlice('\n')
	if err != nil {
		panic(err)
	}
	return params{k, nums}
}

func writeResults(bw *bufio.Writer, results results) {
	if results.ok {
		bw.WriteString("YES\n")
	} else {
		bw.WriteString("NO\n")
	}
}

func slowSolve(params params) results {
	k := params.k
	nums := params.nums

	if debugEnable {
		log.Println("---")
	}

	for i := 0; i < len(nums)-k+1; i++ {
		if debugEnable {
			log.Printf("%2d\n", nums)
		}

		v := nums[i]
		nums[i] = 0

		// XXX slow
		for j := i + 1; j < i+k; j++ {
			if nums[j] < v {
				return results{false}
			}
			nums[j] -= v
		}
	}

	if debugEnable {
		log.Printf("%2d\n", nums)
	}

	for i := len(nums) - k; i < len(nums); i++ {
		if nums[i] != 0 {
			return results{false}
		}
	}

	return results{true}
}

func solve(params params) results {
	k := params.k
	nums := params.nums
	var diff []int

	if debugEnable {
		defer func() {
			log.Println("k:", k)
			log.Println("nums:", nums)
			log.Println("diff:", diff)
		}()
	}

	if k == 1 {
		return results{true}
	}

	diff = make([]int, len(nums))
	ded := 0

	for i := 0; i < len(nums); i++ {
		diff[i] = nums[i] - ded

		if diff[i] < 0 {
			return results{false}

		} else if i > len(diff)-k && diff[i] != 0 {
			return results{false}
		}

		ded += diff[i]
		if j := i - k + 1; j >= 0 {
			ded -= diff[j]
		}
	}

	return results{true}
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
