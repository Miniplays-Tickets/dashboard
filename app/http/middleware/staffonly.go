package middleware

import (
	"github.com/Miniplays-Tickets/dashboard/config"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func StaffOnly(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	isBotStaff, err := dbclient.Client.BotStaff.IsStaff(ctx, userId)
	if err != nil {
		ctx.JSON(401, utils.ErrorStr("Keine Berechtigungen"))
		ctx.Abort()
		return
	}

	if !isBotStaff {
		ctx.JSON(401, utils.ErrorStr("Keine Berechtigungen"))
		ctx.Abort()
		return
	}

	if !utils.Contains(config.Conf.Admins, userId) {
		ctx.JSON(401, utils.ErrorStr("Keine Berechtigungen"))
		ctx.Abort()
		return
	}
}
