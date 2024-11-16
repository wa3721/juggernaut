package newticket

import (
	"github.com/gin-gonic/gin"
	logmgr "judgement/config/log"
	"net/http"
)

func NewTicketHandler(c *gin.Context) {
	logmgr.Log.Info("success")
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
