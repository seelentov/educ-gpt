package unit

import (
	"educ-gpt/config/dic"
	"educ-gpt/models"
	"educ-gpt/services"
	"testing"
	"time"
)

var (
	roadMapSrv services.RoadmapService
)

func TestCanInitRoadmapService(t *testing.T) {
	roadMapSrv = dic.RoadmapService()
}

func TestCreateAndGetTheme(t *testing.T) {

	theme := &models.Theme{
		Title:   "Test Theme",
		TopicID: 1,
	}

	err := roadMapSrv.CreateThemes([]*models.Theme{theme})
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	retrievedTheme, err := roadMapSrv.GetTheme(0, theme.ID, false)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	if retrievedTheme == nil {
		t.Error("Expected theme to be not nil")
		return
	}

	if retrievedTheme.Title != theme.Title {
		t.Errorf("Expected theme title to be %s, but got %s", theme.Title, retrievedTheme.Title)
	}
}

func TestCreateAndGetTopic(t *testing.T) {

	theme := &models.Theme{
		Title:   "Test Theme for Topic",
		TopicID: 1,
	}

	err := roadMapSrv.CreateThemes([]*models.Theme{theme})
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	retrievedTheme, err := roadMapSrv.GetTheme(0, theme.ID, true)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	if retrievedTheme.Topic == nil {
		t.Error("Expected topic to be not nil")
		return
	}

	if retrievedTheme.Topic.ID != theme.TopicID {
		t.Errorf("Expected topic ID to be %d, but got %d", theme.TopicID, retrievedTheme.Topic.ID)
	}
}

func TestCreateAndGetProblems(t *testing.T) {

	problem := &models.Problem{
		Question:  "Test Problem",
		ThemeID:   1,
		Languages: "Python",
		IsTheory:  false,
	}

	createdProblems, err := roadMapSrv.CreateProblems([]*models.Problem{problem})
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	if len(createdProblems) == 0 {
		t.Error("Expected at least one problem to be created")
		return
	}

	retrievedProblem, err := roadMapSrv.GetProblem(createdProblems[0].ID)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	if retrievedProblem == nil {
		t.Error("Expected problem to be not nil")
		return
	}

	if retrievedProblem.Question != problem.Question {
		t.Errorf("Expected problem question to be %s, but got %s", problem.Question, retrievedProblem.Question)
	}
}

func TestIncrementUserScoreAndAddAnswer(t *testing.T) {

	problem := &models.Problem{
		Question:  "Test Problem for Score",
		ThemeID:   1,
		Languages: "Python",
		IsTheory:  false,
	}

	createdProblems, err := roadMapSrv.CreateProblems([]*models.Problem{problem})
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	userID := uint(1)
	score := uint(10)
	err = roadMapSrv.IncrementUserScoreAndAddAnswer(userID, createdProblems[0].ID, score)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	theme, err := roadMapSrv.GetTheme(userID, problem.ThemeID, false)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	if theme.Score != score {
		t.Errorf("Expected user score to be %d, but got %d", score, theme.Score)
	}
}

func TestDeleteProblem(t *testing.T) {

	problem := &models.Problem{
		Question:  "Test Problem for Deletion",
		ThemeID:   1,
		Languages: "Python",
		IsTheory:  false,
	}

	createdProblems, err := roadMapSrv.CreateProblems([]*models.Problem{problem})
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	err = roadMapSrv.DeleteProblem(createdProblems[0].ID)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	_, err = roadMapSrv.GetProblem(createdProblems[0].ID)
	if err == nil {
		t.Error("Expected error when getting deleted problem, but got nil")
	}
}

func TestClearProblems(t *testing.T) {

	oldProblem := &models.Problem{
		Question:  "Old Problem",
		ThemeID:   1,
		Languages: "Python",
		IsTheory:  false,
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}

	_, err := roadMapSrv.CreateProblems([]*models.Problem{oldProblem})
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	err = roadMapSrv.ClearProblems()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	_, err = roadMapSrv.GetProblem(oldProblem.ID)
	if err == nil {
		t.Error("Expected error when getting deleted problem, but got nil")
	}
}

func TestClearThemes(t *testing.T) {

	oldTheme := &models.Theme{
		Title:     "Old Theme",
		TopicID:   1,
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}

	err := roadMapSrv.CreateThemes([]*models.Theme{oldTheme})
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	err = roadMapSrv.ClearThemes()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
		return
	}

	_, err = roadMapSrv.GetTheme(0, oldTheme.ID, false)
	if err == nil {
		t.Error("Expected error when getting deleted theme, but got nil")
	}
}
