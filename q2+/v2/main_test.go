package main

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

const testDataDir = "../test_data/"

func Test_run(t *testing.T) {
	test_run(t, solve)
}

func Test_run_fullset(t *testing.T) {
	test_run_fullset(t, solve)
}

type runTest struct {
	name    string
	in      io.Reader
	wantOut string
	debug   bool
}

func (tt runTest) do(t *testing.T, solve solveFunc) {
	t.Run(tt.name, func(t *testing.T) {
		defer func(v bool) { debugEnable = true }(debugEnable)
		debugEnable = tt.debug

		w := &strings.Builder{}
		run(tt.in, w, solve)
		if gotW := w.String(); gotW != tt.wantOut {
			t.Errorf("run() = %v, want %v", gotW, tt.wantOut)
		}
	})
}

func test_run(t *testing.T, solve solveFunc) {
	tests := []runTest{
		{
			"1",
			strings.NewReader(`3
How old is A?
B is 2 years younger than Caac
A is 44 years older than B
B is 55 years old
How old is B?
A is 2 years older than Ca
B is 10 years younger than Ca
A is 2 years old
How old is Bd?
C is 23 years younger than Bd
C is 38 years younger than A
Bd is 27 years old
`),
			`99
-10
27
`,
			true,
		},
		{
			"2",
			strings.NewReader(`1
How old is I?
I is 2 years old
You is 3 years old
We is 4 years old
`),
			`2
`,
			true,
		},
		{
			"3",
			strings.NewReader(`1
How old is I?
We is 40 years old
You is the same age as We
I is the same age as You
`),
			`40
`,
			true,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt.do(t, solve)
	}
}

func test_run_fullset(t *testing.T, solve solveFunc) {
	files, err := os.ReadDir(testDataDir)
	if err != nil {
		panic(err)
	}

	var testNums []int
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".a") {
			continue
		}

		testName := strings.TrimSuffix(fileName, ".a")
		testNum, err := strconv.Atoi(testName)
		if err != nil {
			t.Log(err)
			continue
		}

		testNums = append(testNums, testNum)
	}

	if len(testNums) == 0 {
		t.Log("no any test")
		return
	}

	slices.Sort(testNums)

	for _, testNum := range testNums {
		func() {
			testName := strconv.Itoa(testNum)
			testPath := filepath.Join(testDataDir, testName)

			testFile, err := os.Open(testPath)
			if err != nil {
				t.Fatal(err)
			}
			defer testFile.Close()

			wantOut, err := os.ReadFile(testPath + ".a")
			if err != nil {
				t.Fatal(err)
			}

			runTest{
				name:    testName,
				in:      testFile,
				wantOut: unsafe.String(unsafe.SliceData(wantOut), len(wantOut)),
				debug:   false,
			}.do(t, solve)
		}()
	}
}
