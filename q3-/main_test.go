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

const testDataDir = "./test_data/"

func Test_run(t *testing.T) {
	test_run(t, solve)
}

func Test_run_slow(t *testing.T) {
	test_run(t, slowSolve)
}

func Test_run_fullset(t *testing.T) {
	test_run_fullset(t, solve)
}

func Test_run_fullset_slow(t *testing.T) {
	test_run_fullset(t, slowSolve)
}

type runTest struct {
	name    string
	in      io.Reader
	wantOut string
	debug   bool
}

func test_run(t *testing.T, solve solveFunc) {
	tests := []runTest{
		{
			"1",
			strings.NewReader(`1
4 2
1 2 3 2
`),
			`YES
`,
			true,
		},
		{
			"2",
			strings.NewReader(`1
3 2
1 3 1
		`),
			`NO
`,
			true,
		},
		{
			"3",
			strings.NewReader(`1
4 2
9083 9870 8557 2302
`),
			`NO
`,
			true,
		},
		{
			"5",
			strings.NewReader(`7
4 2
6 16 20 10
9 6
5 7 7 7 7 7 1 0 0
12 9
2 3 2 1 3 3 3 1 2 2 2 1
10 6
1 4 4 4 9 9 8 5 5 5
3 1
1 12 3
8 6
2 4 9 9 9 9 7 5
53 37
6 8 8 9 12 12 12 12 15 15 15 18 19 19 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 22 16 14 14 13 10 10 10 10 7 7 7 4 3 2 0 0
`),
			`YES
NO
NO
YES
YES
YES
NO
`,
			true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.do(t, solve)
	}
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
