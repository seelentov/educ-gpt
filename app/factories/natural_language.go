package factories

import "educ-gpt/services"

type ChatType int

const (
	GPT ChatType = iota
)

type NaturalLanguageServiceFactory interface {
	Get(t ChatType) services.NaturalLanguageService
}
