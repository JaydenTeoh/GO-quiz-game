package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", " a csv file in the format of 'question, answer'")
	flag.Parse()

	// open the csv file, use pointer * because flag.String returns a pointer to the value of the flag
	// returns a pointer to an os.File, which is an io.Reader (implements the Read method)
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}

	// use csv.NewReader function that implements the io.Reader interface
	r := csv.NewReader(file)
	lines, err := r.ReadAll() // read all the lines in r
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}
	problems := parseLines(lines)

	correct := 0 // counter to keep track of number of correct answers

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		var answer string
		//scan the input in the CLI and stores in the answer variable
		//stops the printing of the next question and waits for input for answer to current question
		fmt.Scanf("%s\n", &answer)
		if answer == p.answer {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

// create array of problems via a CSV 2D slice: ([][]string) => []problem
func parseLines(lines [][]string) []problem {
	// we know the exact size of the array (num of problems) => no need to use empty slice as it would be extra work to resizen the slice and append
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]), // strings.Trimspace remove unneccessary spaces that might come in the answers in the CSV file.
		}
	}
	return ret
}

// expect the problems to come in this format (question stirng, answer string) regardless of file format: CSV, JSON etc.
type problem struct {
	question string
	answer   string
}

// create reusable exit function for errors
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
