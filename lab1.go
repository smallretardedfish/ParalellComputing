package main

import (
	"fmt"
	"github.com/ParallelComputing/plotting"
	"log"
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
	//	size := len(matr)
	//	defer estimateTime(1, size)()
	for i := range matr {
		fillSlice(matr[i], maxVal)
	}
}
func FillMatrix(matr [][]int, maxVal int) {
	for i := range matr {
		fillSlice(matr[i], maxVal)
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
	//defer estimateTime(numOfThreads, size)()

	diff := numOfThreads - size%numOfThreads
	for i := 0; i < numOfThreads; i++ {
		if i >= diff {
			step = size/numOfThreads + 1
		} else {
			step = size / numOfThreads
		}
		end = end + step
		part := matr[start:end]
		go func(part [][]int) { // OKAY LET`S GO
			defer wg.Done()
			fillMatrix(part, maxVal)
		}(part)
		//fmt.Println(start, " - ", end)
		start = end
	}
	wg.Wait()
}

func PlaceMinOfColumnOnDiagonal(matr [][]int) {
	//estimateTime()()
	colCount := len(matr[0])
	rowCount := len(matr)
	for j := 0; j < colCount; j++ {
		var min = matr[0][j]
		for i := 0; i < rowCount; i++ {
			if matr[i][j] < min {
				min = matr[i][j]
			}
		}
		matr[rowCount-j-1][j] = min
	}
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
	matrixSizes := []int{
		5, 50, 250, 300, 350,
		400, 500, 600, 650, 700,
		750, 800, 900, 850, 900,
		1000, 1300, 1500, 2000, 2500}
	amountOfGoroutines := []int{
		1, 2, 3, 4, 5,
		11, 12, 100, 50, 15,
		200, 500, 750, 2000,
		5000,
	}

	for _, size := range matrixSizes {
		Array2DParallel := NewMatrix(size)
		for _, numOfThreads := range amountOfGoroutines {
			t := measureTimeParallel(func() {
				fillMatrixParallel(Array2DParallel, 100, numOfThreads)
			})
			s.Add(strconv.Itoa(numOfThreads), plotting.Stat{Size: size, T: t})
			plotting.CreatePlot(s)
			//fmt.Println(size)
		}
	}
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
		//fillMatrixParallel(Array2DParallel, 100, numOfThreads)
		measureTimeParallel(func() {
			fillMatrixParallel(Array2DParallel, 100, numOfThreads)
		})
		time.Sleep(time.Second * 2)
		//print2DSlice(Array2DParallel)

		Array2D := NewMatrix(size)
		fmt.Println("sequential filling goes here: ")
		measureTimeParallel(func() {
			fillMatrix(Array2D, 100)
		})
		//print2DSlice(Array2D)

		//fmt.Println("placing column`s minimums on the anti-diagonal goes here: ")
		//PlaceMinOfColumnOnDiagonal(Array2D)
		////print2DSlice(Array2D)
	}
	return nil
}

func main() {
	//err := FillManually()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//plotting.MakePlot()

	storage := plotting.NewStorage()
	WriteTable(storage)
	time.Sleep(time.Second * 3)
	//fmt.Println(storage)
}
