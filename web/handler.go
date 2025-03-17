package web

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/mulan/code"
	"github.com/virzz/mulan/rsp"
)

func ErrCodeHandler(c *gin.Context) { c.JSON(200, rsp.S(code.Codes)) }
