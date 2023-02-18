// Generated by gosql: DO NOT EDIT.
package controllers

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
)

func responseWithPayload(payload interface{}, errCode, message interface{}, ok bool) gin.H {
	var m interface{}
	var e interface{}

	if errCode != nil {
		e = errCode
	}

	if message != nil {
		m = message
	}

	return gin.H{"payload": payload, "error": e, "message": m, "ok": ok}
}

func bindAndValidateJSON(context *gin.Context, obj interface{}) error {
	if err := context.ShouldBindJSON(obj); err != nil {
		return err
	}

	if err := validator.Validate(obj); err != nil {
		return err
	}

	return nil
}
