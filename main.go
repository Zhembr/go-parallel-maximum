// The program compares sequential and concurrent maximum searches.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

// generateRandomElements returns size positive random integers.
// A non-positive size yields an empty slice.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return []int{}
	}
	data := make([]int, size)
	for i := range data {
		// Intn excludes the upper bound, so adding one cannot overflow.
		data[i] = rand.Intn(int(^uint(0)>>1)) + 1
	}
	return data
}

// maximum returns the largest value, or zero when data is empty.
func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}
	maxValue := data[0]
	for _, value := range data[1:] {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

// maxChunks searches disjoint chunks in eight goroutines and combines results.
// Empty chunks are excluded from the final reduction when len(data) < CHUNKS.
func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}
	maxima := make([]int, CHUNKS)
	var workers sync.WaitGroup
	chunkSize, remainder := len(data)/CHUNKS, len(data)%CHUNKS
	start := 0
	for chunk := 0; chunk < CHUNKS; chunk++ {
		end := start + chunkSize
		if chunk < remainder {
			end++
		}
		workers.Add(1)
		go func(index int, values []int) {
			defer workers.Done()
			maxima[index] = maximum(values)
		}(chunk, data[start:end])
		start = end
	}
	workers.Wait()
	activeChunks := min(len(data), CHUNKS)
	return maximum(maxima[:activeChunks])
}

func main() {
	fmt.Printf("Генерируем %d целых положительных чисел\n", SIZE)
	data := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	started := time.Now()
	maxValue := maximum(data)
	elapsed := time.Since(started).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed)

	fmt.Printf("Ищем максимальное значение в %d горутинах\n", CHUNKS)
	started = time.Now()
	maxValue = maxChunks(data)
	elapsed = time.Since(started).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed)
}
