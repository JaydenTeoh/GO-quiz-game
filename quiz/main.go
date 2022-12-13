package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", " a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
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

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0 // counter to keep track of number of correct answers

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		answerC := make(chan string) //make a channel answerC that receives string values
		go func() {                  /*create goroutine such that the Scanf statement can run concurrently with the select statement
			=> it won't be blocking; if timer runs out, the final score will print even if current answer is unanswered/not entered*/
			var answer string
			//scan the input in the CLI and stores in the answer variable
			//stops the printing of the next question and waits for input for answer to current question
			fmt.Scanf("%s\n", &answer)
			answerC <- answer //send answer string to the answerC channel
		}()

		select {
		case <-timer.C: //timer.C blocks on the timer's channel `C` until it sends a value indicating that the timer expired.
			fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerC: //case if answerC receives a value & initialize a variable named answer with that value
			if answer == p.answer {
				correct++
			}
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
