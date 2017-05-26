package sudoku

import (
	"errors"
	"strconv"
	"strings"
)

const digits string = "123456789"
const rows string = "ABCDEFGHI"
const cols string = digits

var squares = Cross(rows, cols)
var unitlist = CreateUnitList(rows, cols)
var units = CreateUnits(squares, unitlist)
var peers = CreatePeers(units)

// Cross the product of each elements in strings A and B together
func Cross(A, B string) []string {
	res := make([]string, len(A)*len(B))

	i := 0
	for _, a := range A {
		for _, b := range B {
			res[i] = string(a) + string(b)
			i++
		}
	}

	return res
}

// CreateUnitList list the 20 peers of each sudoku squares
func CreateUnitList(rows, cols string) [][]string {
	res := make([][]string, len(rows)*3)

	i := 0
	for _, c := range cols {
		// A1 B1 C1 D1 E1 F1 G1 H1 I1...
		res[i] = Cross(rows, string(c))
		i++
	}

	for _, r := range rows {
		res[i] = Cross(string(r), cols)
		i++
	}

	rs := []string{`ABC`, `DEF`, `GHI`}
	cs := []string{`123`, `456`, `789`}
	for _, r := range rs {
		for _, c := range cs {
			res[i] = Cross((r), string(c))
			i++
		}
	}

	return res
}

// CreateUnits find the units of each squares
func CreateUnits(squares []string, unitList [][]string) map[string][][]string {
	units := make(map[string][][]string, len(squares))

	for _, s := range squares {
		unit := make([][]string, 3)
		i := 0
		for _, u := range unitList {
			// For each squares of the unit
			for _, su := range u {
				if s == su {
					unit[i] = u
					i++
					break
				}
			}
		}
		units[s] = unit
	}

	return units
}

// CreatePeers find the 20 peers of a square
func CreatePeers(units map[string][][]string) map[string][]string {
	peers := make(map[string][]string, len(units))

	for s, ul := range units {
		peer := make([]string, 20)
		i := 0
		for _, u := range ul {
			for _, su := range u {
				if s != su {
					peer[i] = su
				}
			}
			i++
		}
		peers[s] = peer
	}

	return peers
}

// GridValues match all the sudoku values to its square
func GridValues(grid string) (map[string]string, error) {
	values := make(map[string]string, len(grid))
	chars := make([]string, len(grid))

	// For each square
	for i := 0; i < len(grid); i++ {
		// Value of the square
		str := grid[i : i+1]
		// Valid that the square value is a digit from 1 to 9 or '0' or '.' for empties
		// and adds it to the sudoku list of values.
		if strings.Contains(digits, str) || strings.Contains("0.", str) {
			chars[i] = str
		}
	}

	if len(chars) != 81 {
		return nil, errors.New("Invalid grid size: expected grid size of 81 found grid size of " + strconv.Itoa(len(chars)))
	}

	// Map the square value to it's corresponding key (A1, B5, D8,...)
	for i := 0; i < len(grid); i++ {
		values[squares[i]] = chars[i]
	}

	return values, nil
}
