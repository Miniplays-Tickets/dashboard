package api

import (
	"context"

	"github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/rpc/cache"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	cache2 "github.com/rxdn/gdl/cache"
)

type (
	blacklistAddResponse struct {
		Success  bool   `json:"success"`
		Resolved bool   `json:"resolved"`
		Id       uint64 `json:"id,string"`
		Username string `json:"username"`
	}

	blacklistAddBody struct {
		EntityType entityType `json:"entity_type"`
		Snowflake  uint64     `json:"snowflake,string"`
	}

	entityType int
)

const (
	entityTypeUser entityType = iota
	entityTypeRole
)

func AddBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var body blacklistAddBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if body.EntityType == entityTypeUser {
		// Max of 250 blacklisted users
		count, err := database.Client.Blacklist.GetBlacklistedCount(ctx, guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if count >= 250 {
			ctx.JSON(400, utils.ErrorStr("Blacklist Limit von 250 erreicht: Ziehe in Betracht, stattdessen eine Rolle zu verwenden"))
			return
		}

		// TODO: Use proper context
		permLevel, err := utils.GetPermissionLevel(context.Background(), guildId, body.Snowflake)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if permLevel > permission.Everyone {
			ctx.JSON(400, utils.ErrorStr("Du kannst keine Mitglieder eines Teams auf die schwarze Liste setzen!"))
			return
		}

		if err := database.Client.Blacklist.Add(ctx, guildId, body.Snowflake); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// Resolve user
		// TODO: Use proper context
		user, err := cache.Instance.GetUser(context.Background(), body.Snowflake)
		if err != nil {
			if errors.Is(err, cache2.ErrNotFound) {
				ctx.JSON(200, blacklistAddResponse{
					Success:  true,
					Resolved: false,
					Id:       body.Snowflake,
				})
				return
			} else {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		}

		ctx.JSON(200, blacklistAddResponse{
			Success:  true,
			Resolved: true,
			Id:       body.Snowflake,
			Username: user.Username,
		})
	} else if body.EntityType == entityTypeRole {
		// Max of 50 blacklisted roles
		count, err := database.Client.RoleBlacklist.GetBlacklistedCount(ctx, guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if count >= 50 {
			ctx.JSON(400, utils.ErrorStr("Blacklist Limit von 50 erreicht"))
			return
		}

		if err := database.Client.RoleBlacklist.Add(ctx, guildId, body.Snowflake); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		ctx.JSON(200, blacklistAddResponse{
			Success: true,
			Id:      body.Snowflake,
		})
	} else {
		ctx.JSON(400, utils.ErrorStr("Fehler 30"))
		return
	}
}
