package main

import (
	"fmt"
	"github.com/ParallelComputing/lab4/async"
	"math/rand"
	"time"
)

const variant = 21

/*21) Створити 2 масиви (або колекції) з випадковими числами.
У першому масиві - залишити елементи які більше 0.2 максимального значення масиву,
в другому залишити елементи кратні 10.
Відсортувати масиви і злити в один відсортований масив тільки ті елементи,
які входять в перший масив і не входять в другий.*/

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

func main() {
	//arr1Promise := async.DoAsync[[]int](
	//	func() []int {
	//		arr := GenerateSlice(variant)
	//		return arr
	//	})
	arr2Promise := async.DoAsync[[]int](
		func() []int {
			arr := GenerateSlice(variant)
			return arr
		})
	res := ProcessSlice(*arr2Promise, DivisibleBy10)
	fmt.Println(res)
}
