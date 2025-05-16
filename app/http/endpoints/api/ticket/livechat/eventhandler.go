package livechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Miniplays-Tickets/dashboard/botcontext"
	"github.com/Miniplays-Tickets/dashboard/config"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/internal/api"
	"github.com/Miniplays-Tickets/dashboard/rpc"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/TicketsBot-cloud/common/premium"
	"github.com/golang-jwt/jwt"
)

func (c *Client) HandleEvent(event Event) error {
	switch event.Type {
	case EventTypeAuth:
		var data AuthData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			c.Write(NewErrorMessage("Fehler 14"))
			_ = c.Ws.Close()
			c.Flush()
			return err
		}

		if err := c.handleAuthEvent(data); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) handleAuthEvent(data AuthData) error {
	if c.Authenticated {
		return api.NewErrorWithMessage(http.StatusBadRequest, errors.New("Bereits Authentifiziert"), "Bereits Authentifiziert")
	}

	token, err := jwt.Parse(data.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Fehler 15: %v", token.Header["alg"])
		}

		return []byte(config.Conf.Server.Secret), nil
	})
	if err != nil {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Fehler 16")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Fehler 17")
	}

	userIdStr, ok := claims["userid"].(string)
	if !ok {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Fehler 18")
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Fehler 19")
	}

	// Get the ticket
	ticket, err := dbclient.Client.Tickets.Get(context.Background(), c.TicketId, c.GuildId)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Fehler beim abrufen der Ticket Daten")
	}

	if ticket.Id == 0 || ticket.GuildId == 0 || ticket.GuildId != c.GuildId {
		return api.NewErrorWithMessage(http.StatusNotFound, err, "Ticket nicht gefunden")
	}

	// Verify the user has permissions to be here
	hasPermission, requestErr := utils.HasPermissionToViewTicket(context.Background(), c.GuildId, userId, ticket)
	if requestErr != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Fehler beim abrufen von Berechtigungen")
	}

	if !hasPermission {
		return api.NewErrorWithMessage(http.StatusForbidden, err, "You do not have permission to view this ticket")
	}

	// Check premium
	botContext, err := botcontext.ContextForGuild(c.GuildId)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Fehler beim abrufen vom Bot context")
	}

	// Verify the guild is premium
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(context.Background(), c.GuildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Fehler beim abrufen vom Premium Tier")
	}

	if premiumTier == premium.None {
		return api.NewErrorWithMessage(http.StatusPaymentRequired, err, "Live-chat benötigt Premium")
	}

	c.Authenticated = true

	c.Write(Event{
		Type: EventTypeAuthenticated,
	})

	return nil
}
