package sudoku_test

import (
	"fmt"
	"testing"

	"github.com/laurentlp/sudoku-solver/sudoku"
)

const digits string = "123456789"
const rows string = "ABCDEFGHI"
const cols string = digits

var squares []string
var unitList [][]string
var squareUnits map[string][][]string
var squarePeers map[string][]string

func TestCross(t *testing.T) {
	squares = sudoku.Cross(rows, cols)

	if len(squares) != 81 {
		t.Error("Expected 81 squares")
	}
}

func TestUnitList(t *testing.T) {
	unitList = sudoku.CreateUnitList(rows, cols)

	var tot int
	for _, u := range unitList {
		tot += len(u)
	}

	if tot != 81*3 {
		t.Error("Expected 243 squares got : ", tot)
	}
}

func TestCreateUnits(t *testing.T) {
	squareUnits = sudoku.CreateUnits(squares, unitList)
	for _, u := range squareUnits {
		if len(u) != 3 {
			t.Error("Unit length expected 3 got : ", len(u))
		}
	}
}

func TestUnit(t *testing.T) {
	unit := [][]string{
		[]string{"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2", "I2"},
		[]string{"C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9"},
		[]string{"A1", "A2", "A3", "B1", "B2", "B3", "C1", "C2", "C3"},
	}

	c2 := squareUnits["C2"]
	if !compareSlices3D(c2, unit) {
		t.Error("An error occured while creating units. Some returned values are wrong !")
	}
}

func TestCreatePeers(t *testing.T) {
	squarePeers = sudoku.CreatePeers(squareUnits)
	for k, p := range squarePeers {
		if len(p) != 20 {
			t.Error("Peers length expected 20 got : ", len(squarePeers[k]), " for the square ", k)
		}
	}
}

func TestPeers(t *testing.T) {
	unit := []string{
		"A2", "B2", "D2", "E2", "F2", "G2", "H2", "I2",
		"C1", "C3", "C4", "C5", "C6", "C7", "C8", "C9",
		"A1", "A3", "B1", "B3",
	}

	c2 := squarePeers["C2"]
	if compareSlices2D(c2, unit) {
		t.Error("An error occured while creating peers. Some returned values are wrong !")
	}
}

func TestGrid(t *testing.T) {
	grid := "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......" // 81 values
	result := map[string]string{}

	for i, s := range squares {
		result[s] = string(grid[i])
	}

	values, err := sudoku.GridValues(grid)

	if len(values) != len(grid) {
		t.Error(err)
	}

	for k, v := range values {
		if v != result[k] {
			t.Error("Value is invalid. Expected a value of " + result[k] + " got one of " + v)
		}
	}
}

func TestGridErr(t *testing.T) {
	grid := "4.....8.5.3..........7......2.....6.....8.4......1...." // 54 values

	_, err := sudoku.GridValues(grid)

	if err.Error() != ("Invalid grid size: expected grid size of 81 found grid size of " + fmt.Sprintf("%d", len(grid))) {
		t.Error("An error was supposed to be returned")
	}
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
func compareSlices2D(A, B []string) bool {
	if len(A) != len(B) {
		return false
	}

	if (A == nil) != (B == nil) {
		return false
	}

	B = B[:len(A)] // Bounds-checking elimination
	for i, a := range A {
		if a != B[i] {
			return false
		}
	}

	return true
}
