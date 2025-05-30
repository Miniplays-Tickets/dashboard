package middleware

import (
	"context"
	"time"

	"github.com/Miniplays-Tickets/dashboard/database"
	"github.com/Miniplays-Tickets/dashboard/utils"
	"github.com/gin-gonic/gin"
)

// UpdateLastSeen We store the last time a user was seen in the dashboard so that we can delete their data if they
// haven't logged in for 30 days.
func UpdateLastSeen(req *gin.Context) {
	userId, ok := req.Keys["userid"].(uint64) // ok=false if not present
	if !ok {
		req.AbortWithStatusJSON(500, utils.ErrorStr("Fehler 2"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1500)
	defer cancel()

	if err := database.Client.DashboardUsers.UpdateLastSeen(ctx, userId); err != nil {
		req.AbortWithStatusJSON(500, utils.ErrorStr(err.Error()))
		return
	}

	req.Next()
}
