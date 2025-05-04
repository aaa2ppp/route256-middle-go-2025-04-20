package main

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

const testDataDir = "./test_data"

func Test_run(t *testing.T) {
	test_run(t, solve)
}

func Test_run_set1(t *testing.T) {
	test_run_fullset(t, solve, testDataDir)
}

func Test_run_set2(t *testing.T) {
	test_run_fullset(t, solve, testDataDir+"2")
}
func Test_run_set3(t *testing.T) {
	test_run_fullset(t, solve, testDataDir+"3")
}

func Test_run_fullset(t *testing.T) {
	test_run_fullset(t, solve, testDataDir)
	test_run_fullset(t, solve, testDataDir+"2")
	test_run_fullset(t, solve, testDataDir+"3")
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
5 3 2
..#
##.
.##
##.
.##
5 16
10 2
`),
			`###
.##
##.
.##

##
##
.#
##

`,
			true,
		},
		{
			"2",
			strings.NewReader(`1
8 9 3
#########
#########
#########
#########
#########
#########
#########
#########
21 7
14 28
23 7
`),
			`######
######
######
######
######
######
######
######

#\.....
##\....
###\...
####\..
#####\.
######\
.######

#\...
##\..
###\.
####>
###/.
##/..
./...

`,
			true,
		},
		{
			"3",
			strings.NewReader(`2
3 5 1
#####
#####
#####
5 10
1 1 1
#
4 3
`),
			`####
####
####

#

`,
			true,
		},
		{
			"8x",
			strings.NewReader(`1
1 9 7
#....##.#
16 5
2 9
8 5
5 8
4 1
1 8
4 5
`),
			`####

###

###

###

###

###

###

`,
			true,
		},
		// TODO: Add test cases.
	}

	if len(tests) == 0 {
		t.Log("no any test")
		return
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

func test_run_fullset(t *testing.T, solve solveFunc, testDataDir string) {
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

func Test_getQPoint(t *testing.T) {
	type args struct {
		n int
		m int
		q int
	}
	tests := []struct {
		name   string
		args   args
		want   point
		wantOk bool
	}{
		{
			"1",
			args{1, 1, 0},
			point{0, 0},
			true,
		},
		{
			"2",
			args{1, 1, 4},
			point{0, 0},
			false,
		},
		{
			"3",
			args{2, 3, 3},
			point{0, 3},
			true,
		},
		{
			"3",
			args{2, 3, 5},
			point{2, 3},
			true,
		},
		{
			"3",
			args{2, 3, 9},
			point{1, 0},
			true,
		},
		{
			"3",
			args{3, 5, 15},
			point{1, 0},
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := getAxisPoint(tt.args.n, tt.args.m, tt.args.q)
			if ok != tt.wantOk {
				t.Errorf("getQPoint() ok = %v, want %v", ok, tt.wantOk)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getQPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
