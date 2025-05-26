package main

import "fmt"

type Mark int

const (
	None Mark = iota
	Star
	Prime
)

type Hungarian struct {
	matrix     [][]int
	rows, cols int
	marked     [][]Mark
	rowCovered []bool
	colCovered []bool
	rowPath    []int
	colPath    []int
}

func newHungarian(matrix [][]int) *Hungarian {
	h := Hungarian{
		matrix:     matrix,
		rows:       len(matrix),
		cols:       len(matrix[0]),
		marked:     make([][]Mark, len(matrix)),
		rowCovered: make([]bool, len(matrix)),
		colCovered: make([]bool, len(matrix[0])),
		z0row:      0,
		z0col:      0,
		rowPath:    []int{},
		colPath:    []int{},
	}

	for i := 0; i < h.rows; i++ {
		h.marked[i] = make([]Mark, h.cols)
	}

	return &h
}

// Substract the minimum element of a row from all elements in the row
func (h *Hungarian) Step1() {
	for i := 0; i < h.rows; i++ {
		min := h.matrix[i][0]
		for j := 1; j < h.cols; j++ {
			if h.matrix[i][j] < min {
				min = h.matrix[i][j]
			}
		}

		for j := 1; j < h.cols; j++ {
			h.matrix[i][j] -= min
		}
	}

	h.Step2()
}

// Repeat the row procedure for each column
func (h *Hungarian) Step2() {
	for j := 0; j < h.cols; j++ {
		min := h.matrix[0][j]
		for i := 1; i < h.rows; i++ {
			if h.matrix[i][j] < min {
				min = h.matrix[i][j]
			}
		}

		for i := 1; i < h.cols; i++ {
			h.matrix[i][j] -= min
		}
	}

	h.Step3()
}

func (h *Hungarian) clearCovers() {
	for i := 0; i < h.rows; i++ {
		h.rowCovered[i] = false
	}

	for i := 0; i < h.cols; i++ {
		h.colCovered[i] = false
	}
}

// All zeroes must be covered by marking as few rows or columns as possible.
// We assign tasks by starring a zero, only one starred zero can exist in
// each row or each column.
func (h *Hungarian) Step3() {
	for i := 0; i < h.rows; i++ {
		for j := 0; j < h.cols; j++ {
			if !h.rowCovered[i] && !h.colCovered[j] {
				h.marked[i][j] = Star
				h.rowCovered[i] = true
				h.colCovered[j] = true
			}
		}
	}
	h.Step4()
}

// Cover all columns containing a zero
func (h *Hungarian) Step4()  {
	h.clearCovers()
	for j := 0; j < h.cols; j++ {
		for i := 0; j < h.rows; i++ {
			if h.marked[i][j] == Star {
				h.colCovered[j] = true
				break
			}
		}
	}

find:
	zeroI, zeroJ := h.findNonCoveredZeroAndPrimeIt()
	if zeroI == -1 && zeroJ == -1  {
		h.Step5()
	} else {
		found := false
		for j := 0; j < h.cols; j++ {
			if h.marked[zeroI][j] == Star {
				h.colCovered[j] = false
				h.rowCovered[zeroI] = true
				found = true
			}
		}

		if !found {
			// todo: implement augmented path from wikipedia
		    fmt.Println(found)
		}
		goto find
	}
}

func (h *Hungarian) findNonCoveredZeroAndPrimeIt() (int, int) {
	for i := 0; i < h.rows; i++ {
		for j := 0; j < h.cols; j++ {
			if !h.rowCovered[i] && !h.colCovered[j] {
				h.marked[i][j] = Prime
				return i, j
			}
		}
	}
	return -1, -1
}

func (h *Hungarian) Step5() {
	count := 0
	for i := 0; i < h.rows; i++ {
		for j := 0; j < h.cols; j++ {
			if h.marked[i][j] == Star {
				count++
				fmt.Printf("%d %d\n", i, j)
			}
		}
	}

	if count == min(h.rows, h.cols) {
		fmt.Println("Assignment")
	}
}

func main() {

}
