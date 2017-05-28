package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/braintree/manners"
	"github.com/gorilla/mux"
	"github.com/laurentlp/sudoku-solver/api/bundles/sudoku_bundle"
	"github.com/laurentlp/sudoku-solver/api/errors"
)

func main() {

	// load error messages
	if err := errors.LoadErrors("./api/errors/error_templates.yaml"); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	// Controllers declaration
	sudoku := &sudokubundle.SudokuController{}

	// Router declaration and prefix
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1/").Subrouter()

	// Routes handling
	s.HandleFunc("/sudoku", sudoku.Solve).Methods("POST")

	// Create new Gracefulserver and bind listin address and handlers
	go func(r http.Handler) {
		manners.ListenAndServe(":8080", r)

		log.Println("The api is shutting down...")

	}(r)

	fmt.Println("Api started on localhost:8080")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press enter to stop server...")
	reader.ReadString('\n')

	fmt.Print("Server stopping...")
	manners.Close()
}
