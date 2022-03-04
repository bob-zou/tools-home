package http

import (
	"net/http"
	"tools-home/internal/server/http/common"

	"github.com/gin-gonic/gin"
)

func currentUser(c *gin.Context) {
	c.JSON(http.StatusOK, common.Reply{
		Data: map[string]interface{}{"name": "游客", "username": "guest"},
	})
}
