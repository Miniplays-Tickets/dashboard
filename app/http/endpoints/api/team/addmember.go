package api

import (
	"strconv"

	dbclient "github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func AddMember(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	snowflake, err := strconv.ParseUint(ctx.Param("snowflake"), 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// get entity type
	typeParsed, err := strconv.Atoi(ctx.Query("type"))
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	entityType, ok := entityTypes[typeParsed]
	if !ok {
		ctx.JSON(400, utils.ErrorStr("Fehler 23"))
		return
	}

	if entityType == entityTypeUser {
		ctx.JSON(400, utils.ErrorStr("Nur Rollen können zu Teams hinzugefügt werden"))
		return
	}

	if entityType == entityTypeRole && snowflake == guildId {
		ctx.JSON(400, utils.ErrorStr("Du kannst nicht @everyone zu einem Team hinzufügen"))
		return
	}

	teamId := ctx.Param("teamid")
	if teamId == "default" {
		addDefaultMember(ctx, guildId, snowflake, entityType)
	} else {
		parsed, err := strconv.Atoi(teamId)
		if err != nil {
			ctx.JSON(400, utils.ErrorStr("Ungültige Team ID"))
			return
		}

		addTeamMember(ctx, parsed, guildId, snowflake, entityType)
	}
}

func addDefaultMember(ctx *gin.Context, guildId, snowflake uint64, entityType entityType) {
	var err error
	switch entityType {
	case entityTypeUser:
		err = dbclient.Client.Permissions.AddSupport(ctx, guildId, snowflake)
	case entityTypeRole:
		err = dbclient.Client.RolePermissions.AddSupport(ctx, guildId, snowflake)
	}

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

func addTeamMember(ctx *gin.Context, teamId int, guildId, snowflake uint64, entityType entityType) {
	exists, err := dbclient.Client.SupportTeam.Exists(ctx, teamId, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !exists {
		ctx.JSON(404, utils.ErrorStr("Support Team mit der angegebenen ID nicht gefunden"))
		return
	}

	switch entityType {
	case entityTypeUser:
		err = dbclient.Client.SupportTeamMembers.Add(ctx, teamId, snowflake)
	case entityTypeRole:
		err = dbclient.Client.SupportTeamRoles.Add(ctx, teamId, snowflake)
	}

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}
