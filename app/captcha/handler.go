package captcha

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//这里处理业务逻辑

func CaptchaHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
