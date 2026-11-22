package main

import (
	"math/rand"
	"reflect"
	"slices"
	"strconv"
	"testing"
)

func TestGenerateRandomElements(t *testing.T) {
	for _, size := range []int{-10, 0, 1, 7, 8, 9, 1000} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			data := generateRandomElements(size)
			if len(data) != max(0, size) {
				t.Fatalf("length = %d for size %d", len(data), size)
			}
			for i, value := range data {
				if value <= 0 {
					t.Fatalf("data[%d] = %d, want a positive integer", i, value)
				}
			}
		})
	}
}

func TestMaximumAndMaxChunks(t *testing.T) {
	minInt := -int(^uint(0)>>1) - 1
	maxInt := int(^uint(0) >> 1)
	cases := []struct {
		name string
		data []int
		want int
	}{
		{"nil", nil, 0},
		{"empty", []int{}, 0},
		{"one", []int{42}, 42},
		{"zero", []int{0}, 0},
		{"negative singleton", []int{-8}, -8},
		{"fewer than chunks", []int{-9, -2, -5}, -2},
		{"equal", []int{5, 5, 5, 5, 5, 5, 5, 5}, 5},
		{"tail maximum", []int{1, 2, 3, 4, 5, 6, 7, 8, 99}, 99},
		{"head maximum", []int{99, 8, 7, 6, 5, 4, 3, 2, 1}, 99},
		{"negative uneven", []int{-9, -8, -7, -6, -5, -4, -3, -2, -1}, -1},
		{"integer extremes", []int{minInt, 0, maxInt}, maxInt},
		{"minimum integer", []int{minInt, minInt}, minInt},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before := slices.Clone(tc.data)
			for name, search := range map[string]func([]int) int{"maximum": maximum, "maxChunks": maxChunks} {
				if got := search(tc.data); got != tc.want {
					t.Errorf("%s() = %d, want %d", name, got, tc.want)
				}
				if !reflect.DeepEqual(tc.data, before) {
					t.Errorf("%s changed its input", name)
				}
			}
		})
	}
}

func TestMaxChunksCoversEveryPosition(t *testing.T) {
	for length := 1; length <= 3*CHUNKS+1; length++ {
		for position := 0; position < length; position++ {
			data := make([]int, length)
			data[position] = 100
			if got := maxChunks(data); got != 100 {
				t.Fatalf("length %d, position %d: got %d", length, position, got)
			}
		}
	}
}

func TestMaxChunksMatchesReference(t *testing.T) {
	random := rand.New(rand.NewSource(265847))
	for length := 1; length < 1000; length++ {
		data := make([]int, length)
		for i := range data {
			data[i] = random.Intn(20001) - 10000
		}
		if got, want := maxChunks(data), slices.Max(data); got != want {
			t.Fatalf("length %d: got %d, want %d", length, got, want)
		}
	}
}

func FuzzMaximumsAgree(f *testing.F) {
	f.Add([]byte{1, 9, 2, 3, 4, 5, 6, 7, 8})
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, raw []byte) {
		data := make([]int, len(raw))
		for i, value := range raw {
			data[i] = int(int8(value))
		}
		want := 0
		if len(data) > 0 {
			want = slices.Max(data)
		}
		if maximum(data) != want || maxChunks(data) != want {
			t.Fatalf("searches disagree with reference for %v", data)
		}
	})
}

func BenchmarkMaximum(b *testing.B) {
	data := generateRandomElements(10_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maximum(data)
	}
}

func BenchmarkMaxChunks(b *testing.B) {
	data := generateRandomElements(10_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maxChunks(data)
	}
}
