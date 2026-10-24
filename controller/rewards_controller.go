package controller

import (
	"net/http"

	"rewards-ledger/model"
	"rewards-ledger/service"
	"rewards-ledger/utils"

	"github.com/gin-gonic/gin"
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
	var req model.CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	reward, err := rc.rewardService.CreateReward(c.Request.Context(), req)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, reward)
}
