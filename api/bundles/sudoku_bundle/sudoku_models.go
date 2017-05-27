package sudokubundle

// Sudoku struct
type Sudoku struct {
	Sudoku string `json:"sudoku"`
	Solved bool   `json:"solved"`
}

// NewSudoku create a new sudoku
func NewSudoku(sudoku string, solved bool) *Sudoku {
	return &Sudoku{
		Sudoku: sudoku,
		Solved: solved,
	}
}
