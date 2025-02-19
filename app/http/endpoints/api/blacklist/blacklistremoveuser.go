package api

import (
	"strconv"

	"github.com/Dev-Miniplays/Ticketsv2-dashboard/database"
	"github.com/Dev-Miniplays/Ticketsv2-dashboard/utils"
	"github.com/gin-gonic/gin"
)

func RemoveUserBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	userId, err := strconv.ParseUint(ctx.Param("user"), 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if err := database.Client.Blacklist.Remove(ctx, guildId, userId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}
