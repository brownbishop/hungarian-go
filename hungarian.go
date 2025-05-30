package hungarian

import (
	"fmt"
	"math"
)

type Mark int

const (
	None Mark = iota
	Star
	Prime
)

type Zero struct {
	r, c int
}

type hungarian struct {
	matrix     [][]int
	n          int
	marked     [][]Mark
	rowCovered []bool
	colCovered []bool
}

func printMatrix[T any](matrix [][]T) {
	for _, v := range matrix {
		for _, i := range v {
			fmt.Printf("%3v ", i)
		}
		fmt.Println("")
	}
}

func newHungarian(matrix [][]int) *hungarian {
	m := len(matrix)
	n := len(matrix[0])
	N := max(m, n)
	h := hungarian{
		matrix:     make([][]int, N),
		n:          n,
		marked:     make([][]Mark, N),
		rowCovered: make([]bool, N),
		colCovered: make([]bool, N),
	}

	for i := 0; i < N; i++ {
		h.marked[i] = make([]Mark, N)
		h.matrix[i] = make([]int, N)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			h.matrix[i][j] = matrix[i][j]
		}
	}

	return &h
}

func (h *hungarian) Compute() {
	steps := map[int]func() int{
		1: h.step1,
		2: h.step2,
		3: h.step3,
		4: h.step4,
		5: h.step5,
	}

	idx := 1
	for true {
		f := steps[idx]
		if f == nil {
			break
		}
		idx = f()
	}
}

func (h *hungarian) clearCovers() {
	for i := 0; i < h.n; i++ {
		h.rowCovered[i] = false
		h.colCovered[i] = false
	}
}

func (h *hungarian) step1() int {
	for i := 0; i < h.n; i++ {
		min := math.MaxInt
		for j := 0; j < h.n; j++ {
			if h.matrix[i][j] < min {
				min = h.matrix[i][j]
			}
		}

		for j := 0; j < h.n; j++ {
			h.matrix[i][j] -= min
		}
	}

	for j := 0; j < h.n; j++ {
		min := math.MaxInt
		for i := 0; i < h.n; i++ {
			if h.matrix[i][j] < min {
				min = h.matrix[i][j]
			}
		}

		for i := 0; i < h.n; i++ {
			h.matrix[i][j] -= min
		}
	}
	return 2
}

func (h *hungarian) step2() int {
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.n; j++ {
			if (h.matrix[i][j] == 0) && (!h.rowCovered[i]) && (!h.colCovered[j]) {
				h.marked[i][j] = Star
				h.rowCovered[i] = true
				h.colCovered[j] = true
				break
			}
		}
	}

	h.clearCovers()
	return 3
}

func (h *hungarian) step3() int {
	count := 0
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.n; j++ {
			if (h.marked[i][j] == Star) && !h.colCovered[j] {
				count++
				h.colCovered[j] = true
			}
		}
	}

	if count >= h.n {
		return 6
	}

	return 4
}

func (h *hungarian) step4() int {
	for true {
		row, col := h.findAZero()
		if row < 0 {
			return 5
		}
		h.marked[row][col] = Prime
		starCol := h.findStarInRow(row)
		if starCol >= 0 {
			h.rowCovered[row] = true
			h.colCovered[starCol] = false
		} else {
			path := h.makePath(row, col)
			for _, zero := range path {
				switch h.marked[zero.r][zero.c] {
				case Prime:
					h.marked[zero.r][zero.c] = Star
				case Star:
					h.marked[zero.r][zero.c] = None
				}
			}
			h.clearPrimes()
			h.clearCovers()
			return 3
		}
	}
	return 5
}

func (h *hungarian) step5() int {
	minVal := h.findSmallest()
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.n; j++ {
			if !h.rowCovered[i] && !h.colCovered[j] {
				h.matrix[i][j] -= minVal
			}

			if h.rowCovered[i] && h.colCovered[j] {
				h.matrix[i][j] += minVal
			}
		}
	}

	h.clearCovers()
	return 3
}

func (h *hungarian) clearPrimes() {
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.n; j++ {
			if h.marked[i][j] == Prime {
				h.marked[i][j] = None
			}
		}
	}
}
func (h *hungarian) makePath(primeRow, primeCol int) []Zero {
	path := []Zero{{primeRow, primeCol}}
	prime := Zero{primeRow, primeCol}
	var star Zero
	for true {
		row :=  h.findStarInColumn(prime.c)
		if row < 0 {
			break
		}
		star.r = row
		star.c = prime.c
		path = append(path, star)

		col := h.findPrimeInRow(star.r)
		prime.r = star.r
		prime.c = col
		path = append(path, prime)
	}
	return path
}

func (h *hungarian) findAZero() (int, int) {
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.n; j++ {
			if (h.matrix[i][j] == 0) && !h.rowCovered[i] && !h.colCovered[j] {
				return i, j
			}
		}
	}
	return -1, -1
}

func (h *hungarian) findPrimeInRow(row int) int {
	for j := 0; j < h.n; j++ {
		if h.marked[row][j] == Prime {
			return j
		}
	}
	return -1
}

func (h *hungarian) findStarInRow(row int) int {
	for j := 0; j < h.n; j++ {
		if h.marked[row][j] == Star {
			return j
		}
	}
	return -1
}

func (h *hungarian) findStarInColumn(col int) int {
	for i := 0; i < h.n; i++ {
		if h.marked[i][col] == Star {
			return i
		}
	}
	return -1
}

func (h *hungarian) erasePrimes() {
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.n; j++ {
			h.marked[i][j] = None
		}
	}
}

func (h *hungarian) findSmallest() int {
	minVal := math.MaxInt
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.n; j++ {
			if !h.rowCovered[i] && !h.colCovered[j] {
				if minVal > h.matrix[i][j] {
					minVal = h.matrix[i][j]
				}
			}
		}
	}
	return minVal
}

func HunagarianAlgorithm(matrix [][]int) ([]int, error) {
	m := len(matrix)
	n := len(matrix[0])
	if m > n {
		err := fmt.Errorf("This mode is unsupported, please transpose frist")
		return nil, err
	}

	h := newHungarian(matrix)
	h.Compute()

	assignment := make([]int, m)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if h.marked[i][j] == Star {
				assignment[i] = j
				break
			}
		}
	}

	return assignment, nil
}
