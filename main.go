package main

import (
	"bufio"
	"fmt"
	"log"
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
	server := manners.NewServer()
	server.Addr = ":8080"
	server.Handler = r

	go func(s *manners.GracefulServer) {
		err := s.ListenAndServe()

		log.Println("The api is shutting down...")

		if err != nil {
			panic(err)
		}

	}(server)

	fmt.Println("Api started on localhost", server.Addr)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press enter to stop server...")
	reader.ReadString('\n')

	fmt.Print("Server stopping...")
	server.Close()
}
