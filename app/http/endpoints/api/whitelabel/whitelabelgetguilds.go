package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Miniplays-Tickets/dashboard/app"
	"github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/rpc/cache"
	"github.com/Miniplays-Tickets/dashboard/utils"
	cache2 "github.com/TicketsBot-cloud/gdl/cache"
	"github.com/gin-gonic/gin"
)

func WhitelabelGetGuilds(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	bot, err := database.Client.Whitelabel.GetByUserId(c, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to load whitelabel bots"))
		return
	}

	// id -> name
	if bot.BotId == 0 {
		c.JSON(400, utils.ErrorStr("Whitelabel Bot nicht gefunden"))
		return
	}

	ids, err := database.Client.WhitelabelGuilds.GetGuilds(c, bot.BotId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to load whitelabel bots"))
		return
	}

	guilds := make(map[string]string)
	for i, id := range ids {
		if i >= 10 {
			idStr := strconv.FormatUint(id, 10)
			guilds[idStr] = idStr
			continue
		}

		// get guild name
		// TODO: Use proper context
		guild, err := cache.Instance.GetGuild(context.Background(), id)
		if err != nil {
			if errors.Is(err, cache2.ErrNotFound) {
				continue
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to load whitelabel bots"))
				return
			}
		}

		guilds[strconv.FormatUint(id, 10)] = guild.Name
	}

	c.JSON(200, gin.H{
		"success": true,
		"guilds":  guilds,
	})
}
