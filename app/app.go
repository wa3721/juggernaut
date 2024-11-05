package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"judgement/app/captcha"
	"judgement/app/newreply"
	"judgement/app/newticket"
	"judgement/app/remind"
	logmgr "judgement/log"
	"time"
)

const (
	verifycode          = "验证码"
	allocationreminders = "工单分配提醒"
	newticketreminder   = "新工单提醒"
	newReplyreminder    = "新回复提醒"
)

//NEW judgement 实例

func NewApp(app string) *gin.Engine {
	fmt.Println(`
   |\  \|\  \|\  \|\   ___ \|\   ____\|\  ___ \ |\   _ \  _   \|\  ___ \ |\   ___  \|\___   ___\ 
   \ \  \ \  \\\  \ \  \_|\ \ \  \___|\ \   __/|\ \  \\\__\ \  \ \   __/|\ \  \\ \  \|___ \  \_| 
 __ \ \  \ \  \\\  \ \  \ \\ \ \  \  __\ \  \_|/_\ \  \\|__| \  \ \  \_|/_\ \  \\ \  \   \ \  \  
|\  \\_\  \ \  \\\  \ \  \_\\ \ \  \|\  \ \  \_|\ \ \  \    \ \  \ \  \_|\ \ \  \\ \  \   \ \  \ 
\ \________\ \_______\ \_______\ \_______\ \_______\ \__\    \ \__\ \_______\ \__\\ \__\   \ \__\
 \|________|\|_______|\|_______|\|_______|\|_______|\|__|     \|__|\|_______|\|__| \|__|    \|__|
	`)
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	backend := r.Group("/" + app)
	{ //验证码转发路由
		backend.POST("/captcha", logMiddleware(verifycode), captcha.CaptchaHandler)
		udesk := backend.Group("/udesk")
		{
			//分配工单提醒路由
			udesk.POST("/remind", logMiddleware(allocationreminders), remind.RemindHandler)
			//新工单提醒路由
			udesk.POST("/newticket", logMiddleware(newticketreminder), newticket.NewTicketHandler)
			//新回复提醒路由
			udesk.POST("/newreply", logMiddleware(newReplyreminder), newreply.NewreplyHandler)
		}

	}
	return r
}

// 路由中间件添加日志

func logMiddleware(app string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		logmgr.Log.WithFields(logrus.Fields{
			"app":     app,
			"path":    c.Request.URL.Path,
			"status":  c.Writer.Status(),
			"latency": end.Sub(start),
		}).Info("[GIN] Processing completed!")

	}
}
