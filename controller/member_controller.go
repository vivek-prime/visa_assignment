package controller

import (
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
		utils.LogRequest(ctx.Request.Context(), "create_member rejected invalid_payload error=%q", err)

		utils.RespondWithError(ctx, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	utils.LogRequest(ctx.Request.Context(), "create_member accepted")
	member, err := mc.memberService.CreateMember(ctx.Request.Context(), req)
	if err != nil {
		utils.LogRequest(ctx.Request.Context(), "create_member failed error=%q", err)
		utils.HandleServiceError(ctx, err)
		return
	}

	utils.LogRequest(ctx.Request.Context(), "create_member response member_id=%d", member.MemberID)
	ctx.JSON(http.StatusCreated, member)
}

// GetMember handles GET /members/:memberId
func (mc *MemberController) GetMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("memberId"), 10, 64)
	if err != nil || memberID <= 0 {
		utils.LogRequest(c.Request.Context(), "get_member rejected invalid_member_id value=%q", c.Param("memberId"))
		utils.RespondWithError(c, http.StatusBadRequest, "invalid member_id parameter")
		return
	}

	utils.LogRequest(c.Request.Context(), "get_member member_id=%d", memberID)
	member, err := mc.memberService.GetMemberWithBalance(c.Request.Context(), memberID)
	if err != nil {
		utils.LogRequest(c.Request.Context(), "get_member failed member_id=%d error=%q", memberID, err)
		utils.HandleServiceError(c, err)
		return
	}

	utils.LogRequest(c.Request.Context(), "get_member response : %+v", member)
	c.JSON(http.StatusOK, member)
}

// GetMemberRewards handles GET /members/:memberId/rewards
func (mc *MemberController) GetMemberRewards(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("memberId"), 10, 64)
	if err != nil || memberID <= 0 {
		utils.LogRequest(c.Request.Context(), "get_member_rewards rejected invalid_member_id value=%q", c.Param("memberId"))
		utils.RespondWithError(c, http.StatusBadRequest, "invalid member_id parameter")
		return
	}

	utils.LogRequest(c.Request.Context(), "get_member_rewards , member_id=%d", memberID)
	rewards, err := mc.rewardService.GetRewardsByMemberID(c.Request.Context(), memberID)
	if err != nil {
		utils.LogRequest(c.Request.Context(), "get_member_rewards failed member_id=%d error=%q", memberID, err)
		utils.HandleServiceError(c, err)
		return
	}

	utils.LogRequest(c.Request.Context(), "get_member_rewards response : %+v", rewards)
	c.JSON(http.StatusOK, rewards)
}
