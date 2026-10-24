package controller

import (
	"net/http"
	"strconv"

	dto "visa_assignment/models/dto"
	"visa_assignment/service"
	"visa_assignment/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
		zerolog.Ctx(ctx.Request.Context()).Error().
			Err(err).
			Msg("Create member request rejected: invalid payload")

		utils.RespondWithError(ctx, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	zerolog.Ctx(ctx.Request.Context()).Info().
		Interface("request", req).
		Msg("Create member request accepted")
	member, err := mc.memberService.CreateMember(ctx.Request.Context(), req)
	if err != nil {
		zerolog.Ctx(ctx.Request.Context()).Error().
			Err(err).
			Msg("Failed to create member")
		utils.HandleServiceError(ctx, err)
		return
	}

	zerolog.Ctx(ctx.Request.Context()).Info().
		Int64("member_id", member.MemberID).
		Msg("Member created")
	ctx.JSON(http.StatusCreated, member)
}

// GetMember handles GET /members/:memberId
func (mc *MemberController) GetMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("memberId"), 10, 64)
	if err != nil || memberID <= 0 {
		zerolog.Ctx(c.Request.Context()).Error().
			Str("member_id", c.Param("memberId")).
			Msg("Get member request rejected: invalid member ID")
		utils.RespondWithError(c, http.StatusBadRequest, "invalid member_id parameter")
		return
	}

	member, err := mc.memberService.GetMemberDetails(c.Request.Context(), memberID)
	if err != nil {
		zerolog.Ctx(c.Request.Context()).Error().
			Int64("member_id", memberID).
			Err(err).
			Msg("Failed to get member")
		utils.HandleServiceError(c, err)
		return
	}

	zerolog.Ctx(c.Request.Context()).Info().
		Interface("response", member).
		Int64("member_id", member.MemberID).
		Msg("Member details retrieved")
	c.JSON(http.StatusOK, member)
}

// GetMemberRewards handles GET /members/:memberId/rewards
func (mc *MemberController) GetMemberRewards(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("memberId"), 10, 64)
	if err != nil || memberID <= 0 {
		zerolog.Ctx(c.Request.Context()).Error().
			Str("member_id", c.Param("memberId")).
			Msg("List member rewards request rejected: invalid member ID")
		utils.RespondWithError(c, http.StatusBadRequest, "invalid member_id parameter")
		return
	}

	rewards, err := mc.rewardService.GetRewardsByMemberID(c.Request.Context(), memberID)
	if err != nil {
		zerolog.Ctx(c.Request.Context()).Error().
			Int64("member_id", memberID).
			Err(err).
			Msg("Failed to list member rewards")
		utils.HandleServiceError(c, err)
		return
	}

	zerolog.Ctx(c.Request.Context()).Info().
		Int64("member_id", memberID).
		Interface("rewards", rewards).
		Msg("Member rewards retrieved")
	c.JSON(http.StatusOK, rewards)
}
