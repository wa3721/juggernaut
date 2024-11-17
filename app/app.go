package app

import (
	"fmt"
	"judgement/app/captcha"
	"judgement/app/newticket"
	"judgement/app/remind"
	"judgement/app/reply"
	"judgement/config"
	"judgement/config/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	verifycode          = "验证码"
	allocationreminders = "工单分配提醒"
	newticketreminder   = "新工单提醒"
	replyreminder       = "新回复提醒"

	authUser     = "admin"
	authPassword = "123456"
)

//New judgement 实例

func NewApp(app string) *gin.Engine {
	fmt.Println(` 
_________          _______  _______  _______  _______  _        _______          _________
\__    _/|\     /|(  ____ \(  ____ \(  ____ \(  ____ )( (    /|(  ___  )|\     /|\__   __/
   )  (  | )   ( || (    \/| (    \/| (    \/| (    )||  \  ( || (   ) || )   ( |   ) (   
   |  |  | |   | || |      | |      | (__    | (____)||   \ | || (___) || |   | |   | |   
   |  |  | |   | || | ____ | | ____ |  __)   |     __)| (\ \) ||  ___  || |   | |   | |   
   |  |  | |   | || | \_  )| | \_  )| (      | (\ (   | | \   || (   ) || |   | |   | |   
|\_)  )  | (___) || (___) || (___) || (____/\| ) \ \__| )  \  || )   ( || (___) |   | |   
(____/   (_______)(_______)(_______)(_______/|/   \__/|/    )_)|/     \|(_______)   )_(
	`)
	// 加载配置
	config.NewConfig("./config.yaml").LoadConfig()
	logmgr.LoadLogConfig()
	//加载reply代理
	//每个客服有一个专门用于接受消息的代理，这个代理的作用是将接受到的消息发送到根据受理客服人的名字将消息发送到对应人的通道中
	reply.InitAssigneeAgent()
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	backend := r.Group("/" + app)
	{ //验证码转发路由
		backend.POST("/captcha", logMiddleware(verifycode), captcha.CaptchaHandler)
		udesk := backend.Group("/udesk")
		{
			//分配工单提醒路由
			udesk.POST("/remind", logMiddleware(allocationreminders), gin.BasicAuth(gin.Accounts{authUser: authPassword}), remind.RemindHandler)
			//新工单提醒路由
			udesk.POST("/newticket", logMiddleware(newticketreminder), gin.BasicAuth(gin.Accounts{authUser: authPassword}), newticket.NewTicketHandler)
			//新回复提醒路由
			udesk.POST("/reply", logMiddleware(replyreminder), gin.BasicAuth(gin.Accounts{authUser: authPassword}), reply.ReplyHandler)
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
