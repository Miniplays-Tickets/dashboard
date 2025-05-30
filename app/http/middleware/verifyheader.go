package middleware

import (
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func VerifyXTicketsHeader(ctx *gin.Context) {
	if ctx.GetHeader("x-tickets") != "true" {
		ctx.AbortWithStatusJSON(400, utils.ErrorStr("Fehlender x-tickets header"))
	}
}
