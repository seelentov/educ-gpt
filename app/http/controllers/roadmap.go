package controllers

import (
	"educ-gpt/http/dtos"
	"educ-gpt/http/httputils"
	"educ-gpt/http/httputils/valid"
	"educ-gpt/models"
	"educ-gpt/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type RoadmapController struct {
	userSrv    services.UserService
	nlSrv      services.NaturalLanguageService
	promptSrv  services.PromptService
	roadmapSrv services.RoadmapService
}

func (r RoadmapController) GetTopics(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	topics, err := r.roadmapSrv.GetTopics(userid, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, topics)
}

func (r RoadmapController) GetThemes(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := r.userSrv.GetById(userid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	topicId, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	topic, err := r.roadmapSrv.GetTopic(userid, uint(topicId), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userStats := make([]*models.Theme, 0)
	for _, theme := range topic.Themes {
		if theme.Score != 0 {
			userStats = append(userStats, theme)
		}
	}

	prompt, err := r.promptSrv.GetThemes(topic.Title, topic.Themes, userStats)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var target []string

	err = r.nlSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	sortedThemes := make([]*models.Theme, len(target))
	newThemes := make([]*models.Theme, 0)

	for i := range target {
		title := target[i]

		exist := false
		existIndex := 0

		for j := range topic.Themes {
			if topic.Themes[j].Title == title {
				exist = true
				existIndex = j
				break
			}
		}

		if exist {
			sortedThemes[i] = topic.Themes[existIndex]
		} else {
			sortedThemes[i] = &models.Theme{Title: title, TopicID: topic.ID}
			newThemes = append(newThemes, sortedThemes[i])
		}

	}

	if len(newThemes) > 0 {
		err = r.roadmapSrv.CreateThemes(newThemes)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	ctx.JSON(http.StatusOK, sortedThemes)
}

func (r RoadmapController) GetTheme(ctx *gin.Context) {
	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := r.userSrv.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	topicId, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	topic, err := r.roadmapSrv.GetTopic(userId, uint(topicId), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	themeId, err := strconv.ParseUint(ctx.Param("theme_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	theme, err := r.roadmapSrv.GetTheme(userId, uint(themeId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	prompt, err := r.promptSrv.GetTheme(topic.Title, theme.Title, theme, topic.Themes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var target services.PromptThemeRequest

	err = r.nlSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	problems, err := r.roadmapSrv.CreateProblems(target.Problems, uint(themeId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"text":     target.Text,
		"problems": problems,
	})
}

func (r RoadmapController) GetProblems(ctx *gin.Context) {
	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := r.userSrv.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	topicId, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	topic, err := r.roadmapSrv.GetTopic(userId, uint(topicId), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	themeId, err := strconv.ParseUint(ctx.Param("theme_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	theme, err := r.roadmapSrv.GetTheme(userId, uint(themeId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	prompt, err := r.promptSrv.GetProblems(10, topic.Title, theme.Title, theme, topic.Themes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var target services.PromptProblemsRequest

	err = r.nlSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	problems, err := r.roadmapSrv.CreateProblems(target.Problems, uint(themeId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, problems)
}

func (r RoadmapController) IncrementUserScoreAndAddAnswer(ctx *gin.Context) {
	var req dtos.IncreaseUserScoreAndAddAnswerRequest

	userId, err := httputils.GetUserId(ctx)
	if err != nil {

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := r.userSrv.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	problem, err := r.roadmapSrv.GetProblem(req.ProblemId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	prompt, err := r.promptSrv.VerifyAnswer(problem.Question, req.Answer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var target services.PromptProblemRequest

	err = r.nlSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if target.Ok {
		err = r.roadmapSrv.IncrementUserScoreAndAddAnswer(userId, req.ProblemId, 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	ctx.JSON(http.StatusOK, target)
}

func (r RoadmapController) GetDailyChallenge(ctx *gin.Context) {

}

func NewRoadmapController(
	userSrv services.UserService,
	nlSrv services.GptService,
	promptSrv services.PromptService,
	roadmapSrv services.RoadmapService,
) *RoadmapController {
	return &RoadmapController{userSrv, nlSrv, promptSrv, roadmapSrv}
}
