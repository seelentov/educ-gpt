package services

import (
	"educ-gpt/models"
	"fmt"
	"strings"
)

type PromptServiceImpl struct{}

func (p PromptServiceImpl) CompileCode(code string) string {
	prompt1 := "Пришли мне в ответном сообщении, что выведет в консоль этот код: "
	prompt2 := `Ответ должен быть в формате JSON такого вида: {"result":"(текст в терминале при запуске этого кода)"}!`

	sb := strings.Builder{}

	sb.WriteString(prompt1)
	sb.WriteString(code)
	sb.WriteString(prompt2)

	return sb.String()
}

func (p PromptServiceImpl) GetThemes(topic string, existedThemes []*models.Theme, userStats []*models.Theme) string {
	prompt1 := "Я собираю список разделов для обучения программированию по теме: "
	prompt2 := ". В списке уже есть эти разделы, не создавай новые разделы в списке, которые являются аналогами существующих: "
	prompt3 := ". Учитывай прогресс пользователя. Вот список разделов, которые он уже изучил, и количество решенных задач по каждому разделу: "
	prompt4 := `. Верни список в виде JSON массива строк, который будет включать новые от тебя разделы и предоставленные мной (если я его предоставил). `
	prompt42 := ` Верни список в виде JSON массива строк.`
	prompt5 := `Отсортируй разделы по сложности: от простых к сложным. Пример ответа: [\"Основы синтаксиса\",\"Работа с функциями\", ...]. Новых разделов должно быть не менее 50, т.е. 50 твоих разделов + предоставленные мной. Они должны охватывать как начинающий уровень, так и продвинутый и более высокий по сложности. В ответе должен быть только НЕ ПУСТОЙ массив JSON со строками!`

	exist := make([]string, 0)

	if existedThemes != nil {
		for i := range existedThemes {
			if existedThemes[i].Score > 0 {
				exist = append(exist, existedThemes[i].Title)
			}
		}
	}

	var stats []string

	if userStats != nil {
		stats = make([]string, len(userStats))
		for i := range userStats {
			stats[i] = fmt.Sprintf("[%s,%v]", userStats[i].Title, userStats[i].Score)
		}
	}

	sb := strings.Builder{}

	sb.WriteString(prompt1)
	sb.WriteString(topic)

	if len(exist) != 0 {
		sb.WriteString(prompt3)
		sb.WriteRune('[')
		sb.WriteString(strings.Join(stats, ", "))
		sb.WriteRune(']')
	}

	sb.WriteString(prompt2)
	sb.WriteString(strings.Join(exist, ", "))

	if len(stats) != 0 {
		sb.WriteString(prompt3)
		sb.WriteString(strings.Join(stats, ", "))
	}

	if len(exist) != 0 {
		sb.WriteString(prompt4)
	} else {
		sb.WriteString(prompt42)
	}

	sb.WriteString(prompt5)

	return sb.String()
}

func (p PromptServiceImpl) GetTheme(topic string, theme string, userThemeStats *models.Theme, userAllStats []*models.Theme) string {
	prompt1 := "Расскажи подробно по теме: "
	prompt2 := ". Твой ответ должен включать в себя примеры кода и теории не менее 30000 символов, не считая кода"
	prompt3 := ". Учитывай прогресс пользователя. Вот список тем, которые он уже изучил, и количество решенных задач по каждой теме: "
	prompt4 := ". Подготовь 10 задач по этой теме."
	prompt5 := ". Вот список задач, которые уже выполнил пользователь, их не должно быть в списке задач от тебя:"
	prompt6 := ` Ответ должен быть json в виде {\"text\": \"<теоретический текст>\", \"problems\":[\"<задача 1>\", \"<задача 2>\", ...]}. Если тема не связанна с конкретной технологией или языком программирования, то технология или язык программирования может быть любым. В ответе должен быть только JSON объект с text, который будет включать основной ответ и problems с задачами в виде строк!`

	sb := strings.Builder{}

	sb.WriteString(prompt1)
	sb.WriteString(topic)
	sb.WriteRune('.')
	sb.WriteString(theme)
	sb.WriteString(prompt2)

	if userAllStats != nil && len(userAllStats) > 0 {
		sb.WriteString(prompt3)
		sb.WriteRune('[')
		for i := range userAllStats {
			sb.WriteString(fmt.Sprintf("%s : %v задач, ", userAllStats[i].Title, userAllStats[i].Score))
		}
		sb.WriteRune(']')
	}

	sb.WriteString(prompt4)

	if userThemeStats != nil && userThemeStats.ResolvedProblems != "" {
		sb.WriteString(prompt5)
		sb.WriteRune('[')
		sb.WriteString(userThemeStats.ResolvedProblems)
		sb.WriteRune(']')
	}
	sb.WriteString(prompt6)

	return sb.String()
}

func (p PromptServiceImpl) GetProblems(count int, topic string, theme string, userThemeStats *models.Theme, userAllStats []*models.Theme) string {
	prompt1 := fmt.Sprintf("Подготовь 10 задачи в количестве %v по теме: ", count)
	prompt2 := ". Учитывай прогресс пользователя. Вот список тем, которые он уже изучил, и количество решенных задач по каждой теме: "
	prompt3 := ". Вот список задач, которые уже выполнил пользователь, их не должно быть в списке задач от тебя:"
	prompt4 := `Ответ должен быть json в виде {\"text\": \"<теоретический текст>\", \"problems\":[\"<задача 1>\", \"<задача 2>\", ...]}. Если тема не связанна с конкретной технологией или языком программирования, то технология или язык программирования может быть любым. В ответе должен быть только JSON объект с text, который будет включать основной ответ и problems с задачами в виде строк! Так же текст должен быть правиьно размечен (текст, заголовки, жирный, курсив, программный код), т.к. я буду рендерить его в markdown компоненте.`

	sb := strings.Builder{}

	sb.WriteString(prompt1)
	sb.WriteString(topic)
	sb.WriteRune('.')
	sb.WriteString(theme)
	sb.WriteRune(' ')

	notEmptyUserAllStats := make([]*models.Theme, 0)

	if userAllStats != nil {
		for i := range userAllStats {
			if userAllStats[i].Score > 0 {
				notEmptyUserAllStats = append(notEmptyUserAllStats, userAllStats[i])
			}
		}
	}

	if len(notEmptyUserAllStats) > 0 {
		sb.WriteString(prompt2)
		sb.WriteRune('[')
		for i := range notEmptyUserAllStats {
			sb.WriteString(fmt.Sprintf("%s : %v задач, ", notEmptyUserAllStats[i].Title, notEmptyUserAllStats[i].Score))
		}
		sb.WriteRune(']')
	}

	sb.WriteString(prompt4)

	if userThemeStats != nil && userThemeStats.ResolvedProblems != "" {
		sb.WriteString(prompt3)
		sb.WriteRune('[')
		sb.WriteString(userThemeStats.ResolvedProblems)
		sb.WriteRune(']')
	}

	sb.WriteString(prompt4)

	return sb.String()
}

func (p PromptServiceImpl) VerifyAnswer(problem string, answer string) string {
	prompt1 := "Я получил от тебя задачу: "
	prompt2 := " Вот мое решение: "
	prompt3 := `. Соответствует ли мой ответ требованиям задачи? Я жду от тебя ответ в формате JSON : {ok: <Булево значение, соответствует ли решение задаче>, message:<Если ok==false, то тут должно быть короткое пояснение в виде строки, если ты не принимаешь задачу или подсказку к улучшению кода, если задача выполнена>}. Если ответ пользователя решает задачу то прими ее, подсказав, как решение можно улучшить в message. В ответе должен быть только JSON в описанным мной ранее формате!`

	sb := strings.Builder{}

	answer = strings.ReplaceAll(answer, `"`, `\"`)

	sb.WriteString(prompt1)
	sb.WriteString(problem)
	sb.WriteString(prompt2)
	sb.WriteString(answer)
	sb.WriteString(prompt3)

	return sb.String()
}

func NewPromptServiceImpl() PromptService {
	return &PromptServiceImpl{}
}
