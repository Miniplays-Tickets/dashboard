package middleware

import (
	"github.com/Dev-Miniplays/Ticketsv2-dashboard/config"
	"github.com/Dev-Miniplays/Ticketsv2-dashboard/utils"
	"github.com/gin-gonic/gin"
)

func AdminOnly(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	if !utils.Contains(config.Conf.Admins, userId) {
		ctx.JSON(401, utils.ErrorStr("Keine Berechtigungen"))
		ctx.Abort()
		return
	}
}
