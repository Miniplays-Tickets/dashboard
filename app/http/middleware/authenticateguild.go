package middleware

import (
	"context"
	"errors"
	"strconv"

	"github.com/Dev-Miniplays/Ticketsv2-dashboard/rpc/cache"
	"github.com/Dev-Miniplays/Ticketsv2-dashboard/utils"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/gin-gonic/gin"
	cache2 "github.com/rxdn/gdl/cache"
)

// requires AuthenticateCookie middleware to be run before
func AuthenticateGuild(requiredPermissionLevel permission.PermissionLevel) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if guildId, ok := ctx.Params.Get("id"); ok {
			parsed, err := strconv.ParseUint(guildId, 10, 64)
			if err != nil {
				ctx.JSON(400, utils.ErrorStr("Ungültige Guild ID"))
				ctx.Abort()
				return
			}

			ctx.Keys["guildid"] = parsed

			// TODO: Do we need this? Only really serves as a check whether the bot is in the server
			// TODO: Use proper context
			if _, err := cache.Instance.GetGuildOwner(context.Background(), parsed); err != nil {
				if errors.Is(err, cache2.ErrNotFound) {
					ctx.JSON(404, utils.ErrorStr("Guild not found"))
					ctx.Abort()
				} else {
					ctx.JSON(500, utils.ErrorJson(err))
					ctx.Abort()
				}

				return
			}

			// Verify the user has permissions to be here
			userId := ctx.Keys["userid"].(uint64)

			// TODO: Use proper context
			permLevel, err := utils.GetPermissionLevel(context.Background(), parsed, userId)
			if err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				ctx.Abort()
				return
			}

			if permLevel < requiredPermissionLevel {
				ctx.JSON(403, utils.ErrorStr("Keine Berechtigungen"))
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(400, utils.ErrorStr("Ungültige Guild ID"))
			ctx.Abort()
			return
		}
	}
}
