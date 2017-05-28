# sudoku-solver

[![wercker status](https://app.wercker.com/status/3ee8307fdb3876b6c0f7a504caf5daef/s/master "wercker status")](https://app.wercker.com/project/byKey/3ee8307fdb3876b6c0f7a504caf5daef)
[![Coverage Status](https://coveralls.io/repos/github/laurentlp/sudoku-solver/badge.svg)](https://coveralls.io/github/laurentlp/sudoku-solver)
[![Go Report Card](https://goreportcard.com/badge/github.com/laurentlp/sudoku-solver)](https://goreportcard.com/report/github.com/laurentlp/sudoku-solver)

A simple golang sudoku solver using Peter Norvig algorithm

## Installation

```shell
$ go get github.com/laurentlp/sudoku-solver
```

## To install the dependencies (Tests + API)

### Install glide (Package Management for Golang)

```shell
$ curl https://glide.sh/get | sh
```

Or via Homebrew on MacOS

```shell
$ brew install glide
```

Then to install the dependencies the dependencies

```shell
$ glide i
```

## Running the tests

### Test the solver

```shell
$ cd $GOPATH/src/github.com/laurentlp/sudoku-solver
$ go test -v ./solver
```

### Test everything

Simply run the `test.sh` file like so

```shell
$ cd $GOPATH/src/github.com/laurentlp/sudoku-solver
$ ./test.sh
```

*Before, make sure the file is executable

```bash
$ chmod 777 test.sh
```

## API

To use the API, simply run the `main.go` file in the command line.
Then enter your sudoku in the body of the request as shown bellow :

Make a POST request to `http://localhost:8080/sudoku`

with the body (JSON) :

```json
{
    "sudoku" : "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
}
```

Result :

![solved.jpg from the examples folder](https://raw.githubusercontent.com/laurentlp/sudoku-solver/master/examples/solved.jpeg)

## Basic usage of the solver

First create a separate go project in which you will need a `main.go` file

```golang
package main

import (
    "fmt"

    "github.com/laurentlp/sudoku-solver/solver"
)

func main() {
    // Directly write a sudoku grid or use a file as shown in the examples folder
    grid := "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
    resolved, err := solver.Solve(grid)

    if err != nil {
        fmt.Print(err)
    }

    solver.Display(resolved)
}
```

Now to run the file just hit `go run main.go` in the command line at the root of your project.

You should see an output looking like this :

```bash
$ go run main.go

4 1 7 | 3 6 9 | 8 2 5
6 3 2 | 1 5 8 | 9 4 7
9 5 8 | 7 2 4 | 3 1 6
------+-------+-------
8 2 5 | 4 3 7 | 1 6 9
7 9 1 | 5 8 6 | 4 3 2
3 4 6 | 9 1 2 | 7 5 8
------+-------+-------
2 8 9 | 6 4 3 | 5 7 1
5 7 3 | 2 9 1 | 6 8 4
1 6 4 | 8 7 5 | 2 9 3
```