package remind

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RemindHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
