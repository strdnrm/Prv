package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type table struct {
	SortMethod func([]int) (int, int, time.Duration)
	Name       string
	Arr        []int
}

func main() {
	var size int
	fmt.Print("Введите размер массива: ")
	fmt.Scan(&size)
	// size = 10

	var min, max int
	fmt.Print("Введите диапазон случайных чисел (минимум и максимум): ")
	fmt.Scan(&min, &max)
	// min = 0
	// max = 100

	arr := generateRandomArray(size, min, max)

	fmt.Println("Исходный массив:", arr)

	tt := []table{
		{
			SortMethod: bubbleSort,
			Name:       "Пузырек",
			Arr:        make([]int, len(arr)),
		},
		{
			SortMethod: shellSort,
			Name:       "Шелл",
			Arr:        make([]int, len(arr)),
		},
		{
			SortMethod: quick,
			Name:       "Быстрая",
			Arr:        make([]int, len(arr)),
		},
	}

	for i := range tt {
		copy(tt[i].Arr, arr)
	}

	startTime := time.Now()
	for i := range tt {
		parComparisons, parSwaps, parTime := tt[i].SortMethod(tt[i].Arr)
		fmt.Printf("%s:\nКоличество сравнений: %d\nКоличество перестановок: %d\nЗатраченное время: %s\n", tt[i].Name, parComparisons, parSwaps, parTime.String())
		// fmt.Printf("Резуьтат: %s\n", fmt.Sprint(tt[i].Arr))
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Время последовательной работы всех сортировок: %s\n", elapsed)

	for i := range tt {
		tt[i].Arr = make([]int, len(arr))
		copy(tt[i].Arr, arr)
	}

	var wg sync.WaitGroup
	startAsync := time.Now()
	wg.Add(len(tt))
	for i := range tt {
		go func(i int) {
			defer wg.Done()
			parComparisons, parSwaps, parTime := tt[i].SortMethod(tt[i].Arr)
			fmt.Printf("%s - Поток %d:\nКоличество сравнений: %d\nКоличество перестановок: %d\nЗатраченное время: %s\n", tt[i].Name, i, parComparisons, parSwaps, parTime.String())
			// fmt.Printf("Резуьтат: %s\n", fmt.Sprint(tt[i].Arr))
		}(i)
	}
	wg.Wait()
	elapsed = time.Since(startAsync)
	fmt.Printf("Время паралелльной работы всех сортировок: %s\n", elapsed)
}

func generateRandomArray(size, min, max int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(max-min+1) + min
	}

	return arr
}

func bubbleSort(arr []int) (int, int, time.Duration) {
	comparisons := 0
	swaps := 0

	start := time.Now()

	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			comparisons++
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swaps++
			}
		}
	}

	elapsed := time.Since(start)

	return comparisons, swaps, elapsed
}

func shellSort(arr []int) (int, int, time.Duration) {
	comparisons := 0
	swaps := 0

	start := time.Now()

	n := len(arr)
	gap := n / 2
	for gap > 0 {
		for i := gap; i < n; i++ {
			temp := arr[i]
			j := i
			comparisons++
			for j >= gap && arr[j-gap] > temp {
				arr[j] = arr[j-gap]
				j -= gap
				swaps++
			}
			arr[j] = temp
		}
		gap /= 2
	}

	elapsed := time.Since(start)

	return comparisons, swaps, elapsed
}

func quick(arr []int) (int, int, time.Duration) {

	start := time.Now()
	comparisons, swaps := quickSort(arr, 0, len(arr)-1)
	elapsed := time.Since(start)
	return comparisons, swaps, elapsed
}

func quickSort(arr []int, low, high int) (int, int) {
	comparisons := 0
	swaps := 0

	if low < high {
		pivotIndex, cm, sw := partition(arr, low, high)
		comparisons += cm
		swaps += sw
		quickSort(arr, low, pivotIndex-1)
		quickSort(arr, pivotIndex+1, high)
	}

	return comparisons, swaps
}

func partition(arr []int, low, high int) (int, int, int) {
	comparisons := 0
	swaps := 0
	pivot := arr[high]
	i := low - 1

	for j := low; j <= high-1; j++ {
		comparisons++
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
			swaps++
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]

	return i + 1, comparisons, swaps
}
