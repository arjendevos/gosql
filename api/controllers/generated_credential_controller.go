// Generated by gosql: DO NOT EDIT.
package controllers

import (
	"context"
	"database/sql"
	"http"

	"github.com/gin-gonic/gin"

	"github.com/arjendevos/gosql/models/dm"
	"github.com/arjendevos/gosqlmodels/am"
)

type CredentialController struct {
	*Client
}

func (c *CredentialController) List(ctx *gin.Context) {
	queryMods, err := ParseCredentialListQueryToMods(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responseWithPayload(nil, "invalid_request", "Invalid request", false))
		return
	}

	payload, err := dm.Credentials(queryMods...).All(ctx, c.db)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responseWithPayload(nil, "generic", "Something went wrong", false))
		return
	}

	ctx.JSON(http.StatusOK, responseWithPayload(am.SqlBoilerCredentialsToApiCredentials(payload), nil, nil, true))
}
