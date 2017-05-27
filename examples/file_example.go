package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/laurentlp/sudoku-solver/solver"
)

func main() {
	// File containing sudoku(s)
	file, err := os.Open("./sudoku/_tests/hardest.txt")
	if err != nil {
		panic("Couldn't open file hardest.txt : " + err.Error())
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		fmt.Println("problem(hard): ", string(line))
		resolved, err := solver.Solve(string(line))
		if err != nil {
			break
		}

		solver.Display(resolved)
	}
}
