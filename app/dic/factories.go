package dic

import (
	"educ-gpt/factories"
	"educ-gpt/logger"
)

var naturalLanguageServiceFactory factories.NaturalLanguageServiceFactory

func NaturalLanguageServiceFactory() factories.NaturalLanguageServiceFactory {
	if naturalLanguageServiceFactory == nil {
		naturalLanguageServiceFactory = factories.NewNaturalLanguageServiceFactoryImpl(
			GptService(),
		)
		logger.Logger().Debug("NaturalLanguageServiceFactory initialized")
	}

	return naturalLanguageServiceFactory
}
