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

// Solve func return a solved sudoku
func (s *SudokuController) Solve(w http.ResponseWriter, r *http.Request) {

	var model Sudoku
	err := s.MapJSON(w, r, &model)
	if err == nil {

		res, err := solver.Solve(model.Sudoku)

		if err != nil {
			errors.BadRequest(err.Error()).Send(w)
			return
		}

		solved := err == nil

		solvedSudoku := Sudoku{toString(res), solved}
		s.SendJSON(
			w,
			r,
			solvedSudoku,
			http.StatusOK,
		)
		return
	}
	err.Send(w)
}

// ToString convert the solved sudoku (map[string]string) to as string of values
func toString(solvedSudoku map[string]string) (res string) {
	keys := []string{}
	for k := range solvedSudoku {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		res += solvedSudoku[k]
	}

	return res
}
