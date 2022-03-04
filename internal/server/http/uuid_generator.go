package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func generateUuid(c *gin.Context) {
	u, err := uuid.NewV4()
	if err != nil {
		renderBadRequest(c, err)
	}
	renderOK(c, u.String())
}
