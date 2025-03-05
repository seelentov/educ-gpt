package dtos

import "fmt"

type IncreaseUserScoreAndAddAnswerRequest struct {
	ProblemId uint   `json:"problem_id" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}

type Book struct {
	Title  string
	Author string
	Year   uint8
}

func (b Book) String() string {
	return fmt.Sprintf("%s %s %v", b.Title, b.Author, b.Year)
}

func PrintBookData(b *Book) {
	fmt.Println(b)
}
