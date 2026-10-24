package controller

import (
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
		utils.LogRequest(c.Request.Context(), "create_reward rejected invalid_payload error=%q", err)
		utils.RespondWithError(c, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	utils.LogRequest(
		c.Request.Context(),
		"create_reward accepted : %+v", req)
	reward, err := rc.rewardService.CreateReward(c.Request.Context(), req)
	if err != nil {
		utils.LogRequest(c.Request.Context(), "create_reward failed member_id=%d error=%q", req.MemberID, err)
		utils.HandleServiceError(c, err)
		return
	}

	utils.LogRequest(c.Request.Context(), "create_reward response : %+v", reward)
	c.JSON(http.StatusCreated, reward)
}
