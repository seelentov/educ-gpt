package services

import (
	"educ-gpt/models"
	"fmt"
	"strings"
)

type PromptServiceImpl struct{}

func (p PromptServiceImpl) CompileCode(code string, language string) string {
	prompt1 := "Пришли мне в ответном сообщении, что выведет в консоль этот код на языке " + language + ": "
	prompt2 := " Ответ должен быть в формате JSON такого вида: {result:(текст в терминале при запуске этого кода)}!"

	sb := strings.Builder{}

	sb.WriteString(prompt1)
	sb.WriteString("<начало кода>")
	sb.WriteString(code)
	sb.WriteString("<конец кода>")
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
	prompt2 := ". Твой ответ должен включать в себя примеры кода и теории не менее 15000 символов, если тема является более узкоспециализированной и не менее 50000 если тема обобщенная, и включает в себя многое, о чем можно было бы рассказать, не считая кода"
	prompt3 := ". Учитывай прогресс пользователя. Вот список тем, которые он уже изучил, и количество решенных задач по каждой теме: "
	prompt4 := ". Подготовь 10 задач по этой теме, в этих 10 задачах должны быть как задачи на решение с помощью написания кода, так и теоретические в соотношении 1 к 1"
	prompt5 := ". Вот список задач, которые уже выполнил пользователь, их не должно быть в списке задач от тебя:"
	prompt6 := ` Ответ должен быть json в виде {\"text\": \"<теоретический текст>\", \"problems\":[{\"question\": \"<текст задачи>\", \"languages\":\"<список языков программирования через символ ;, на которых можно было бы решить эту задачу. Если тема является конкретным языком программирования, то поле должно содержать только этот язык. Если вопрос теоретический, то это будет пустая строка>\", \"is_theory\":<булевое значение, является ли задача теоретической>}, ...]}. В ответе должен быть только JSON объект с text, который будет включать основной ответ и problems с задачами, без нумерации! Так же учти, что задачи будут решаться в онлайн редакторе кода, если они предполагают решение через написание кода, либо через текст, если вопрос теоретический. Если задача не теоретическая, может быть решена на множестве языков программирования и тема не закреплена за конкретный язык, то в languages добавь список этих языков программирования. Сами задачи не должны быть закреплены за определенный язык программирования, если тема может относиться к множеству языков. Например подобной задачи, как: Создайте простое RESTful API на Node.js, которое возвращает список пользователей, быть не должно, она должна быть: Создайте простое RESTful API, которое возвращает список пользователей, и в languages указан список языков, на которых ее можно было бы решить.Пример списка языков программирования: \"Python;JavaScript;Java;C++;C#;PHP;TypeScript;Swift;Go;Kotlin;Ruby;Rust;SQL;R;Perl;Dart;Scala;Haskell;Lua;Objective-C;Shell;PowerShell;Assembly;MATLAB;Groovy;Elixir;Clojure;F#;Erlang;VBA;Delphi;Ada;Lisp;Fortran;Prolog;Cobol;Bash;Racket;Julia;Crystal;Nim;OCaml;D;Vala;Smalltalk;ABAP;ActionScript;Apex;ColdFusion;Eiffel;LabVIEW;PL/SQL;SAS;Scheme;Tcl;Verilog;VHDL;Zig\"`

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
	prompt1 := fmt.Sprintf(". Подготовь %v задач по теме: %s.%s , в этих задачах должны быть как задачи на решение с помощью написания кода, так и теоретические в соотношении 1 к 1", count, topic, topic)
	prompt2 := ". Учитывай прогресс пользователя. Вот список тем, которые он уже изучил, и количество решенных задач по каждой теме: "
	prompt3 := ". Вот список задач, которые уже выполнил пользователь, их не должно быть в списке задач от тебя:"
	prompt4 := ` Ответ должен быть json в виде [{\"question\": \"<текст задачи>\", \"languages\":\"<список языков программирования через символ ;, на которых можно было бы решить эту задачу. Если тема является конкретным языком программирования, то поле должно содержать только этот язык. Если вопрос теоретический, то это будет пустая строка>\", \"is_theory\":<булевое значение, является ли задача теоретической>}, ...]. В ответе должен быть только JSON массив с задачами, без нумерации! Так же учти, что задачи будут решаться в онлайн редакторе кода, если они предполагают решение через написание кода, либо через текст, если вопрос теоретический. Если задача не теоретическая, может быть решена на множестве языков программирования и тема не закреплена за конкретный язык, то в languages добавь список этих языков программирования. Сами задачи не должны быть закреплены за определенный язык программирования, если тема может относиться к множеству языков. Например подобной задачи, как: Создайте простое RESTful API на Node.js, которое возвращает список пользователей, быть не должно, она должна быть: Создайте простое RESTful API, которое возвращает список пользователей, и в languages указан список языков, на которых ее можно было бы решить. Пример списка языков программирования: \"Python;JavaScript;Java;C++;C#;PHP;TypeScript;Swift;Go;Kotlin;Ruby;Rust;SQL;R;Perl;Dart;Scala;Haskell;Lua;Objective-C;Shell;PowerShell;Assembly;MATLAB;Groovy;Elixir;Clojure;F#;Erlang;VBA;Delphi;Ada;Lisp;Fortran;Prolog;Cobol;Bash;Racket;Julia;Crystal;Nim;OCaml;D;Vala;Smalltalk;ABAP;ActionScript;Apex;ColdFusion;Eiffel;LabVIEW;PL/SQL;SAS;Scheme;Tcl;Verilog;VHDL;Zig\"`

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

func (p PromptServiceImpl) VerifyAnswer(problem string, answer string, language string) string {
	prompt1 := "Я получил от тебя задачу: "
	prompt2 := " Вот мое решение: "
	prompt3 := ` Соответствует ли мой ответ требованиям задачи? Я жду от тебя ответ в формате JSON : {ok: <Булево значение, соответствует ли решение задаче>, message:<Если ok==false, то тут должно быть пояснение в виде строки, если ты не принимаешь задачу или подсказку к более правильному ответу, если задача выполнена>}. Если ответ пользователя решает задачу то прими ее, подсказав, как решение можно улучшить в message. В ответе должен быть только JSON в описанным мной ранее формате!`
	prompt4 := ". Задача выполнена на " + language + ". Учти это при ответе."

	sb := strings.Builder{}

	sb.WriteString(prompt1)
	sb.WriteString(problem)
	sb.WriteString(prompt2)
	sb.WriteString("<начало кода>")
	sb.WriteString(answer)
	sb.WriteString("<конец кода>")
	sb.WriteString(prompt3)

	if language != "" {
		sb.WriteString(prompt4)
	}

	return sb.String()
}

func NewPromptServiceImpl() PromptService {
	return &PromptServiceImpl{}
}
