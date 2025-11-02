package forms

import (
	"net/http"
	"strings"

	"github.com/Miniplays-Tickets/dashboard/app"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/TicketsBot-cloud/database"
	"github.com/gin-gonic/gin"
)

type createFormBody struct {
	Title string `json:"title"`
}

func CreateForm(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

	var data createFormBody
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, utils.ErrorStr("Fehler 29"))
		return
	}

	// Validate title is not empty or whitespace-only
	if len(strings.TrimSpace(data.Title)) == 0 {
		c.JSON(400, utils.ErrorStr("Form title cannot be empty"))
		return
	}

	if len(data.Title) > 45 {
		c.JSON(400, utils.ErrorStr("Form Titel muss 45 oder weniger Zeichen haben (aktuell: %d Zeichen)", len(data.Title)))
		return
	}

	// 26^50 chance of collision
	customId, err := utils.RandString(30)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to generate unique form ID"))
		return
	}

	id, err := dbclient.Client.Forms.Create(c, guildId, data.Title, customId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to create form in database"))
		return
	}

	form := database.Form{
		Id:       id,
		GuildId:  guildId,
		Title:    data.Title,
		CustomId: customId,
	}

	c.JSON(200, form)
}
