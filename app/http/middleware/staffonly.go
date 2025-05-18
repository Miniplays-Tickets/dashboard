package middleware

import (
	"github.com/Miniplays-Tickets/dashboard/config"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func StaffOnly(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	isBotStaff, _ := dbclient.Client.BotStaff.IsStaff(ctx, userId)

	if isBotStaff || utils.Contains(config.Conf.Admins, userId) {
		return
	}

	ctx.JSON(401, utils.ErrorStr("Keine Berechtigungen"))
	ctx.Abort()
	return
}
