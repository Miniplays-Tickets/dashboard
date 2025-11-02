package api

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Miniplays-Tickets/dashboard/botcontext"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/rpc"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/Miniplays-Tickets/dashboard/utils/types"
	"github.com/TicketsBot-cloud/common/premium"
	"github.com/TicketsBot-cloud/database"
	"github.com/TicketsBot-cloud/gdl/objects/interaction"
	"github.com/TicketsBot-cloud/gdl/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type tag struct {
	Id              string             `json:"id" validate:"required,min=1,max=16"`
	UseGuildCommand bool               `json:"use_guild_command"`
	Content         *string            `json:"content" validate:"omitempty,min=1,max=4096"`
	UseEmbed        bool               `json:"use_embed"`
	Embed           *types.CustomEmbed `json:"embed" validate:"omitempty,dive"`
}

var (
	validate          = validator.New()
	slashCommandRegex = regexp.MustCompile(`^[-_a-zA-Z0-9]{1,32}$`)
)

func CreateTag(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	// Max of 200 tags
	count, err := dbclient.Client.Tag.GetTagCount(ctx, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorStr(fmt.Sprintf("Failed to fetch tag from database: %v", err)))
		return
	}

	if count >= 200 {
		ctx.JSON(400, utils.ErrorStr("Du kannst nur bis zu 200 Tag haben"))
		return
	}

	var data tag
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid request data. Please check your input and try again."))
		return
	}

	data.Id = strings.ToLower(data.Id)

	if !data.UseEmbed {
		data.Embed = nil
	}

	// TODO: Limit command amount
	if err := validate.Struct(data); err != nil {
		var validationErrors validator.ValidationErrors
		if ok := errors.As(err, &validationErrors); !ok {
			ctx.JSON(500, utils.ErrorStr("Es ist ein Fehler aufgetreten die Integration zu validieren"))
			return
		}

		formatted := "Deine Eingabe enthält die folgenden Fehler:\n" + utils.FormatValidationErrors(validationErrors)
		ctx.JSON(400, utils.ErrorStr(formatted))
		return
	}

	if !data.verifyId() {
		ctx.JSON(400, utils.ErrorStr("Tag IDs müssen alphabetisch sein (inklusive Bindestriche und Unterstriche) und müssen zwische 1-16 Zeichen lang sein."))
		return
	}

	if !data.verifyContent() {
		ctx.JSON(400, utils.ErrorStr("Du hast keinen Inhalt für den Tag angegeben"))
		return
	}

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorStr("Unable to connect to Discord. Please try again later."))
		return
	}

	if data.UseGuildCommand {
		premiumTier, err := rpc.PremiumClient.GetTierByGuildId(ctx, guildId, true, botContext.Token, botContext.RateLimiter)
		if err != nil {
			ctx.JSON(500, utils.ErrorStr("Unable to verify premium status. Please try again."))
			return
		}

		if premiumTier < premium.Premium {
			ctx.JSON(400, utils.ErrorStr("Premium wird benötigt um eigene Befehle zu nutzen"))
			return
		}
	}

	var embed *database.CustomEmbedWithFields
	if data.Embed != nil {
		customEmbed, fields := data.Embed.IntoDatabaseStruct()
		embed = &database.CustomEmbedWithFields{
			CustomEmbed: customEmbed,
			Fields:      fields,
		}
	}

	var applicationCommandId *uint64
	if data.UseGuildCommand {
		cmd, err := botContext.CreateGuildCommand(ctx, guildId, rest.CreateCommandData{
			Name:        data.Id,
			Description: fmt.Sprintf("Alias für /tag %s", data.Id),
			Options:     nil,
			Type:        interaction.ApplicationCommandTypeChatInput,
		})

		if err != nil {
			ctx.JSON(500, utils.ErrorStr("Failed to create tag. Please try again."))
			return
		}

		applicationCommandId = &cmd.Id
	}

	wrapped := database.Tag{
		Id:                   data.Id,
		GuildId:              guildId,
		Content:              data.Content,
		Embed:                embed,
		ApplicationCommandId: applicationCommandId,
	}

	if err := dbclient.Client.Tag.Set(ctx, wrapped); err != nil {
		ctx.JSON(500, utils.ErrorStr("Failed to create tag. Please try again."))
		return
	}

	ctx.Status(204)
}

func (t *tag) verifyId() bool {
	if len(t.Id) == 0 || len(t.Id) > 16 || strings.Contains(t.Id, " ") {
		return false
	}

	if t.UseGuildCommand {
		return slashCommandRegex.MatchString(t.Id)
	} else {
		return true
	}
}

func (t *tag) verifyContent() bool {
	if t.Content != nil { // validator ensures that if this is not nil, > 0 length
		return true
	}

	if t.Embed != nil {
		if t.Embed.Description != nil || len(t.Embed.Fields) > 0 || t.Embed.ImageUrl != nil || t.Embed.ThumbnailUrl != nil {
			return true
		}
	}

	return false
}
