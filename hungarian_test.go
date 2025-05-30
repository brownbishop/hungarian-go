package hungarian_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/brownbishop/hungarian-go"
)

const TEST_N = 12

func Test(t *testing.T) {
	for i := 0; i < TEST_N; i++ {
		testI(t, i)
	}
}

func testI(t *testing.T, i int) {
	costMatrix, err := ReadCostMatrix(fmt.Sprintf("test_data/test%d.in", i))
	if err != nil {
		t.Error(err)
	}
	correctCost, err := ReadCorrectCost(fmt.Sprintf("test_data/correct%d.out", i))
	if err != nil {
		t.Error(err)
	}

	assignment, err := hungarian.HunagarianAlgorithm(costMatrix)
	if err != nil {
		t.Error(err)
	}

	totalCost := 0
	for i := range assignment {
		totalCost += costMatrix[i][assignment[i]]
	}

	if totalCost != correctCost {
		t.Errorf("Test %d: wrong cost, expected %d, got %d", i, correctCost, totalCost)
	}
}

func ReadCostMatrix(fileName string) ([][]int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var m, n int
	_, err = fmt.Fscan(file, &m, &n)
	if err != nil {
		return nil, err
	}
	cost := make([][]int, m)
	for i := range m {
		cost[i] = make([]int, n)
		for j := range n {
			_, err := fmt.Fscan(file,  &cost[i][j])
			if err != nil {
				return nil, err
			}
		}
	}
	return cost, nil
}

func ReadCorrectCost(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}

	var cost int
	_, err = fmt.Fscan(file, &cost)
	if err != nil {
		return 0, err
	}

	file.Close()
	return cost, nil
}

