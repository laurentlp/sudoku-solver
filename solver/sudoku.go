package solver

import (
	"errors"
	"fmt"
	"strings"
)

const digits string = "123456789"
const rows string = "ABCDEFGHI"
const cols string = digits

var squares = cross(rows, cols)
var unitlist = createUnitList(rows, cols)
var units = createUnits(squares, unitlist)
var peers = createPeers(units)

// Cross the product of each elements in strings A and B together
func cross(A, B string) []string {
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
func createUnitList(rows, cols string) [][]string {
	res := make([][]string, len(rows)*3)

	i := 0
	for _, c := range cols {
		// A1 B1 C1 D1 E1 F1 G1 H1 I1...
		res[i] = cross(rows, string(c))
		i++
	}

	for _, r := range rows {
		res[i] = cross(string(r), cols)
		i++
	}

	rs := []string{`ABC`, `DEF`, `GHI`}
	cs := []string{`123`, `456`, `789`}
	for _, r := range rs {
		for _, c := range cs {
			res[i] = cross((r), string(c))
			i++
		}
	}

	return res
}

// CreateUnits find the units of each squares
func createUnits(squares []string, unitList [][]string) map[string][][]string {
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
func createPeers(units map[string][][]string) map[string]map[string]bool {
	peers := make(map[string]map[string]bool, len(units))

	for s, ul := range units {
		peer := make(map[string]bool, 20)
		for _, u := range ul {
			for _, su := range u {
				if s != su {
					peer[su] = true
				}
			}
		}
		peers[s] = peer
	}

	return peers
}

// GridValues match all the sudoku values to its square
func gridValues(grid string) (map[string]string, error) {
	values := make(map[string]string, len(grid))
	chars := make([]string, len(grid))

	// The number of clues given in the grid
	nbClues := 0
	var diffDigits []string

	// For each square
	for i := 0; i < len(grid); i++ {
		// Value of the square
		str := grid[i : i+1]

		// Valid that the square value is a digit from 1 to 9 ('0' or '.' for empties)
		// and add it to the sudoku list of values.
		if strings.Contains(digits, str) || strings.Contains("0.", str) {
			chars[i] = str
		}
		if strings.Contains(digits, str) {
			nbClues++
			if !contains(diffDigits, str) {
				diffDigits = append(diffDigits, str)
			}
		}
	}

	if len(chars) != 81 {

		return nil, fmt.Errorf("Invalid grid size: expected grid size of 81 found grid size of %d", len(chars))
	} else if nbClues < 17 {
		return nil, fmt.Errorf("Invalid number of squares filled: expected a minimum of 17 clues found %d", nbClues)
	} else if len(diffDigits) < 8 {
		return nil, fmt.Errorf("Invalid number of different clues digits: expected a minimum of 8 different digits found %d", len(diffDigits))
	}

	// Map the square value to it's corresponding key (A1, B5, D8,...)
	for i := 0; i < len(grid); i++ {
		values[squares[i]] = chars[i]
	}

	return values, nil
}

// ParseGrid convert a grid to a dict of possible values, {square: digits}, or
// return nil if a contradiction is detected.
func parseGrid(grid string) (values map[string]string, err error) {
	values = make(map[string]string, len(squares))
	for _, s := range squares {
		values[s] = digits
	}

	gr, err := gridValues(grid)
	for s, v := range gr {
		if strings.Contains(digits, v) {
			values = assign(values, s, v)
			if values == nil {
				return nil, nil
			}
		}
	}
	return values, err
}

// Eliminate removes d from values[s]; propagate when values or places <= 2.
// Return values, except return False if a contradiction is detected.
func eliminate(values map[string]string, s string, v string) (map[string]string, bool) {
	// The value is already eliminated
	if !strings.Contains(values[s], v) {
		return values, true
	}

	// Remove all occurrences of the value (v) from the square possible values
	values[s] = strings.Replace(values[s], v, "", -1)

	// If a square (s) is reduced to one value (v2), then eliminate the value from the peers.
	if len(values[s]) == 0 {
		return nil, false
	} else if len(values[s]) == 1 {
		v2 := values[s]

		for s2 := range peers[s] {
			if _, ok := eliminate(values, s2, v2); !ok {
				return nil, false
			}
		}
	}

	// If a unit (u) has only one possible place for a value (v), then put it there.
	for _, u := range units[s] {
		dplaces := []string{}
		for _, s := range u {
			if strings.Contains(values[s], v) {
				dplaces = append(dplaces, s)
			}
		}

		if len(dplaces) == 0 {
			return nil, false
		} else if len(dplaces) == 1 {
			if assign(values, dplaces[0], v) == nil {
				return nil, false
			}
		}
	}

	return values, true
}

// Assign eliminate all the other values (except v) from a square possible values and propagate.
func assign(values map[string]string, s string, v string) map[string]string {
	otherValues := strings.Replace(values[s], v, "", -1)
	for _, v := range otherValues {
		if _, ok := eliminate(values, s, string(v)); !ok {
			return nil
		}
	}
	return values
}

// Search using depth-first search and propagation, try all possible values.
func search(values map[string]string) (map[string]string, error) {
	if values == nil {
		return nil, errors.New("The sudoku contains errors and can not be solved")
	}

	// Check if there is only one remaining possibility in every square
	// If true, return the solved sudoku
	solved := true
	for s := range values {
		if len(values[s]) != 1 {
			solved = false
		}
	}

	if solved {
		return values, nil
	}

	// Chose the first unfilled square with the fewest possibilities
	min := len(digits) + 1
	sq := ""
	for _, s := range squares {
		l := len(values[s])
		if l > 1 && l < min {
			sq = s
			min = l
		}
	}

	ch := make(chan map[string]string)
	for _, v := range values[sq] {
		go func(val string) {
			newValues := cloneValues(values)
			value, _ := search(assign(newValues, sq, val))
			if value != nil {
				ch <- value
			}
		}(string(v))
	}
	return <-ch, nil
}

// CloneValues from one map to another
func cloneValues(m map[string]string) map[string]string {
	newMap := make(map[string]string, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// Solve the sudoku in input
func Solve(grid string) (map[string]string, error) {
	pg, err := parseGrid(grid)
	if err != nil {
		return nil, err
	}

	return search(pg)
}

// Display the solved sudoku
func Display(values map[string]string) {
	for i, row := range rows {
		for j, col := range digits {
			if j == 3 || j == 6 {
				fmt.Printf("| ")
			}
			fmt.Printf("%v ", values[string(row)+string(col)])
		}
		fmt.Println()
		if i == 2 || i == 5 {
			fmt.Println("------+-------+-------")
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
