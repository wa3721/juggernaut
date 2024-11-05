package newreply

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewreplyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
