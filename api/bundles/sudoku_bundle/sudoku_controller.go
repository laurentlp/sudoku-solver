package sudokubundle

import (
	"net/http"
	"sort"

	"github.com/laurentlp/sudoku-solver/api/common"
	"github.com/laurentlp/sudoku-solver/api/errors"
	"github.com/laurentlp/sudoku-solver/solver"
)

// SudokuController struct
type SudokuController struct {
	common.Controller
}

// Solve a sudoku and return the output.
// If there is an errors of any sort (bad sudoku, no input,...) it is sent to the client
func (s *SudokuController) Solve(w http.ResponseWriter, r *http.Request) {

	var model Sudoku
	err := s.MapJSON(w, r, &model)
	if err == nil {

		res, err := solver.Solve(model.Sudoku)

		if err != nil {
			s.SendJSON(w, r, errors.BadRequest(err.Error()), http.StatusBadRequest)
			return
		}

		solved := err == nil

		solvedSudoku := NewSudoku(toString(res), solved)
		s.SendJSON(w, r, solvedSudoku, http.StatusOK)
		return
	}
	s.SendJSON(w, r, err, err.Status)
}

// ToString convert the solved sudoku (map[string]string) to as string of values
func toString(solvedSudoku map[string]string) (res string) {
	keys := []string{}
	for k := range solvedSudoku {
		keys = append(keys, k)
	}

	// Sort the map keys to make sure the output is ordered
	sort.Strings(keys)
	for _, k := range keys {
		res += solvedSudoku[k]
	}

	return res
}
