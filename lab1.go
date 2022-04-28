package main

import (
	"fmt"
	"github.com/ParallelComputing/plotting"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//Заповнити квадратну матрицю випадковими числами. На побічній діагоналі розмістити мінімальний елемент стовпчика.
func measureTimeParallel(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

func estimateTime(numOfThreads, size int) func() {
	start := time.Now()
	return func() {
		log.Println(size, numOfThreads, time.Since(start))
	}
}

func fillSlice(slice []int, maxVal int) {
	for i := range slice {
		slice[i] = 1
	}
}

func fillMatrix(matr [][]int, maxVal int) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range matr {
		for j := range matr {
			matr[i][j] = rng.Intn(maxVal)
		}
	}
}

func fillMatrixParallel(matr [][]int, maxVal int, numOfThreads int) {
	var (
		wg    sync.WaitGroup
		step  int
		start int
		end   int
	)
	wg.Add(numOfThreads)
	size := len(matr)

	diff := numOfThreads - size%numOfThreads
	for i := 0; i < numOfThreads; i++ {
		if i >= diff {
			step = size/numOfThreads + 1
		} else {
			step = size / numOfThreads
		}
		end = end + step
		part := matr[start:end]

		worker := func(part [][]int) {
			defer wg.Done()
			fillMatrix(part, maxVal)
		}
		go worker(part) // OKAY LET`S GO
		start = end
	}
	wg.Wait()
}

func MinOfSlice(arr []int) int {
	var min int
	for i := range arr {
		if arr[i] < arr[min] {
			min = i
		}
	}
	return arr[min]
}

func PlaceMinOfColumnOnDiagonal(matr [][]int, colIdx, step int) {
	rowCount := len(matr)
	for j := colIdx; j < colIdx+step; j++ {
		colItems := make([]int, rowCount)
		for i := 0; i < len(matr); i++ {
			colItems[i] = matr[i][j]
		}

		min := MinOfSlice(colItems)
		matr[rowCount-j-1][j] = min
	}
}

func PlaceMinOfColumnOnDiagonalParallel(matr [][]int, numOfThreads int) {
	var (
		wg sync.WaitGroup
	)
	wg.Add(numOfThreads)
	diff := numOfThreads - len(matr)%numOfThreads

	for i, colIdx := 0, 0; i < numOfThreads; i++ {
		var step int
		if i >= diff {
			step = len(matr)/numOfThreads + 1
		} else {
			step = len(matr) / numOfThreads
		}
		worker := func(matr [][]int, colIdx int, step int) {
			defer wg.Done()
			PlaceMinOfColumnOnDiagonal(matr, colIdx, step)
		}
		go worker(matr, colIdx, step)
		colIdx += step
	}
	wg.Wait()
}

func print2DSlice(matrix [][]int) {
	for _, row := range matrix {
		for _, el := range row {
			fmt.Printf("%d\t", el)
		}
		fmt.Println()
	}
}

func NewMatrix(size int) [][]int {
	array := make([][]int, size)
	buffer := make([]int, size*size)
	for i := range array {
		array[i], buffer = buffer[:size], buffer[size:]
	}
	return array
}

//WriteTable populates matrices
func WriteTable(s *plotting.Storage) {
	matrixSizes := []int{300, 500, 600, 650, 700, 1000, 2000, 4000, 5000, 7000, 10000}
	amountOfGoroutines := []int{2, 3, 10, 20, 50, 100, 1000, 10000}

	for _, size := range matrixSizes {
		Array2DParallel := NewMatrix(size)
		for _, numOfThreads := range amountOfGoroutines {
			t := measureTimeParallel(func() {
				fillMatrixParallel(Array2DParallel, 100, numOfThreads)
				PlaceMinOfColumnOnDiagonalParallel(Array2DParallel, numOfThreads)
			})
			s.Add(strconv.Itoa(numOfThreads), plotting.Stat{Size: size, T: t})
		}
	}
	plotting.CreatePlot(s)
}
func FillManually() error {
	for {
		fmt.Println("Enter size of square matrix:")
		var size int
		_, err := fmt.Scanf("%d\n", &size)
		if err != nil {
			log.Fatalln(err)
			return err
		}
		fmt.Println("Enter num of goroutines to fill the matrix: ")

		var numOfThreads int
		_, err = fmt.Scanf("%d\n", &numOfThreads)
		if err != nil {
			log.Fatalln(err)
			return err
		}
		Array2DParallel := NewMatrix(size)
		fmt.Println("parallel filling goes here: ")

		measureTimeParallel(func() {
			fillMatrixParallel(Array2DParallel, 100, numOfThreads)
		})

		Array2D := NewMatrix(size)
		fmt.Println("sequential filling goes here: ")
		measureTimeParallel(func() {
			fillMatrix(Array2D, 100)
		})
		//print2DSlice(Array2D)
		//fmt.Println("placing column`s minimums on the anti-diagonal goes here: ")
		//PlaceMinOfColumnOnDiagonal(Array2D)
	}
	return nil
}

func main() {
	//err := FillManually()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	fmt.Println(runtime.NumCPU())
	storage := plotting.NewStorage()
	WriteTable(storage)
}
