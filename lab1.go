package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//Заповнити квадратну матрицю випадковими числами. На побічній діагоналі розмістити мінімальний елемент стовпчика.
func measureTimeParalell(f func()) {
	start := time.Now()
	f()
	log.Println("execution took: ", time.Since(start))
}

func estimateTime(numOfThreads, size int) func() {
	start := time.Now()
	return func() {
		log.Println(size, numOfThreads, time.Since(start))
	}
}

func fillSlice(slice []int, maxVal int) {
	for i := range slice {
		(slice)[i] = 1
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

		go func() { // OKAY LET`S GO
			defer wg.Done()

			fillMatrix(matr[start:end], maxVal)
		}()
		fmt.Println(start, " - ", end)
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
func WriteTable() {
	matrixSizes := []int{
		5, 50, 250, 300, 350,
		400, 500, 600, 650, 700,
		750, 800, 900, 850, 900,
		1000, 1300, 1500, 2000, 2500}
	amountOfGoroutines := []int{
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
		11, 12, 13, 14, 15,
		16, 17, 18, 19, 20,
	}

	for _, size := range matrixSizes {
		Array2DParallel := NewMatrix(size)
		for _, numOfThreads := range amountOfGoroutines {
			measureTimeParalell(func() {
				fillMatrixParallel(Array2DParallel, 100, numOfThreads)
			})
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
		measureTimeParalell(func() {
			fillMatrixParallel(Array2DParallel, 100, numOfThreads)
		})
		time.Sleep(time.Second * 2)
		//print2DSlice(Array2DParallel)

		Array2D := NewMatrix(size)
		fmt.Println("sequential filling goes here: ")
		measureTimeParalell(func() {
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
	err := FillManually()
	if err != nil {
		log.Println(err)
		return
	}
	//plotting.MakePlot()
	//time.Sleep(time.Second * 1)
	//WriteTable()
}
