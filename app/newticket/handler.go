package newticket

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewTicketHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
