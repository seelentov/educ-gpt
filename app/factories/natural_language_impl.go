package factories

import "educ-gpt/services"

type NaturalLanguageServiceFactoryImpl struct {
	gpt services.GptService
}

func (n NaturalLanguageServiceFactoryImpl) Get(t ChatType) services.NaturalLanguageService {
	switch t {
	case GPT:
		return n.gpt
	default:
		return n.gpt
	}
}

func NewNaturalLanguageServiceFactoryImpl(gpt services.GptService) NaturalLanguageServiceFactory {
	return &NaturalLanguageServiceFactoryImpl{gpt}
}
