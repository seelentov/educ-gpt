package controllers

import (
	"educ-gpt/http/dtos"
	"educ-gpt/http/httputils"
	"educ-gpt/http/validator"
	"educ-gpt/models"
	"educ-gpt/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RoadmapController struct {
	userSrv    services.UserService
	gptSrv     services.GptService
	promptSrv  services.PromptService
	roadmapSrv services.RoadmapService
}

func (r RoadmapController) GetTopics(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	topics, err := r.roadmapSrv.GetTopics(userid, true)
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

	err = r.gptSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	topicId, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	topic, err := r.roadmapSrv.GetTopic(0, uint(topicId), false)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	prompt, err := r.promptSrv.GetTheme(topic.Title, theme.Title, theme)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var target services.PromptThemeRequest

	err = r.gptSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err = r.roadmapSrv.CreateProblems(target.Problems)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, target)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validator.ParseValidationErrors(err)})
		return
	}

	themeId, err := strconv.ParseUint(ctx.Param("theme_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	problem, err := r.roadmapSrv.GetProblem(req.ProblemId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	prompt, err := r.promptSrv.VerifyAnswer(problem, req.Answer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var target services.PromptProblemRequest

	err = r.gptSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if target.Ok {
		err = r.roadmapSrv.IncrementUserScoreAndAddAnswer(userId, uint(themeId), req.Answer, 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	ctx.JSON(http.StatusOK, target)
}

func NewRoadmapController(
	userSrv services.UserService,
	gptSrv services.GptService,
	promptSrv services.PromptService,
	roadmapSrv services.RoadmapService,
) *RoadmapController {
	return &RoadmapController{userSrv, gptSrv, promptSrv, roadmapSrv}
}
