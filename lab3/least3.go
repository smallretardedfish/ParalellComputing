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
	First  int
	Second int
	Third  int
}

func GenerateSlice(scalar int) []int {
	rand.Seed(time.Now().Unix())
	res := make([]int, 100000*scalar)
	for i := range res {
		res[i] = rand.Intn(1000)
	}
	return res
}

func LeastThreeSerial(slice []int) Result {
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

func LeastThreeBlocking(slice []int, workerCount int) Result {
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
	worker := func(slice []int, jobs <-chan int, result *Result, mu *sync.RWMutex) {
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

func ThreeLeastAtomic(slice []int, workersCount int) Result {
	var v atomic.Value
	v.Store(&Result{
		First:  math.MaxInt,
		Second: math.MaxInt,
		Third:  math.MaxInt,
	})

	jobs := make(chan int, len(slice))
	go func() {
		for i := range slice {
			jobs <- i
		}
		close(jobs)
	}()

	worker := func(slice []int, jobs <-chan int, result atomic.Value) {
		for i := range jobs {
			res := v.Load()
			result, _ := res.(Result)
			if slice[i] < result.First {
				//atomic.CompareAndSwapPointer(&unsafe.Pointer(&result.Third), unsafe.Pointer(&result.Second))
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
	}
	var wg sync.WaitGroup
	wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go func() {
			worker(slice, jobs, v)
			defer wg.Done()
		}()
	}
	wg.Wait()
	r := v.Load().(Result)
	return r

}

func main() {
	arr := GenerateSlice(variant)
	fmt.Println(arr)
	min3 := LeastThreeBlocking(arr, 16)
	fmt.Println(min3)
}
