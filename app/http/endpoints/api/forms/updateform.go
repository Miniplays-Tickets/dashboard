package forms

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Miniplays-Tickets/dashboard/app"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func UpdateForm(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

	var data createFormBody
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, utils.ErrorStr("Invalid request data. Please check your input and try again."))
		return
	}

	// Validate title is not empty or whitespace-only
	if len(strings.TrimSpace(data.Title)) == 0 {
		c.JSON(400, utils.ErrorStr("Form title cannot be empty"))
		return
	}

	if len(data.Title) > 45 {
		c.JSON(400, utils.ErrorStr("Form title must be 45 characters or less (current: %d characters)", len(data.Title)))
		return
	}

	formId, err := strconv.Atoi(c.Param("form_id"))
	if err != nil {
		c.JSON(400, utils.ErrorStr("Ungültige Formular ID: %s", c.Param("form_id")))
		return
	}

	form, ok, err := dbclient.Client.Forms.Get(c, formId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to fetch form from database"))
		return
	}

	if !ok {
		c.JSON(404, utils.ErrorStr("Formular #%d nicht gefunden", formId))
		return
	}

	if form.GuildId != guildId {
		c.JSON(403, utils.ErrorStr("Formular #%d gehört nicht zu dieser Guild %d", formId, guildId))
		return
	}

	if err := dbclient.Client.Forms.UpdateTitle(c, formId, data.Title); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to update form title in database"))
		return
	}

	c.JSON(200, utils.SuccessResponse)
}
