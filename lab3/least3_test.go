package main

import "testing"

var testArray = GenerateSlice(variant)

var testCases = []struct {
	slice []int64
	res   Result
}{
	{
		[]int64{8, 23, 5, 3, 8, 56, -1, 7, 214},
		Result{
			First:  -1,
			Second: 3,
			Third:  5,
		},
	},
	{
		[]int64{64, 7348, 742, 47236, 6, 6, 34},
		Result{
			First:  6,
			Second: 6,
			Third:  34,
		},
	},
	{
		[]int64{-1, 0, 3, 31, 577, 534},
		Result{
			First:  -1,
			Second: 0,
			Third:  3,
		},
	},
}

func RunTestAnyThreeLeastFunc(t *testing.T, numOfThreads int, f func([]int64, int) Result) {
	for _, tc := range testCases {
		got := f(tc.slice, numOfThreads)
		want := tc.res
		if got != want {
			t.Errorf("LeastThreeSerial(%v) want:%+v, got:%+v", tc.slice, tc.res, got)
		}
	}
}

func TestLeastThreeSerial(t *testing.T) {
	RunTestAnyThreeLeastFunc(t, 1, LeastThreeBlocking)
}

func TestLeastThreeBlocking(t *testing.T) {
	RunTestAnyThreeLeastFunc(t, 16, LeastThreeBlocking)
}

func TestLeastThreeAtomic(t *testing.T) {
	RunTestAnyThreeLeastFunc(t, 16, LeastThreeAtomic)
}

func BenchmarkLeastThreeAtomic(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		LeastThreeAtomic(testArray, 250)
	}
	b.StopTimer()
}

func BenchmarkLeastThreeBlocking(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		LeastThreeBlocking(testArray, 250)
	}
	b.StopTimer()
}

//
