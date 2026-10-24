package controller

import (
	"log"
	"net/http"
	"strconv"

	dto "visa_assignment/models/dto"
	"visa_assignment/service"
	"visa_assignment/utils"

	"github.com/gin-gonic/gin"
)

type MemberController struct {
	memberService service.MemberService
	rewardService service.RewardService
}

func NewMemberController(memberService service.MemberService, rewardService service.RewardService) *MemberController {
	return &MemberController{
		memberService: memberService,
		rewardService: rewardService,
	}
}

// CreateMember handles POST /members
func (mc *MemberController) CreateMember(ctx *gin.Context) {
	var req dto.MemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON: %v", err)

		utils.RespondWithError(ctx, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	member, err := mc.memberService.CreateMember(ctx.Request.Context(), req)
	if err != nil {
		utils.HandleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, member)
}

// GetMember handles GET /members/:memberId
func (mc *MemberController) GetMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("memberId"), 10, 64)
	if err != nil || memberID <= 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid member_id parameter")
		return
	}

	member, err := mc.memberService.GetMemberWithBalance(c.Request.Context(), memberID)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, member)
}

// GetMemberRewards handles GET /members/:memberId/rewards
func (mc *MemberController) GetMemberRewards(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("memberId"), 10, 64)
	if err != nil || memberID <= 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid member_id parameter")
		return
	}

	rewards, err := mc.rewardService.GetRewardsByMemberID(c.Request.Context(), memberID)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, rewards)
}
