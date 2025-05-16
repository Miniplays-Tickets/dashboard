package api

import (
	"github.com/Miniplays-Tickets/dashboard/botcontext"
	"github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

type deleteBody struct {
	TagId string `json:"tag_id"`
}

func DeleteTag(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var body deleteBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// Increase max length for characters from other alphabets
	if body.TagId == "" || len(body.TagId) > 100 {
		ctx.JSON(400, utils.ErrorStr("Ungültiger Tag"))
		return
	}

	// Fetch tag to see if we need to delete a guild command
	tag, exists, err := database.Client.Tag.Get(ctx, guildId, body.TagId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !exists {
		ctx.JSON(404, utils.ErrorStr("Tag nicht gefunden"))
		return
	}

	if tag.ApplicationCommandId != nil {
		botContext, err := botcontext.ContextForGuild(guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if err := botContext.DeleteGuildCommand(ctx, guildId, *tag.ApplicationCommandId); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}
	}

	if err := database.Client.Tag.Delete(ctx, guildId, body.TagId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}
