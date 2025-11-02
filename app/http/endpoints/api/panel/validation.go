package api

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Miniplays-Tickets/dashboard/app"
	"github.com/Miniplays-Tickets/dashboard/app/http/validation"
	"github.com/Miniplays-Tickets/dashboard/app/http/validation/defaults"
	"github.com/Miniplays-Tickets/dashboard/botcontext"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/Miniplays-Tickets/dashboard/utils/types"
	"github.com/TicketsBot-cloud/database"
	"github.com/TicketsBot-cloud/gdl/objects/channel"
	"github.com/TicketsBot-cloud/gdl/objects/guild"
	"github.com/TicketsBot-cloud/gdl/objects/interaction/component"
)

func ApplyPanelDefaults(data *panelBody) {
	for _, applicator := range DefaultApplicators(data) {
		if applicator.ShouldApply() {
			applicator.Apply()
		}
	}
}

func DefaultApplicators(data *panelBody) []defaults.DefaultApplicator {
	return []defaults.DefaultApplicator{
		defaults.NewDefaultApplicator(defaults.EmptyStringCheck, &data.Title, "Öffne ein Ticket!"),
		defaults.NewDefaultApplicator(defaults.EmptyStringCheck, &data.Content, "Durch Klicken auf die Schaltfläche wird ein Ticket für dich eröffnet"),
		defaults.NewDefaultApplicator[*string](defaults.NilOrEmptyStringCheck, &data.ImageUrl, nil),
		defaults.NewDefaultApplicator[*string](defaults.NilOrEmptyStringCheck, &data.ThumbnailUrl, nil),
		defaults.NewDefaultApplicator(defaults.EmptyStringCheck, &data.ButtonLabel, data.Title),
		defaults.NewDefaultApplicator(defaults.EmptyStringCheck, &data.ButtonLabel, "Öffne ein Ticket!"), // Title could have been blank
		defaults.NewDefaultApplicator[*string](defaults.NilOrEmptyStringCheck, &data.NamingScheme, nil),
	}
}

type PanelValidationContext struct {
	Data       panelBody
	GuildId    uint64
	IsPremium  bool
	BotContext *botcontext.BotContext
	Channels   []channel.Channel
	Roles      []guild.Role
}

func ValidatePanelBody(validationContext PanelValidationContext) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	return validation.Validate(ctx, validationContext, panelValidators()...)
}

func panelValidators() []validation.Validator[PanelValidationContext] {
	return []validation.Validator[PanelValidationContext]{
		validateTitle,
		validateContent,
		validateChannelId,
		validateCategory,
		validateEmoji,
		validateImageUrl,
		validateThumbnailUrl,
		validateButtonStyle,
		validateButtonLabel,
		validateFormId,
		validateExitSurveyFormId,
		validateTeams,
		validateNamingScheme,
		validateWelcomeMessage,
		validateAccessControlList,
		validatePendingCategory,
	}
}

func validateTitle(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if len(ctx.Data.Title) > 80 {
			return validation.NewInvalidInputError("Der Titel des Panels darf nicht länger als 80 Zeichen sein")
		}

		return nil
	}
}

func validateContent(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if len(ctx.Data.Content) > 4096 {
			return validation.NewInvalidInputError("Der Inhalt des Panels darf nicht länger als 4096 Zeichen sein")
		}

		return nil
	}
}

func validateChannelId(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		for _, ch := range ctx.Channels {
			if ch.Id == ctx.Data.ChannelId && (ch.Type == channel.ChannelTypeGuildText || ch.Type == channel.ChannelTypeGuildNews) {
				return nil
			}
		}

		return validation.NewInvalidInputError("Panel Kanal nicht gefunden")
	}
}

func validateCategory(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		for _, ch := range ctx.Channels {
			if ch.Id == ctx.Data.CategoryId && ch.Type == channel.ChannelTypeGuildCategory {
				return nil
			}
		}

		return validation.NewInvalidInputError("Invalid ticket category")
	}
}

func validateEmoji(c PanelValidationContext) validation.ValidationFunc {
	return func() error {
		emoji := c.Data.Emoji

		if emoji.IsCustomEmoji {
			if emoji.Id == nil {
				return validation.NewInvalidInputError("Eingenes emoji hat eine fehlende ID")
			}

			ctx, cancel := context.WithTimeout(context.Background(), app.DefaultTimeout)
			defer cancel()

			resolvedEmoji, err := c.BotContext.GetGuildEmoji(ctx, c.GuildId, *emoji.Id)
			if err != nil {
				return err
			}

			if resolvedEmoji.Id.Value == 0 {
				return validation.NewInvalidInputError("Emoji nicht gefunden")
			}

			if resolvedEmoji.Name != emoji.Name {
				return validation.NewInvalidInputError("Emoji Name stimmt nicht")
			}
		} else {
			if len(emoji.Name) == 0 {
				return validation.NewInvalidInputError("Emoji Name ist leer")
			}

			// Convert from :emoji: to unicode if we need to
			name := strings.TrimSpace(emoji.Name)
			name = strings.Replace(name, ":", "", -1)

			unicode, ok := utils.GetEmoji(name)
			if !ok {
				return validation.NewInvalidInputError("Ungültiger emoji")
			}

			emoji.Name = unicode
		}

		return nil
	}
}

var urlRegex = regexp.MustCompile(`^https?://([-a-zA-Z0-9@:%._+~#=]{1,256})\.[a-zA-Z0-9()]{1,63}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)$`)

func validateNullableUrl(url *string) validation.ValidationFunc {
	return func() error {
		if url != nil && (len(*url) > 255 || !urlRegex.MatchString(*url)) {
			return validation.NewInvalidInputError("Ungültige Bild-URL. Muss mit .gif, .jpg, .jpeg, .png oder .webp enden")
		}

		return nil
	}
}

func validateImageUrl(ctx PanelValidationContext) validation.ValidationFunc {
	return validateNullableUrl(ctx.Data.ImageUrl)
}

func validateThumbnailUrl(ctx PanelValidationContext) validation.ValidationFunc {
	return validateNullableUrl(ctx.Data.ThumbnailUrl)
}

func validateButtonStyle(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if ctx.Data.ButtonStyle < component.ButtonStylePrimary && ctx.Data.ButtonStyle > component.ButtonStyleDanger {
			return validation.NewInvalidInputError("Ungültiger Button Style")
		}

		return nil
	}
}

func validateButtonLabel(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if len(ctx.Data.ButtonLabel) > 80 {
			return validation.NewInvalidInputError("Die Beschriftung des Buttons darf nicht länger als 80 Zeichen sein")
		}

		return nil
	}
}

func validatedNullableFormId(guildId uint64, formId *int) validation.ValidationFunc {
	return func() error {
		if formId == nil {
			return nil
		}

		form, ok, err := dbclient.Client.Forms.Get(context.Background(), *formId)
		if err != nil {
			return err
		}

		if !ok {
			return validation.NewInvalidInputError("Formular nicht gefunden")
		}

		if form.GuildId != guildId {
			return validation.NewInvalidInputError("Guild ID stimmt beim Validieren des Formulars nicht überein")
		}

		return nil
	}
}

func validateFormId(ctx PanelValidationContext) validation.ValidationFunc {
	return validatedNullableFormId(ctx.GuildId, ctx.Data.FormId)
}

// Check premium on the worker side to maintain settings if user unsubscribes and later resubscribes
func validateExitSurveyFormId(ctx PanelValidationContext) validation.ValidationFunc {
	return validatedNullableFormId(ctx.GuildId, ctx.Data.ExitSurveyFormId)
}

func validatePendingCategory(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if ctx.Data.PendingCategory == nil {
			return nil
		}

		if !ctx.IsPremium {
			return validation.NewInvalidInputError("Die Kategorie „Auf Antwort warten“ ist eine Premium-Funktion")
		}

		for _, ch := range ctx.Channels {
			if ch.Id == *ctx.Data.PendingCategory && ch.Type == channel.ChannelTypeGuildCategory {
				return nil
			}
		}

		return validation.NewInvalidInputError("Ungültige Kategorie „Auf Antwort warten“")
	}
}

func validateTeams(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		// Query does not work nicely if there are no teams created in the guild, but if the user submits no teams,
		// then the input is guaranteed to be valid. Teams array excludes default team.
		if len(ctx.Data.Teams) == 0 {
			return nil
		}

		ok, err := dbclient.Client.SupportTeam.AllTeamsExistForGuild(context.Background(), ctx.GuildId, ctx.Data.Teams)
		if err != nil {
			return err
		}

		if !ok {
			return validation.NewInvalidInputError("Ungültiges Team")
		}

		return nil
	}
}

var placeholderPattern = regexp.MustCompile(`%(\w+)%`)

// Discord filters out illegal characters (such as +, $, ") when creating the channel for us
func validateNamingScheme(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if ctx.Data.NamingScheme == nil {
			return nil
		}

		if len(*ctx.Data.NamingScheme) > 100 {
			return validation.NewInvalidInputError("Das Namensschema darf nicht länger als 100 Zeichen sein")
		}

		// Validate placeholders used
		validPlaceholders := []string{"id", "username", "nickname", "id_padded"}
		for _, match := range placeholderPattern.FindAllStringSubmatch(*ctx.Data.NamingScheme, -1) {
			if len(match) < 2 { // Infallible
				return errors.New("Fehler 26")
			}

			placeholder := match[1]
			if !utils.Contains(validPlaceholders, placeholder) {
				return validation.NewInvalidInputError(fmt.Sprintf("Ungültiger Platzhalter im Namensschema: %s", placeholder))
			}
		}

		return nil
	}
}

func validateWelcomeMessage(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		return validateEmbed(ctx.Data.WelcomeMessage)
	}
}

func validateAccessControlList(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		acl := ctx.Data.AccessControlList

		if len(acl) == 0 {
			return validation.NewInvalidInputError("Die Zugriffskontrollliste ist leer")
		}

		if len(acl) > 10 {
			return validation.NewInvalidInputError("Die Zugriffskontrollliste kann nicht mehr als 10 Rollen haben")
		}

		roles := utils.ToSet(utils.Map(ctx.Roles, utils.RoleToId))

		if roles.Size() != len(ctx.Roles) {
			return validation.NewInvalidInputError("Doppelte Rollen in der Zugriffskontrollliste")
		}

		everyoneRoleFound := false
		for _, rule := range acl {
			if rule.RoleId == ctx.GuildId {
				everyoneRoleFound = true
			}

			if rule.Action != database.AccessControlActionDeny && rule.Action != database.AccessControlActionAllow {
				return validation.NewInvalidInputErrorf("Ungültige Zugriffskontrollaktion \"%s\"", rule.Action)
			}

			if !roles.Contains(rule.RoleId) {
				return validation.NewInvalidInputErrorf("Ungültige Rolle %d in der Zugriffskontrollliste, nicht in der Guild gefunden", rule.RoleId)
			}
		}

		if !everyoneRoleFound {
			return validation.NewInvalidInputError("Zugriffskontrollliste enthält nicht @everyone")
		}

		return nil
	}
}

func validateEmbed(e *types.CustomEmbed) error {
	if e == nil || e.Title != nil || e.Description != nil || len(e.Fields) > 0 || e.ImageUrl != nil || e.ThumbnailUrl != nil {
		if e.ImageUrl != nil && (len(*e.ImageUrl) > 255 || !urlRegex.MatchString(*e.ImageUrl)) {
			if *e.ImageUrl == "%avatar_url%" {
				// Ignore validation as it is a placeholder
				return nil
			}

			return validation.NewInvalidInputError("Invalid URL")
		}

		if e.ThumbnailUrl != nil && (len(*e.ThumbnailUrl) > 255 || !urlRegex.MatchString(*e.ThumbnailUrl)) {
			if *e.ThumbnailUrl == "%avatar_url%" {
				// Ignore validation as it is a placeholder
				return nil
			}

			return validation.NewInvalidInputError("Invalid URL")
		}

		return nil
	}

	return validation.NewInvalidInputError("Deine eingebettete Nachricht enthält keinen Inhalt")
}
