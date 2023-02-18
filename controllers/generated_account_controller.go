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

type AccountController struct {
	*Client
}

func (c *AccountController) List(ctx *gin.Context) {
	queryMods, err := ParseAccountListQueryToMods(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responseWithPayload(nil, "invalid_request", "Invalid request", false))
		return
	}

	payload, err := dm.Accounts(queryMods...).All(ctx, c.db)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responseWithPayload(nil, "generic", "Something went wrong", false))
		return
	}

	ctx.JSON(http.StatusOK, responseWithPayload(am.SqlBoilerAccountsToApiAccounts(payload), nil, nil, true))
}