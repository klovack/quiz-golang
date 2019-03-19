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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	showAnswer := flag.Bool("answer", true, "Show answer at every question")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the CSV file: %s\n", *csvFilename))
	}

	problems := parseLines(lines)

	// Set timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

PROBLEM_LOOP:
	for i, p := range problems {
		fmt.Printf("Question %d: %s = ", i+1, p.question)
		answerCh := make(chan string)

		go func() {
			printProblem(answerCh)
		}()

		select {
		case <-timer.C:
			fmt.Printf("\n\nTIME'S UP!!!!\n")
			break PROBLEM_LOOP
		case answer := <-answerCh:
			printResult(&answer, &p, &correct, showAnswer)
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

func printProblem(answerCh chan<- string) {
	var answer string
	_, err := fmt.Scanf("%s\n", &answer)
	if err != nil {
		exit(fmt.Sprintf("Failed to get response\n"))
	}
	answerCh <- answer
}

func printResult(answer *string, p *problem, correct *int, showAnswer *bool) {
	isCorrect := strings.ToLower(*answer) == strings.ToLower(p.answer)
	if isCorrect {
		*correct++
	}

	if *showAnswer {
		if isCorrect {
			fmt.Println("Correct!!!")
		}
		fmt.Printf("Your answer: %s.\nRight answer: %s.\n\n", *answer, p.answer)
	}
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
