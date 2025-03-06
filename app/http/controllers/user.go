package controllers

import "educ-gpt/services"

type FavoriteController struct {
	userSrv    services.UserService
	roadmapSrv services.RoadmapService
	nlSrv      services.NaturalLanguageService
	promptSrv  services.PromptService
}

func (f FavoriteController) GetFavorites() {

}

func (f FavoriteController) AddFavorite() {

}

func (f FavoriteController) GetRandom() {

}
