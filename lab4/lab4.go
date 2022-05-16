package main

import (
	"fmt"
	"github.com/ParallelComputing/lab4/async"
	"math/rand"
	"sort"
	"time"
)

const variant = 21

/*21) Створити 2 масиви (або колекції) з випадковими числами.
У першому масиві - залишити елементи які більше 0.2 максимального значення масиву,
в другому залишити елементи кратні 10.
Відсортувати масиви і злити в один відсортований масив тільки ті елементи,
які входять в перший масив і не входять в другий.*/

type SetInt map[int]struct{}

func (si SetInt) Contains(item int) bool {
	_, ok := si[item]
	return ok
}

func RandInt(lower, upper int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng := upper - lower
	return r.Intn(rng) + lower
}

func ProcessSlice(slicePromise async.Promise[[]int], f func(async.Promise[[]int]) []int) []int {
	return f(slicePromise)
}

func Max(slice []int) int {
	max := slice[0]
	for i := range slice {
		if slice[i] > max {
			max = slice[i]
		}
	}
	return max
}

func LargerThen02Max(promise async.Promise[[]int]) []int {
	var res []int
	slice := promise.Await()
	max := Max(slice)

	for i := range slice {
		if float64(slice[i]) > 0.2*float64(max) {
			res = append(res, slice[i])
		}
	}
	return res
}

func DivisibleBy10(promise async.Promise[[]int]) []int {
	var res []int
	slice := promise.Await()

	for i := range slice {
		if slice[i]%10 == 0 {
			res = append(res, slice[i])
		}
	}
	return res
}

func GenerateSlice(scalar int) []int {
	res := make([]int, 100*scalar)
	for i := range res {
		res[i] = RandInt(-1000, 4912)
	}
	return res
}

func SliceToSet(slice []int) SetInt {
	res := make(SetInt)
	for i := range slice {
		res[slice[i]] = struct{}{}
	}
	return res
}

func MergeIfPresentInFirst(sl1, sl2 []int) []int {
	res := make([]int, 0)
	i, j := 0, 0

	set1 := SliceToSet(sl1)
	set2 := SliceToSet(sl2)

	for i < len(sl1) || j < len(sl2) {
		var itemToAdd int
		if j >= len(sl2) {
			itemToAdd = sl1[i]
			i++
		} else if i >= len(sl1) {
			itemToAdd = sl2[j]
			j++
		} else if sl1[i] <= sl2[j] {
			itemToAdd = sl1[i]
			i++
		} else if sl2[j] <= sl1[i] {
			itemToAdd = sl2[j]
			j++
		}
		if set1.Contains(itemToAdd) && !set2.Contains(itemToAdd) {
			res = append(res, itemToAdd)
		}
	}
	return res
}

func MergeIfInFirstPromise(slice1Pr, slice2Pr *async.Promise[[]int]) []int {
	var sl1 []int
	var sl2 []int

	done := make(chan bool)
	go func() {
		sl1 = slice1Pr.Await()
		sl1 = slice2Pr.Await()
		done <- true
	}()
	<-done

	return MergeIfPresentInFirst(sl1, sl2)
}

func main() {
	arr1Promise := async.DoAsync[[]int](
		func() []int {
			arr := GenerateSlice(variant)
			return arr
		})
	arr2Promise := async.DoAsync[[]int](
		func() []int {
			arr := GenerateSlice(variant)
			return arr
		})
	pr1 := async.DoAsync[[]int](
		func() []int {
			return ProcessSlice(*arr1Promise, LargerThen02Max)
		})
	pr2 := async.DoAsync[[]int](
		func() []int {
			return ProcessSlice(*arr2Promise, DivisibleBy10)
		})
	sortedPromise1 := async.DoAsync[[]int](
		func() []int {
			s1 := pr1.Await()
			sort.Ints(s1)
			return s1
		})
	sortedPromise2 := async.DoAsync[[]int](
		func() []int {
			s2 := pr2.Await()
			sort.Ints(s2)
			return s2
		})

	result := async.DoAsync[[]int](
		func() []int {
			return MergeIfInFirstPromise(sortedPromise1, sortedPromise2)
		})

	fmt.Println(result.Await())
}
