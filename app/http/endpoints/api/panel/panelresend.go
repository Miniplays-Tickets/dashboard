package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Miniplays-Tickets/dashboard/botcontext"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/rpc"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/TicketsBot-cloud/common/premium"
	"github.com/TicketsBot-cloud/gdl/rest"
	"github.com/TicketsBot-cloud/gdl/rest/request"
	"github.com/gin-gonic/gin"
)

func ResendPanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorStr("Unable to connect to Discord. Please try again later."))
		return
	}

	panelId, err := strconv.Atoi(ctx.Param("panelid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Failed to send message. Please try again."))
		return
	}

	// get existing
	panel, err := dbclient.Client.Panel.GetById(ctx, panelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorStr("Unable to load panel. Please try again."))
		return
	}

	if panel.PanelId == 0 {
		ctx.JSON(404, utils.ErrorStr(fmt.Sprintf("Panel nicht gefunden: %d", panelId)))
		return
	}

	// check guild ID matches
	if panel.GuildId != guildId {
		ctx.JSON(403, utils.ErrorStr("Guild ID stimmt nicht"))
		return
	}

	if panel.ForceDisabled {
		ctx.JSON(400, utils.ErrorStr("Dieses Panel ist Deaktivert und kann nicht bearbeitet werden: Reaktiviere Premium um dieses Panel wieder zu aktivieren"))
		return
	}

	// delete old message
	// TODO: Use proper context
	if err := rest.DeleteMessage(context.Background(), botContext.Token, botContext.RateLimiter, panel.ChannelId, panel.GuildId); err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && !unwrapped.IsClientError() {
			ctx.JSON(500, utils.ErrorStr("Failed to send message. Please try again."))
			return
		}
	}

	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(ctx, guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorStr("Unable to verify premium status. Please try again."))
		return
	}

	messageData := panelIntoMessageData(panel, premiumTier > premium.None)
	msgId, err := messageData.send(botContext)
	if err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && unwrapped.StatusCode == 403 {
			ctx.JSON(500, utils.ErrorStr("Ich habe keine Berechtigung, Nachrichten in dem angegebenen Kanal zu senden"))
		} else {
			ctx.JSON(500, utils.ErrorStr("Failed to send message. Please try again."))
		}

		return
	}

	if err = dbclient.Client.Panel.UpdateMessageId(ctx, panel.PanelId, msgId); err != nil {
		ctx.JSON(500, utils.ErrorStr("Failed to send message. Please try again."))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}
