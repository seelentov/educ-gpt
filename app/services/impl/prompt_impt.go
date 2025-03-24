package impl

import (
	"educ-gpt/models"
	"educ-gpt/services"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PromptServiceImpl struct{}

func readPromptFile(filename string) string {
	path := filepath.Join("./resources/prompt", filename)
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Ошибка при чтении файла промпта: %v", err)
	}
	return string(content)
}

func (p PromptServiceImpl) CompileCode(code string, language string) string {
	prompt := readPromptFile("compile_code_prompt.txt")
	return fmt.Sprintf(prompt, language) + fmt.Sprintf("<начало кода>%s<конец кода>", code)
}

func (p PromptServiceImpl) GetThemes(topic string, existedThemes []*models.Theme, userStats []*models.Theme) string {
	prompt := readPromptFile("get_themes_prompt.txt")

	exist := make([]string, 0)
	for i := range existedThemes {
		if existedThemes[i].Score > 0 {
			exist = append(exist, existedThemes[i].Title)
		}
	}

	stats := "["
	for i := range userStats {
		stats += fmt.Sprintf("%s : %v задач, ", userStats[i].Title, userStats[i].Score)
	}
	stats += "]"

	return fmt.Sprintf(prompt, topic, strings.Join(exist, ", "), stats)
}

func (p PromptServiceImpl) GetTheme(topic string, theme string, userStats *models.Theme, userAllStats []*models.Theme) string {
	prompt := readPromptFile("get_theme_prompt.txt")

	notEmptyUserAllStats := make([]*models.Theme, 0)
	for i := range userAllStats {
		if userAllStats[i].Score > 0 {
			notEmptyUserAllStats = append(notEmptyUserAllStats, userAllStats[i])
		}
	}

	stats := ""
	if len(notEmptyUserAllStats) > 0 {
		stats = "["
		for i := range notEmptyUserAllStats {
			stats += fmt.Sprintf("%s : %v задач, ", notEmptyUserAllStats[i].Title, notEmptyUserAllStats[i].Score)
		}
		stats += "]"
	}

	resolvedProblems := ""
	if userStats != nil && userStats.ResolvedProblems != "" {
		resolvedProblems = "[" + userStats.ResolvedProblems + "]"
	}

	return fmt.Sprintf(prompt, topic+" "+theme, stats, resolvedProblems)
}

func (p PromptServiceImpl) GetProblems(count int, topic string, theme string, userThemeStats *models.Theme, userAllStats []*models.Theme) string {
	prompt := readPromptFile("get_problems_prompt.txt")

	notEmptyUserAllStats := make([]*models.Theme, 0)
	for i := range userAllStats {
		if userAllStats[i].Score > 0 {
			notEmptyUserAllStats = append(notEmptyUserAllStats, userAllStats[i])
		}
	}

	stats := ""
	if len(notEmptyUserAllStats) > 0 {
		stats = "["
		for i := range notEmptyUserAllStats {
			stats += fmt.Sprintf("%s : %v задач, ", notEmptyUserAllStats[i].Title, notEmptyUserAllStats[i].Score)
		}
		stats += "]"
	}

	resolvedProblems := ""
	if userThemeStats != nil && userThemeStats.ResolvedProblems != "" {
		resolvedProblems = "[" + userThemeStats.ResolvedProblems + "]"
	}

	return fmt.Sprintf(prompt, count, topic+" "+theme, stats, resolvedProblems)
}

func (p PromptServiceImpl) VerifyAnswer(problem string, answer string, language string) string {
	prompt := readPromptFile("verify_answer_prompt.txt")
	return fmt.Sprintf(prompt, problem, answer, language)
}

func NewPromptServiceImpl() services.PromptService {
	return &PromptServiceImpl{}
}
