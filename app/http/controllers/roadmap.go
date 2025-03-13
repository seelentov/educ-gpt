package controllers

import (
	"educ-gpt/http/dtos"
	"educ-gpt/models"
	"educ-gpt/services"
	"educ-gpt/utils/httputils"
	"educ-gpt/utils/httputils/valid"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type RoadmapController struct {
	userSrv    services.UserService
	aiSrv      services.AIService
	promptSrv  services.PromptService
	roadmapSrv services.RoadmapService
}

// GetTopics returns a list of topics for the current user
// @Summary      Get topics
// @Description  Returns a list of topics for the current user
// @Tags         roadmap
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer <JWT token>"
// @Success      200 {array} models.Topic "List of topics"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /roadmap [get]
func (r RoadmapController) GetTopics(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)

	topics, err := r.roadmapSrv.GetTopics(userid, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, topics)
}

// GetTopicInfo returns info of topic without authorization
// @Summary      Get topic info
// @Description  Returns info of topic without authorization
// @Tags         roadmap
// @Accept       json
// @Produce      json
// @Success      200 {object} models.Topic "Topic info"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /roadmap/info/topic/{topic_id} [get]
func (r RoadmapController) GetTopicInfo(ctx *gin.Context) {
	topicId, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	topic, err := r.roadmapSrv.GetTopic(0, uint(topicId), false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, topic)
}

// GetThemeInfo returns info of theme without authorization
// @Summary      Get theme info
// @Description  Returns info of theme without authorization
// @Tags         roadmap
// @Accept       json
// @Produce      json
// @Success      200 {object} models.Theme "Theme info"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /roadmap/info/theme/{theme_id} [get]
func (r RoadmapController) GetThemeInfo(ctx *gin.Context) {
	themeId, err := strconv.ParseUint(ctx.Param("theme_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	theme, err := r.roadmapSrv.GetTheme(0, uint(themeId), true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, theme)
}

// GetThemes returns a list of themes for a specific topic
// @Summary      Get themes
// @Description  Returns a list of themes for a specific topic, sorted by user progress and AI recommendations
// @Tags         roadmap
// @Accept       json
// @Produce      json
// @Param        topic_id path int true "Topic ID"
// @Param Authorization header string true "Bearer <JWT token>"
// @Success      200 {array} models.Theme "List of themes"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      404 {object} dtos.ErrorResponse "Topic not found"
// @Failure      409 {object} dtos.ErrorResponse "AI request error"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /roadmap/{topic_id} [get]
func (r RoadmapController) GetThemes(ctx *gin.Context) {
	userid, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := r.userSrv.GetById(userid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	topicId, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	topic, err := r.roadmapSrv.GetTopic(userid, uint(topicId), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	userStats := make([]*models.Theme, 0)
	for _, theme := range topic.Themes {
		if theme.Score != 0 {
			userStats = append(userStats, theme)
		}
	}

	prompt := r.promptSrv.GetThemes(topic.Title, topic.Themes, userStats)

	var target []string

	err = r.aiSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		if errors.Is(err, services.ErrAIRequestFailed) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
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
			ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
			return
		}
	}

	ctx.JSON(http.StatusOK, sortedThemes)
}

// GetTheme returns detailed information about a specific theme
// @Summary      Get theme details
// @Description  Returns detailed information about a specific theme, including problems and AI-generated content
// @Tags         roadmap
// @Accept       json
// @Produce      json
// @Param        topic_id path int true "Topic ID"
// @Param        theme_id path int true "Theme ID"
// @Param Authorization header string true "Bearer <JWT token>"
// @Success      200 {object} dtos.ThemeResponse "Theme details"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      404 {object} dtos.ErrorResponse "Theme or topic not found"
// @Failure      409 {object} dtos.ErrorResponse "AI request error"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /roadmap/{topic_id}/{theme_id} [get]
func (r RoadmapController) GetTheme(ctx *gin.Context) {
	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := r.userSrv.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	topicId, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	topic, err := r.roadmapSrv.GetTopic(userId, uint(topicId), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	themeId, err := strconv.ParseUint(ctx.Param("theme_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	theme, err := r.roadmapSrv.GetTheme(userId, uint(themeId), false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	prompt := r.promptSrv.GetTheme(topic.Title, theme.Title, theme, topic.Themes)

	var target services.PromptThemeResponse

	err = r.aiSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		if errors.Is(err, services.ErrAIRequestFailed) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	problems, err := r.roadmapSrv.CreateProblems(target.Problems, uint(themeId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, dtos.ThemeResponse{Text: target.Text, Problems: problems})
}

// GetProblems returns a list of problems for a specific theme
// @Summary      Get problems
// @Description  Returns a list of problems for a specific theme, generated by AI
// @Tags         roadmap
// @Accept       json
// @Produce      json
// @Param        topic_id path int true "Topic ID"
// @Param        theme_id path int true "Theme ID"
// @Param Authorization header string true "Bearer <JWT token>"
// @Success      200 {array} models.Problem "List of problems"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      404 {object} dtos.ErrorResponse "Theme or topic not found"
// @Failure      409 {object} dtos.ErrorResponse "AI request error"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /roadmap/problems/{topic_id}/{theme_id} [get]
func (r RoadmapController) GetProblems(ctx *gin.Context) {
	userId, err := httputils.GetUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := r.userSrv.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	topicId, err := strconv.ParseUint(ctx.Param("topic_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	topic, err := r.roadmapSrv.GetTopic(userId, uint(topicId), true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	themeId, err := strconv.ParseUint(ctx.Param("theme_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	theme, err := r.roadmapSrv.GetTheme(userId, uint(themeId), false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	prompt := r.promptSrv.GetProblems(10, topic.Title, theme.Title, theme, topic.Themes)

	var target services.PromptProblemsResponse

	err = r.aiSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		if errors.Is(err, services.ErrAIRequestFailed) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	problems, err := r.roadmapSrv.CreateProblems(target.Problems, uint(themeId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	ctx.JSON(http.StatusOK, problems)
}

// IncrementUserScoreAndAddAnswer increments the user's score and adds an answer to a problem
// @Summary      Increment user score and add answer
// @Description  Increments the user's score and adds an answer to a problem after verifying the answer with AI
// @Tags         roadmap
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer <JWT token>"
// @Param        request body dtos.IncreaseUserScoreAndAddAnswerRequest true "Answer details"
// @Success      200 {object} services.PromptProblemResponse "Answer verification result"
// @Failure      400 {object} dtos.ValidationErrorResponse "Invalid request body"
// @Failure      401 {object} dtos.ErrorResponse "Unauthorized"
// @Failure      404 {object} dtos.ErrorResponse "Problem not found"
// @Failure      409 {object} dtos.ErrorResponse "AI request error"
// @Failure      500 {object} dtos.ErrorResponse "Internal server error"
// @Router       /roadmap/resolve [post]
func (r RoadmapController) IncrementUserScoreAndAddAnswer(ctx *gin.Context) {
	var req dtos.IncreaseUserScoreAndAddAnswerRequest

	userId, err := httputils.GetUserId(ctx)
	if err != nil {

		ctx.JSON(http.StatusUnauthorized, dtos.UnauthorizedResponse())
		return
	}

	user, err := r.userSrv.GetById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valErr validator.ValidationErrors
		ok := errors.As(err, &valErr)

		if ok {
			ctx.JSON(http.StatusBadRequest, dtos.ValidationErrorResponse{Error: valid.ParseValidationErrors(err)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	problem, err := r.roadmapSrv.GetProblem(req.ProblemId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dtos.NotFoundResponse())
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	prompt := r.promptSrv.VerifyAnswer(problem.Question, req.Answer)
	var target services.PromptProblemResponse

	err = r.aiSrv.GetAnswer(user.ChatGptToken, user.ChatGptModel, []*services.DialogItem{{Text: prompt, IsUser: true}}, &target)
	if err != nil {
		if errors.Is(err, services.ErrAIRequestFailed) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
		return
	}

	if target.Ok {
		err = r.roadmapSrv.IncrementUserScoreAndAddAnswer(userId, req.ProblemId, 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dtos.InternalServerErrorResponse())
			return
		}
	}

	ctx.JSON(http.StatusOK, target)
}

func NewRoadmapController(
	userSrv services.UserService,
	nlSrv services.GptService,
	promptSrv services.PromptService,
	roadmapSrv services.RoadmapService,
) *RoadmapController {
	return &RoadmapController{userSrv, nlSrv, promptSrv, roadmapSrv}
}
