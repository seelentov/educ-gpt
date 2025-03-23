package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/models"
	"educ-gpt/services"
	"testing"
)

var (
	promptSrv services.PromptService
)

func TestCanPromptServiceInit(t *testing.T) {
	promptSrv = dic.PromptService()
}

func TestCanCompileCode(t *testing.T) {
	code := "print('Hello, World!')"
	language := "Python"
	expected := "Проанализируй код на языке Python и пришли мне, что будет выведено в консоль при его выполнении. Если в коде есть ошибки, верни сообщение об ошибке. Ответ должен быть в формате JSON: {result: (текст в терминале при запуске этого кода или сообщение об ошибке)}.<начало кода>print('Hello, World!')<конец кода>"

	result := promptSrv.CompileCode(code, language)

	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestCanGetThemes(t *testing.T) {
	topic := "Programming"
	existedThemes := []*models.Theme{
		{Title: "Variables", Score: 5},
		{Title: "Loops", Score: 0},
	}
	userStats := []*models.Theme{
		{Title: "Variables", Score: 3},
		{Title: "Functions", Score: 2},
	}
	expected := "Создай список разделов для обучения программированию по теме: Programming. Учитывай прогресс пользователя: он уже изучил разделы Variables и решил задачи [Variables,3], [Functions,2]. Не дублируй существующие разделы (%!s(MISSING)). Верни список в формате JSON, отсортированный по сложности: от простых к сложным. Пример: [Основы синтаксиса, Работа с функциями, ...]. Разделы должны охватывать все уровни сложности и быть полезными для обучения. В ответе должен быть только JSON массив строк."

	result := promptSrv.GetThemes(topic, existedThemes, userStats)

	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestCanGetTheme(t *testing.T) {
	topic := "Programming"
	theme := "Variables"
	userStats := &models.Theme{Title: "Variables", Score: 3, ResolvedProblems: "1,2,3"}
	userAllStats := []*models.Theme{
		{Title: "Variables", Score: 3},
		{Title: "Functions", Score: 2},
	}
	expected := "Расскажи подробно по теме: Programming. Учитывай прогресс пользователя: он уже изучил темы [Variables : 3 задач, Functions : 2 задач, ] и решил задачи [1,2,3]. Ответ должен включать теорию и примеры кода. Если тема узкая, текст должен быть не менее 30000 символов. Подготовь 10 задач по теме, включая теоретические и практические. Задачи не должны дублировать уже решенные (%!s(MISSING)). Ответ должен быть в формате JSON: {text: <теория с примерами кода>, problems: [{question: <текст задачи>, languages: <список языков через ; или пустая строка для теории>, is_theory: <true/false>}]}. Задачи должны быть адаптированы под уровень пользователя и не привязаны к конкретному языку, если тема универсальна. Пример списка языков: Python;JavaScript;Java;C++;C#;PHP;TypeScript;Swift;Go;Kotlin;Ruby;Rust;SQL;R;Perl;Dart;Scala;Haskell;Lua;Objective-C;Shell;PowerShell;Assembly;MATLAB;Groovy;Elixir;Clojure;F#;Erlang;VBA;Delphi;Ada;Lisp;Fortran;Prolog;Cobol;Bash;Racket;Julia;Crystal;Nim;OCaml;D;Vala;Smalltalk;ABAP;ActionScript;Apex;ColdFusion;Eiffel;LabVIEW;PL/SQL;SAS;Scheme;Tcl;Verilog;VHDL;Zig. Относись к ответам на теоретические вопросы снисходительно: принимай ответ, если он не содержит ошибок, даже если он неполный. Подсказки должны направлять к правильному ответу, но не давать его напрямую."
	result := promptSrv.GetTheme(topic, theme, userStats, userAllStats)

	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestCanGetProblems(t *testing.T) {
	count := 5
	topic := "Programming"
	theme := "Variables"
	userThemeStats := &models.Theme{Title: "Variables", Score: 3, ResolvedProblems: "1,2,3"}
	userAllStats := []*models.Theme{
		{Title: "Variables", Score: 3},
		{Title: "Functions", Score: 2},
	}
	expected := "Подготовь 5 задач по теме: Programming. Учитывай прогресс пользователя: он уже изучил темы Variables и решил задачи [Variables : 3 задач, Functions : 2 задач, ]. Задачи должны включать как теоретические, так и практические вопросы. Если тема не позволяет создать практическую задачу, сделай акцент на теории. Ответ должен быть в формате JSON: [{question: <текст задачи>, languages: <список языков через ; или пустая строка для теории>, is_theory: <true/false>}]. Задачи не должны дублировать уже решенные ([1,2,3]) и должны быть адаптированы под уровень пользователя. Если задача практическая, укажи языки, на которых ее можно решить, но не привязывай задачу к конкретному языку, если тема универсальна. Пример списка языков: Python;JavaScript;Java;C++;C#;PHP;TypeScript;Swift;Go;Kotlin;Ruby;Rust;SQL;R;Perl;Dart;Scala;Haskell;Lua;Objective-C;Shell;PowerShell;Assembly;MATLAB;Groovy;Elixir;Clojure;F#;Erlang;VBA;Delphi;Ada;Lisp;Fortran;Prolog;Cobol;Bash;Racket;Julia;Crystal;Nim;OCaml;D;Vala;Smalltalk;ABAP;ActionScript;Apex;ColdFusion;Eiffel;LabVIEW;PL/SQL;SAS;Scheme;Tcl;Verilog;VHDL;Zig. Относись к ответам на теоретические вопросы снисходительно: принимай ответ, если он не содержит ошибок, даже если он неполный. Подсказки должны направлять к правильному ответу, но не давать его напрямую."

	result := promptSrv.GetProblems(count, topic, theme, userThemeStats, userAllStats)

	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestCanVerifyAnswer(t *testing.T) {
	problem := "Выведи на экран число 4"
	answer := "print(4)"
	language := "Python"
	expected := "Я получил от тебя задачу: Выведи на экран число 4. Вот мое решение: print(4). Проверь, соответствует ли оно требованиям задачи. Учти, что задача выполнена на языке Python. Ответ должен быть в формате JSON: {ok: <true/false>, message: <пояснение, если ok==false, или подсказка по улучшению, если ok==true>}. Проверяй только соответствие решения задаче, не учитывай грамматику или стиль кода."

	result := promptSrv.VerifyAnswer(problem, answer, language)

	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}
