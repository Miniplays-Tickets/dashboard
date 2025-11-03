package forms

import (
	"net/http"
	"strconv"

	"github.com/Miniplays-Tickets/dashboard/app"
	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func DeleteForm(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

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

	if err := dbclient.Client.Forms.Delete(c, formId); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "Failed to delete form from database"))
		return
	}

	c.JSON(200, utils.SuccessResponse)
}
