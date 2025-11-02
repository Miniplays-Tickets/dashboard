package api

import (
	"strconv"

	"github.com/Miniplays-Tickets/dashboard/botcontext"
	"github.com/Miniplays-Tickets/dashboard/rpc"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/TicketsBot-cloud/common/premium"
	"github.com/gin-gonic/gin"
)

func PremiumHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorStr("Unable to connect to Discord. Please try again later."))
		return
	}

	// If error, will default to false
	includeVoting, _ := strconv.ParseBool(ctx.Query("include_voting"))

	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(ctx, guildId, includeVoting, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorStr("Unable to verify premium status. Please try again."))
		return
	}

	ctx.JSON(200, gin.H{
		"premium": premiumTier >= premium.Premium,
		"tier":    premiumTier,
	})
}
