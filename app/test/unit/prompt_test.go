package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/models"
	"educ-gpt/services"
	"fmt"
	"strings"
	"testing"
)

var (
	promptSrv services.PromptService
)

func convertStats(userAllStats []*models.Theme) string {
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
	return stats
}

func convertThemesWithScore(existedThemes []*models.Theme) string {
	exist := make([]string, 0)
	for i := range existedThemes {
		if existedThemes[i].Score > 0 {
			exist = append(exist, existedThemes[i].Title)
		}
	}

	return strings.Join(exist, ", ")
}

func TestCanPromptServiceInit(t *testing.T) {
	promptSrv = dic.PromptService()
}

func TestCanCompileCode(t *testing.T) {
	code := "print('Hello, World!')"
	language := "Python"

	result := promptSrv.CompileCode(code, language)

	if !strings.Contains(result, language) {
		t.Error("Missing language")
		return
	}

	if !strings.Contains(result, code) {
		t.Error("Missing code")
		return
	}

	if strings.Contains(result, "%!s(MISSING)") {
		t.Error("%!s(MISSING)")
		return
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
	result := promptSrv.GetThemes(topic, existedThemes, userStats)

	if !strings.Contains(result, topic) {
		t.Error("Missing topic")
		return
	}

	if !strings.Contains(result, convertStats(userStats)) {
		t.Error("Missing stats")
		return
	}

	if !strings.Contains(result, convertThemesWithScore(existedThemes)) {
		t.Errorf("Missing existed themes")
		return
	}

	if strings.Contains(result, "%!s(MISSING)") {
		t.Error("%!s(MISSING)")
		return
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
	result := promptSrv.GetTheme(topic, theme, userStats, userAllStats)

	if !strings.Contains(result, topic) {
		t.Error("Missing topic")
		return
	}

	if !strings.Contains(result, theme) {
		t.Error("Missing theme")
		return
	}

	if strings.Contains(result, "%!s(MISSING)") {
		t.Error("%!s(MISSING)")
		return
	}

	if !strings.Contains(result, userStats.ResolvedProblems) {
		t.Error("Missing resolved problems")
		return
	}

	if !strings.Contains(result, convertStats(userAllStats)) {
		t.Error("Missing stats")
		return
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

	result := promptSrv.GetProblems(count, topic, theme, userThemeStats, userAllStats)

	if !strings.Contains(result, topic) {
		t.Error("Missing topic")
		return
	}

	if !strings.Contains(result, theme) {
		t.Error("Missing theme")
		return
	}

	if strings.Contains(result, "%!s(MISSING)") {
		t.Error("%!s(MISSING)")
		return
	}

	if !strings.Contains(result, userThemeStats.ResolvedProblems) {
		t.Error("Missing resolved problems")
		return
	}

	if !strings.Contains(result, convertStats(userAllStats)) {
		t.Error("Missing stats")
		return
	}
}

func TestCanVerifyAnswer(t *testing.T) {
	problem := "Выведи на экран число 4"
	answer := "print(4)"
	language := "Python"

	result := promptSrv.VerifyAnswer(problem, answer, language)

	if !strings.Contains(result, problem) {
		t.Error("Missing problem")
		return
	}

	if !strings.Contains(result, language) {
		t.Error("Missing language")
		return
	}

	if !strings.Contains(result, answer) {
		t.Error("Missing answer")
		return
	}

	if strings.Contains(result, "%!s(MISSING)") {
		t.Error("%!s(MISSING)")
		return
	}
}
