package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const variant = 21

type Result struct {
	First  int64
	Second int64
	Third  int64
}

func RandInt(lower, upper int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng := upper - lower
	return int64(r.Intn(int(rng))) + lower
}

func GenerateSlice(scalar int) []int64 {
	res := make([]int64, 100000*scalar)
	for i := range res {
		res[i] = RandInt(-1000000000000, 49123674982)
	}
	return res
}

func LeastThreeSerial(slice []int64) Result {
	result := Result{
		First:  math.MaxInt,
		Second: math.MaxInt,
		Third:  math.MaxInt,
	}
	for i := range slice {

		if slice[i] < result.First {
			result.Third = result.Second
			result.Second = result.First
			result.First = slice[i]
		} else if slice[i] < result.Second {
			result.Third = result.Second
			result.Second = slice[i]
		} else if slice[i] < result.Third {
			result.Third = slice[i]
		}
	}
	return result
}

func LeastThreeBlocking(slice []int64, workerCount int) Result {
	var mu sync.RWMutex
	result := Result{
		First:  math.MaxInt,
		Second: math.MaxInt,
		Third:  math.MaxInt,
	}
	jobs := make(chan int, len(slice))
	go func() {
		for i := range slice {
			jobs <- i
		}
		close(jobs)
	}()
	worker := func(slice []int64, jobs <-chan int, result *Result, mu *sync.RWMutex) {
		for i := range jobs {
			mu.Lock()
			if slice[i] < result.First {
				result.Third = result.Second
				result.Second = result.First
				result.First = slice[i]
			} else if slice[i] < result.Second {
				result.Third = result.Second
				result.Second = slice[i]
			} else if slice[i] < result.Third {
				result.Third = slice[i]
			}
			mu.Unlock()
		}
	}
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			worker(slice, jobs, &result, &mu)
			defer wg.Done()
		}()
	}
	wg.Wait()
	return result
}

func LeastThreeAtomic(slice []int64, workersCount int) Result {
	r := &Result{
		First:  math.MaxInt,
		Second: math.MaxInt,
		Third:  math.MaxInt,
	}

	jobs := make(chan int, len(slice))
	go func() {
		for i := range slice {
			jobs <- i
		}
		close(jobs)
	}()

	worker := func(slice []int64, jobs <-chan int, result *Result) {
		for i := range jobs {
			if slice[i] < result.First {
				atomic.CompareAndSwapInt64(&result.Third, result.Third, result.Second)  //result.Third = result.Second ??
				atomic.CompareAndSwapInt64(&result.Second, result.Second, result.First) //result.Second = result.First
				atomic.CompareAndSwapInt64(&result.First, result.First, slice[i])       //result.First = slice[i]
			} else if slice[i] < result.Second {
				atomic.CompareAndSwapInt64(&result.Third, result.Third, result.Second) //result.Third = result.Second
				atomic.CompareAndSwapInt64(&result.Second, result.Second, slice[i])    //result.Second = slice[i]
			} else if slice[i] < result.Third {
				atomic.CompareAndSwapInt64(&result.Third, result.Third, slice[i]) //result.Third = slice[i]
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go func() {
			worker(slice, jobs, r)
			defer wg.Done()
		}()
	}
	wg.Wait()

	return *r
}

func main() {
	arr := GenerateSlice(variant)
	fmt.Println(arr)
	//toTest := []int64{
	//	1, 3, 5, -100, 534, -65, 0,
	//}
	min3 := LeastThreeBlocking(arr, 128)
	fmt.Println(min3)
}
