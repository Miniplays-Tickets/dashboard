package botstaff

import (
	"strconv"

	"github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func AddBotStaffHandler(ctx *gin.Context) {
	userId, err := strconv.ParseUint(ctx.Param("userid"), 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Failed to process request. Please try again."))
		return
	}

	if err := database.Client.BotStaff.Add(ctx, userId); err != nil {
		ctx.JSON(500, utils.ErrorStr("Failed to process request. Please try again."))
		return
	}

	ctx.Status(204)
}
