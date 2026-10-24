package controller

import (
	"log"
	"net/http"

	dto "visa_assignment/models/dto"
	"visa_assignment/service"
	"visa_assignment/utils"

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
	var req dto.RewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	reward, err := rc.rewardService.CreateReward(c.Request.Context(), req)
	if err != nil {
		log.Printf("Failed to create reward: %v", err)
		utils.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, reward)
}
