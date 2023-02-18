// Generated by gosql: DO NOT EDIT.
package controllers

import (
	"context"
	"database/sql"
	"http"

	"github.com/gin-gonic/gin"

	"github.com/arjendevos/gosql/models/am"
	"github.com/arjendevos/gosql/models/dm"
)

type ScrapedInstagramAccountChainedController struct {
	*Client
}

func (c *ScrapedInstagramAccountChainedController) List(ctx *gin.Context) {
	queryMods, err := ParseScrapedInstagramAccountChainedListQueryToMods(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responseWithPayload(nil, "invalid_request", "Invalid request", false))
		return
	}

	payload, err := dm.ScrapedInstagramAccountChaineds(queryMods...).All(ctx, c.db)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responseWithPayload(nil, "generic", "Something went wrong", false))
		return
	}

	ctx.JSON(http.StatusOK, responseWithPayload(am.SqlBoilerScrapedInstagramAccountChainedsToApiScrapedInstagramAccountChaineds(payload), nil, nil, true))
}
