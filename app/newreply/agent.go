package newreply

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"judgement/auth"
	logmgr "judgement/log"
	"net/http"
	"strings"
	"time"
)

type agent struct {
	agentName    string
	webhookUrl   string
	msgChan      chan *message
	ticketRecord map[string]string
}

type message struct {
	wechatMsg
	flag
}

type wechatMsg struct {
	Title        string `json:"标题"`
	Customer     string `json:"客户"`
	ReplyTime    string `json:"回复时间"`
	ReplyContent string `json:"回复内容"`
	TicketUrl    string `json:"工单地址"`
}

type flag struct {
	TicketID string `json:"工单id"`
	Acceptor string `json:"工单受理人"`
}

// newAgent 创建并返回一个新的 agent 实例
func newAgent(agentName, webhookUrl string) *agent {
	a := &agent{
		agentName:    agentName,
		webhookUrl:   webhookUrl,
		msgChan:      make(chan *message),
		ticketRecord: make(map[string]string, 8),
	}
	go a.run()
	return a
}

func (a *agent) run() {
	logmgr.Log.Println(a.agentName, a.webhookUrl, "agent run")
	for {
		select {
		case msg := <-a.msgChan:
			_, ok := a.ticketRecord[msg.TicketID]
			if !ok {
				a.ticketRecord[msg.TicketID] = msg.ReplyContent
				go msg.handleMessage(a)
			} else {
				a.ticketRecord[msg.TicketID] = msg.ReplyContent
			}

		}
	}
}

func checkReplyLastPerson(ticketId string) bool {
	url := auth.Geturlstring(fmt.Sprintf("https://servicecenter-alauda.udesk.cn//open_api_v1/tickets/%s/replies?", ticketId))
	resp, err := http.Get(url)
	if err != nil {
		logmgr.Log.Errorf("get ticket lastreply person data error!: %v", err)
		return false
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	jsonData := string(body)
	userType := gjson.Get(jsonData, "replies.#.author.user_type").Array()[0]
	if userType.String() == "agent" {
		return true
	}
	return false
}

// 每5min发送一次
func (m *message) handleMessage(agent *agent) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		//检查对应id的最新回复人来决定是否发送
		if !checkReplyLastPerson(m.TicketID) {
			sender := fmt.Sprintf("当前客户有新的回复，请关注并及时回复！\n标题: %s\n客户: %s\n回复时间: %s\n回复内容: %s\n工单地址: %s",
				m.Title,
				m.Customer,
				m.ReplyTime,
				agent.ticketRecord[m.TicketID],
				m.TicketUrl)
			agent.sendMsgToWxWorkRobot(sender)
		} else {
			logmgr.Log.Infof("受理人%v,工单%v已经回复", m.flag.Acceptor, m.TicketID)
			delete(agent.ticketRecord, m.TicketID)
			return
		}
	}
}

func (a *agent) sendMsgToWxWorkRobot(msg string) {
	text := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_mobile_list": ["@all"]}}`, msg)
	logmgr.Log.Infof("send to wechat message: %v\n", text)
	resp, err := http.Post(a.webhookUrl, "application/json", strings.NewReader(text))
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
