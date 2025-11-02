package api

import (
	"context"
	"net/http"

	"github.com/Miniplays-Tickets/dashboard/app"
	"github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/TicketsBot-cloud/gdl/objects/user"
	"github.com/TicketsBot-cloud/gdl/rest"
	"github.com/gin-gonic/gin"
)

type whitelabelResponse struct {
	Id       uint64 `json:"id,string"`
	Username string `json:"username"`
	statusUpdateBody
}

func WhitelabelGet(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	// Check if this is a different token
	bot, err := database.Client.Whitelabel.GetByUserId(c, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to load whitelabel bots"))
		return
	}

	if bot.BotId == 0 {
		c.JSON(404, utils.ErrorStr("Kein Bot Gefunden"))
		return
	}

	// Get status
	status, statusType, _, err := database.Client.WhitelabelStatuses.Get(c, bot.BotId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to load whitelabel bots"))
		return
	}

	username := getBotUsername(c, bot.Token)

	c.JSON(200, whitelabelResponse{
		Id:       bot.BotId,
		Username: username,
		statusUpdateBody: statusUpdateBody{ // Zero values if no status is fine
			Status:     status,
			StatusType: user.ActivityType(statusType),
		},
	})
}

func getBotUsername(ctx context.Context, token string) string {
	user, err := rest.GetCurrentUser(ctx, token, nil)
	if err != nil {
		// TODO: Log error
		return "Unbekannter Benutzer"
	}

	return user.Username
}
