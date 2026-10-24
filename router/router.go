package router

import (
	"net/http"
	"visa_assignment/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	memberController *controller.MemberController,
	rewardController *controller.RewardController,
) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// v1 API group
	v1 := r.Group("/v1")
	{
		// Member routes
		v1.POST("/members", memberController.CreateMember)
		v1.GET("/members/:memberId", memberController.GetMember)
		v1.GET("/members/:memberId/rewards", memberController.GetMemberRewards)

		// Reward routes
		v1.POST("/rewards", rewardController.CreateReward)
	}

	return r
}
