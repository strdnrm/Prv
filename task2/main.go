package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	// runtime.GOMAXPROCS(16)
	var a, b, c, d int
	fmt.Println("Размер A:")
	_, err := fmt.Scan(&a, &b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Размер B:")
	_, err = fmt.Scan(&c, &d)
	if err != nil {
		log.Fatal(err)
	}
	matrixA := generateMatrix(a, b) //[][]int{{1, 2, 3}, {4, 5, 6}}
	matrixB := generateMatrix(c, d) //[][]int{{7, 8}, {9, 10}, {11, 12}}
	fmt.Println("Исходные матрицы:")
	fmt.Println(getMatrixStr(matrixA))
	fmt.Println(getMatrixStr(matrixB))

	multiplySequential(matrixA, matrixB)
	// sequentialResult := multiplySequential(matrixA, matrixB)
	// fmt.Println("Результат последовательного умножения:", sequentialResult)

	multiplyConcurrent(matrixA, matrixB)
	// concurrentResult := multiplyConcurrent(matrixA, matrixB)
	// fmt.Println("Результат асинхронного умножения:", concurrentResult)
}

func getMatrixStr(mat [][]int) string {
	res := ""
	for _, line := range mat {
		for _, elem := range line {
			res += fmt.Sprintf("%s ", strconv.Itoa(elem))
		}
		res += "\n"
	}
	return res
}

func generateMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			matrix[i][j] = rand.Intn(100)
		}
	}

	return matrix
}

func multiplySequential(matrixA, matrixB [][]int) [][]int {

	if len(matrixA[0]) != len(matrixB) {
		fmt.Println("Невозможно выполнить умножение матриц: количество столбцов матрицы A не равно количеству строк матрицы B")
		return nil
	}

	rowsA := len(matrixA)
	colsA := len(matrixA[0])
	colsB := len(matrixB[0])

	result := make([][]int, rowsA)
	start := time.Now()

	for i := 0; i < rowsA; i++ {
		result[i] = make([]int, colsB)
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				result[i][j] += matrixA[i][k] * matrixB[k][j]
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Println("Время выполнения последовательного умножения:", elapsed.String())

	return result
}

func multiplyConcurrent(matrixA, matrixB [][]int) [][]int {

	if len(matrixA[0]) != len(matrixB) {
		fmt.Println("Невозможно выполнить умножение матриц: количество столбцов матрицы A не равно количеству строк матрицы B")
		return nil
	}

	rowsA := len(matrixA)
	colsA := len(matrixA[0])
	colsB := len(matrixB[0])

	result := make([][]int, rowsA)

	var wg sync.WaitGroup
	wg.Add(rowsA * colsB)
	mu := new(sync.RWMutex)
	start := time.Now()

	for i := 0; i < rowsA; i++ {
		result[i] = make([]int, colsB)
		for j := 0; j < colsB; j++ {
			go func(row, col int) {
				defer wg.Done()

				sum := 0
				for k := 0; k < colsA; k++ {
					sum += matrixA[row][k] * matrixB[k][col]
				}
				mu.Lock()
				result[row][col] = sum
				mu.Unlock()
			}(i, j)
		}
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Время выполнения асинхронного умножения:", elapsed.String())

	return result
}
