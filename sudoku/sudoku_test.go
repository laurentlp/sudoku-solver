package sudoku_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/laurentlp/sudoku-solver/sudoku"
)

const grid = "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"

//                     |
//                     v
const invalidGrid = "..757..3.1......2.7...234......8...4..7..4...49....6.5.42...3.....7..9....18....."
const errorGrid = "4.....8.5.3..........7......2.....6.....8.4......1...."
const emptyGrid = ""
const wrongGrid = "..757..3.1....a.2.7...234......8x..4..7..4...49....6.5.42...3e....7..9....18....."

func TestSolve(t *testing.T) {
	solved, err := sudoku.Solve(grid)

	if err != nil {
		t.Error("Error : ", err)
	}

	sudoku.Display(solved)
}

func TestSolveErr(t *testing.T) {
	_, err := sudoku.Solve(errorGrid)

	if err.Error() != "Invalid grid size: expected grid size of 81 found grid size of 54" {
		t.Error("Error : ", err)
	}
}

func TestSolveEmpty(t *testing.T) {
	_, err := sudoku.Solve(emptyGrid)

	if err.Error() != "Invalid grid size: expected grid size of 81 found grid size of 0" {
		t.Error("Error : ", err)
	}
}

func TestSolveInvalid(t *testing.T) {
	_, err := sudoku.Solve(invalidGrid)

	if err.Error() != "The sudoku contains errors and can not be solved" {
		t.Error("Error : ", err)
	}
}

func TestSolveWrong(t *testing.T) {
	_, err := sudoku.Solve(wrongGrid)

	if err.Error() != "The sudoku contains errors and can not be solved" {
		t.Error("Error : ", err)
	}
}

func TestSolveHardest(t *testing.T) {
	solveAll(fromFile("./_tests/hardest.txt"), "hardest", t)
}

func TestSolveTopSudoku(t *testing.T) {
	solveAll(fromFile("./_tests/top95.txt"), "hard", t)
}

func TestSolveEasy(t *testing.T) {
	solveAll(fromFile("./_tests/easy50.txt"), "easy", t)
}

// compareSlices compare values of two 3D arrays.
// Return true if they are the same
func compareSlices3D(A, B [][]string) bool {
	if len(A) != len(B) {
		return false
	}

	if (A == nil) != (B == nil) {
		return false
	}

	B = B[:len(A)] // Bounds-checking elimination
	for i, a := range A {
		for j, v := range a {
			if v != B[i][j] {
				return false
			}
		}
	}

	return true
}

// compareSlices compare values of two 2D arrays.
// Return true if they are the same
// In this case the boolean value of A is ignored
func compareSlices2D(A map[string]bool, B []string) bool {
	if len(A) != len(B) {
		return false
	}

	if (A == nil) != (B == nil) {
		return false
	}

	B = B[:len(A)] // Bounds-checking elimination
	i := 0
	for v, _ := range A {
		if v != B[i] {
			return false
		}
		i++
	}

	return true
}

// timeSolve calculate the time it takes to solve a sudoku
func timeSolve(grid string) (int64, bool) {
	nanosStart := time.Now().UnixNano()

	// Solve the sudoku in input
	_, err := sudoku.Solve(grid)

	duration := time.Now().UnixNano() - nanosStart

	return duration, err == nil
}

// fromFile load a sudoku from a file
func fromFile(filename string) []string {
	dat, _ := ioutil.ReadFile(filename)
	grids := strings.Split(string(dat), "\n")
	return grids[:len(grids)-1]
}

// nanoconv convert unix time to a unit in second
func nanoconv(nanos int64) float64 {
	return float64(nanos) / 1000000000.0
}

// solveAll the sudoku inside a file
func solveAll(grids []string, name string, t *testing.T) {
	times := make([]int64, len(grids))
	results := make([]bool, len(grids))

	for i, grid := range grids {
		t, result := timeSolve(grid)
		times[i] = t
		results[i] = result
	}

	n := len(grids)
	sudokuNumber := len(results)
	if n > 1 {
		// E.g. Solved 49 of 49 easy puzzles (avg 0.0033 secs (304.60 Hz), max 0.0112 secs).
		fmt.Printf("Solved %d of %d %s puzzles (avg %.4f secs (%.2f Hz), max %.4f secs).\n",
			sudokuNumber, // Number of sudoku solved
			n,            // Total number of sudoku to solve
			name,         // The type of the sudoku solved
			nanoconv(sum(times))/float64(n), // Average time to solve this type of sudoku
			float64(n)/nanoconv(sum(times)), // Average hertz used to solve this type of sudoku
			nanoconv(max(times)))            // The maximum time it took to solve one of the sudoku
	}

	if sudokuNumber != n {
		t.Error("Not all the sudoku have been solved")
	}
}

// max returns the maximum value of an array
func max(values []int64) (max int64) {
	max = 0

	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// sum returns the sum of all the number inside an array
func sum(items []int64) (tot int64) {
	tot = 0
	for _, i := range items {
		tot += i
	}
	return tot
}
