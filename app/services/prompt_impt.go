package services

import (
	"educ-gpt/models"
	"fmt"
	"strings"
)

type PromptServiceImpl struct{}

func (p PromptServiceImpl) GetThemes(topic string, existedThemes []*models.Theme, userStats []*models.Theme) (string, error) {
	getThemesPrompt1 := "Я собираю список разделов для обучения программированию по теме: "
	getThemesPrompt2 := ". В списке уже есть эти разделы, не создавай новые разделы в списке, которые являются аналогами существующих: "
	getThemesPrompt3 := ". Учитывай прогресс пользователя. Вот список разделов, которые он уже изучил, и количество решенных задач по каждому разделу: "
	getThemesPrompt4 := `. Верни список в виде JSON массива строк, который будет включать новые от тебя разделы и предоставленные мной (если я его предоставил). `
	getThemesPrompt42 := ` Верни список в виде JSON массива строк.`
	getThemesPrompt5 := `Отсортируй разделы по сложности: от простых к сложным. Пример ответа: [\"Основы синтаксиса\",\"Работа с функциями\", ...]. В ответе должен быть только НЕ ПУСТОЙ массив JSON со строками!`

	exist := make([]string, len(existedThemes))
	for i := range existedThemes {
		exist[i] = existedThemes[i].Title
	}

	stats := make([]string, len(userStats))
	for i := range userStats {
		stats[i] = fmt.Sprintf("[%s,%v]", userStats[i].Title, userStats[i].Score)
	}

	sb := strings.Builder{}

	sb.WriteString(getThemesPrompt1)
	sb.WriteString(topic)

	if len(exist) != 0 {
		sb.WriteString(getThemesPrompt3)
		sb.WriteRune('[')
		sb.WriteString(strings.Join(stats, ", "))
		sb.WriteRune(']')
	}

	sb.WriteString(getThemesPrompt2)
	sb.WriteString(strings.Join(exist, ", "))

	if len(stats) != 0 {
		sb.WriteString(getThemesPrompt3)
		sb.WriteString(strings.Join(stats, ", "))
	}

	if len(exist) != 0 {
		sb.WriteString(getThemesPrompt4)
	} else {
		sb.WriteString(getThemesPrompt42)
	}

	sb.WriteString(getThemesPrompt5)

	return sb.String(), nil
}

func (p PromptServiceImpl) GetTheme(topic string, theme string, userStats []*models.Theme) (string, error) {
	getThemesPrompt1 := "Расскажи подробно по теме: "
	getThemesPrompt2 := ". Твой ответ должен всключать в себя примеры кода и теории не менее 10000 символов"
	getThemesPrompt3 := ". Учитывай прогресс пользователя. Вот список тем, которые он уже изучил, и количество решенных задач по каждой теме: "
	getThemesPrompt4 := ". Подготовь 10 задач по этой теме. Ответ должен быть json в виде {\"text\": <теоретический текст>, \"problems\":[<задача 1>, <задача 2>, ...]}. В ответе должен быть только json!"

	stats := make([]string, len(userStats))
	for i := range userStats {
		stats[i] = fmt.Sprintf("[%s,%v]", userStats[i].Title, userStats[i].Score)
	}

	sb := strings.Builder{}

	sb.WriteString(getThemesPrompt1)
	sb.WriteString(topic)
	sb.WriteRune('.')
	sb.WriteString(theme)
	sb.WriteString(getThemesPrompt2)
	sb.WriteString(getThemesPrompt3)
	sb.WriteString(strings.Join(stats, ", "))
	sb.WriteString(getThemesPrompt4)

	return sb.String(), nil
}

func NewPromptServiceImpl() PromptService {
	return &PromptServiceImpl{}
}
