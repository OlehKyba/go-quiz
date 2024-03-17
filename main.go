package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type problem struct {
	question string
	options  []string
	answer   int
}

func abort(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) ([]problem, error) {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		length := len(line)
		if length < 4 {
			return []problem{}, errors.New(
				fmt.Sprintf(
					"not enough parameters, check if the %d row has: question, minimum 2 options and answer",
					i+1,
				),
			)
		}

		options := line[1 : length-1]
		answer, err := strconv.Atoi(line[length-1])
		if err != nil || answer > len(options) || answer <= 0 {
			return []problem{}, errors.New(
				fmt.Sprintf(
					"invalid answer for %d row: must be an option number, got: %s",
					i+1,
					line[length-1],
				),
			)
		}

		problems[i] = problem{
			question: line[0],
			options:  options,
			answer:   answer,
		}
	}
	return problems, nil
}

func main() {
	fileName := flag.String("csv", "./problems/example.csv", "a csv file in format 'question;...options;answer'")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		abort(fmt.Sprintf("Failed to open the CSV file: %s", *fileName))
	}
	defer func() {
		onFileCloseErr := file.Close()
		if onFileCloseErr != nil {
			abort(fmt.Sprintf("Failed to close the CSV file: %s", *fileName))
		}
	}()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		abort(fmt.Sprintf("Failed to parse the provided CSV file: %s", err.Error()))
	}

	problems, err := parseLines(lines)
	if err != nil {
		abort(err.Error())
	}

	correctAnswers := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s\n", i+1, p.question)
		for j, option := range p.options {
			fmt.Printf("%d) %s\n", j+1, option)
		}

		var userInput string
		_, err := fmt.Scanf("%s\n", &userInput)
		if err != nil {
			fmt.Printf("Failed to read your answer: %s\n", err.Error())
			continue
		}

		answer, err := strconv.Atoi(userInput)
		if err != nil {
			fmt.Printf("Failed to read your answer. Expected number, but got: %s\n", userInput)
		}

		if answer == p.answer {
			correctAnswers++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correctAnswers, len(problems))
}
