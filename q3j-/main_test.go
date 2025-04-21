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

func Test_run_fullset(t *testing.T) {
	test_run_fullset(t, solve)
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
5 7
3030500
0100001
0000000
0200050
4080009
4 5 7
`),
			`5
`,
			true,
		},
		{
			"2",
			strings.NewReader(`1
5 7
3230567
3330005
0000000
1230250
4120009
3 4 6
`),
			`3
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
