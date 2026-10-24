package controller

import (
	"net/http"

	dto "visa_assignment/models/dto"
	"visa_assignment/service"
	"visa_assignment/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type RewardController struct {
	rewardService service.RewardService
}

func NewRewardController(rewardService service.RewardService) *RewardController {
	return &RewardController{
		rewardService: rewardService,
	}
}

// CreateReward handles POST /rewards
func (rc *RewardController) CreateReward(c *gin.Context) {
	var req dto.RewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(c.Request.Context()).Error().
			Err(err).
			Msg("Create reward request rejected: invalid payload")
		utils.RespondWithError(c, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	zerolog.Ctx(c.Request.Context()).Info().
		Interface("request", req).
		Msg("Create reward request accepted")

	reward, err := rc.rewardService.CreateReward(c.Request.Context(), req)
	if err != nil {
		zerolog.Ctx(c.Request.Context()).Error().
			Err(err).
			Msg("Failed to create reward")
		utils.HandleServiceError(c, err)
		return
	}

	zerolog.Ctx(c.Request.Context()).Info().
		Interface("reward_id", reward.RewardID).
		Msg("Reward entry recorded")
	c.JSON(http.StatusCreated, reward)
}
