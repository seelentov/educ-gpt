package services

import (
	"educ-gpt/models"
	"fmt"
	"strings"
)

type PromptServiceImpl struct{}

func (p PromptServiceImpl) GetThemes(topic string, existedThemes []*models.Theme, userStats []*models.Theme) (string, error) {
	getThemesPrompt1 := "Я собираю список разделов для обучения программированию по теме : "
	getThemesPrompt2 := ". Исключи из новых разделов, которые ты подберешь мне следующие разделы (список может быть пустой, я передаю его автоматически): "
	getThemesPrompt3 := ". Учитывай прогресс пользователя. Вот список разделов, которые он уже изучил, и количество решенных задач по каждому разделу: "
	getThemesPrompt4 := ". Верни список в виде JSON массива строк,который будет включать новые от тебя разделы и предоставленные мной (список может быть пустой, я передаю его автоматически).Отсортируй разделы по сложности: от простых к сложным. Пример ответа: [\"Основы синтаксиса\",\"Работа с функциями\", ...]. "

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
	sb.WriteString(getThemesPrompt2)
	sb.WriteString(strings.Join(exist, ", "))

	if len(stats) != 0 {
		sb.WriteString(getThemesPrompt3)
		sb.WriteString(strings.Join(stats, ", "))
	}

	sb.WriteString(getThemesPrompt4)

	return sb.String(), nil
}

func (p PromptServiceImpl) GetTheme(topic string, theme string, userStats []*models.Theme) (string, error) {
	getThemesPrompt1 := "Расскажи подробно по теме: "
	getThemesPrompt2 := ". Твой ответ должен всключать в себя примеры кода и теории не менее 10000 символов"
	getThemesPrompt3 := ". Учитывай прогресс пользователя. Вот список тем, которые он уже изучил, и количество решенных задач по каждой теме: "
	getThemesPrompt4 := ". Подготовь 10 задач по этой теме. Ответ должен быть json в виде {\"text\": <теоретический текст>, \"problems\":[<задача 1>, <задача 2>, ...]}"

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
