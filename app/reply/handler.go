package reply

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	logmgr "judgement/config/log"
	"net/http"
	"strings"
	"time"
)

// 主逻辑

func ReplyHandler(c *gin.Context) {
	msg, err := newMsg(c.Request.Body)
	if err != nil {
		logmgr.Log.Errorf("init reply instance error:%v", err)
		return
	}
	//因为newMsg返回的是接口对象，所以要进行断言才能处理
	switch msg.(type) {
	case *Reply:
		// 处理Reply对象
		// 判断客服处理人在配置成员中是否存在
		if checkAssigneeAgentExist(msg.(*Reply)) {
			agmgr.assignees[msg.(*Reply).Assignee].replyChan <- msg.(*Reply) //发送到对应工单受理人的代理的回复管道中
		} else {
			//不存在，直接返回
			logmgr.Log.Errorf("reply not find agent %v!", msg.(*Reply).Assignee)
			return
		}
	case *Silences:
		// 处理Silences对象
		//发送到对应agent的静默通道中
		agmgr.assignees[msg.(*Silences).Assignee].silenceChan <- msg.(*Silences)
	default:
		//什么对象也不是，不做处理
		logmgr.Log.Errorf("Unknown message type")
	}
}

func checkAssigneeAgentExist(reply *Reply) bool {
	_, ok := agmgr.assignees[reply.Assignee]
	if !ok {
		logmgr.Log.Infof("agent assignee %v not exist!", reply.Assignee)
		return false
	}
	logmgr.Log.Infof("agent assignee %v is exist!", reply.Assignee)
	return true
}

// 生成消息内容
func (r *Reply) generateMessage(LatestComment string) string {
	return fmt.Sprintf(
		"当前客户有新的回复，请及时查看\n工单id：%v\n工单标题：%v\n客户：%v\n回复内容：\n%v\n工单链接：\n%v",
		r.CloudId,
		r.Subject,
		r.TicketUser,
		LatestComment,
		r.WebUrl,
	)
}

func (a *assigneeAgent) sendMsgToWxWorkRobot(ctx context.Context, r *Reply) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		select {
		//主动静默的情况
		case <-ctx.Done():
			logmgr.Log.Infof("cloudid: %v sendMsgToWxWorkRobot: context canceled", r.CloudId)
			//取消之前删掉对应的工单id的回复对象，防止数据无限增加
			delete(a.ticketMgr, r.CloudId)
			return
		default:
			//这里需要判断被动取消的情况，即：有客服回复了工单
			//正常循环发送数据
			//动态读取最新回复,近发送最新的回复
			message := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_mobile_list": ["@all"]}}`, r.generateMessage(a.ticketMgr[r.CloudId].LatestComment))
			logmgr.Log.Infof("send to wechat message: %v", message)
			resp, err := http.Post(a.webhookUrl, "application/json", strings.NewReader(message))
			if err != nil {
				logmgr.Log.Errorf("Error sending message to Wechat Bot: %v", err)
				return
			}
			defer resp.Body.Close()

			// 读取响应
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				logmgr.Log.Errorf("Error reading response body: %v", err)
				return
			}
			logmgr.Log.Info("Message send successfully. Response:", string(respBody))
		}

	}

}
