package main

import (
	"reflect"
	"sort"
	"testing"
)

/*21) Створити 2 масиви (або колекції) з випадковими числами.
У першому масиві - залишити елементи які більше 0.2 максимального значення масиву,
в другому залишити елементи кратні 10.
Відсортувати масиви і злити в один відсортований масив тільки ті елементи,
які входять в перший масив і не входять в другий.*/

var TestCasesMerge = []struct {
	first  []int
	second []int
	result []int
}{
	{first: []int{3, 6, 8, 0, 3, 8, 13},
		second: []int{5, 7, 3, 6, 8, 90},
		result: []int{0, 13},
	},
	{first: []int{57, 83, 113, 549, 80, 47},
		second: []int{5, 7, 3, 6, 8, 90},
		result: []int{47, 57, 80, 83, 113, 549},
	},
	{first: []int{5, 3, 2, 1},
		second: []int{5, 1, 3, 2, 7, 3, 6, 8, 90},
		result: []int{}, // empty slice literal, NOT nil-slice
	},
}
var TestCasesContains = []struct {
	arr    []int
	item   int
	result bool
}{
	{arr: []int{57, 83, 113, 549, 80, 47},
		item:   80,
		result: true,
	},
	{arr: []int{5, 1, 3, 2, 7, 3, 6, 8, 90},
		item:   80,
		result: false,
	},
}

func TestMergeIfInFirst(t *testing.T) {
	for _, tc := range TestCasesMerge {
		sort.Ints(tc.first)
		sort.Ints(tc.second)
		got := MergeIfPresentInFirst(tc.first, tc.second)
		if !reflect.DeepEqual(got, tc.result) {
			t.Errorf("MergeIfInFirst(%#v)\n got: %#v, want: %#v", tc.first, got, tc.result)
		}
	}
}

func TestSetInt_Contains(t *testing.T) {
	for _, tc := range TestCasesContains {
		set := SliceToSet(tc.arr)
		got := set.Contains(tc.item)
		if got != tc.result {
			t.Errorf("TestSetInt_Contains(%#v)\n got: %#v, want: %#v", tc.arr, got, tc.result)
		}
	}
}
