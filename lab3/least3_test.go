package main

import "testing"

var testCases = []struct {
	slice []int
	res   Result
}{
	{
		[]int{8, 23, 5, 3, 8, 56, -1, 7, 214},
		Result{
			First:  -1,
			Second: 3,
			Third:  5,
		},
	},
	{
		[]int{64, 7348, 742, 47236, 6, 6, 34},
		Result{
			First:  6,
			Second: 6,
			Third:  34,
		},
	},
	{
		[]int{-1, 0, 3, 31, 577, 534},
		Result{
			First:  -1,
			Second: 0,
			Third:  3,
		},
	},
}

func TestLeastThreeSerial(t *testing.T) {
	for _, tc := range testCases {
		got := LeastThreeSerial(tc.slice)
		want := tc.res
		if got != want {
			t.Errorf("LeastThreeSerial(%v) want:%+v, got:%+v", tc.slice, tc.res, got)
		}
	}
}

func TestLeastThreeBlocking(t *testing.T) {
	for _, tc := range testCases {
		got := LeastThreeBlocking(tc.slice, 3)
		want := tc.res
		if got != want {
			t.Errorf("LeastThreeBlocking(%v) want:%+v, got:%+v", tc.slice, tc.res, got)
		}
	}
}
