package newreply

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	logmgr "judgement/log"
	"net/http"
	"strings"
)

func NewreplyHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logmgr.Log.Error("Error reading request body!")
	}
	newRpy := newReply()
	json.Unmarshal(body, &newRpy)
	var builder strings.Builder
	message := fmt.Sprintf("当前客户有新的回复，请关注并及时回复！\n\n标题: %s\n客户: %s\n回复时间: %s\n回复内容: %s\n工单链接: %s",
		newRpy.title,
		newRpy.customer,
		newRpy.replyTime,
		newRpy.replyContent,
		newRpy.ticketUrl)
	builder.WriteString(message)
	//msg := builder.String()

}

type reply struct {
	title        string `json:"标题"`
	customer     string `json:"客户"`
	replyTime    string `json:"回复时间"`
	replyContent string `json:"回复内容"`
	ticketUrl    string `json:"工单地址"`
}

func newReply() *reply {
	return &reply{}
}

func sendTowechatBot(msg, wxWorkRobotURL string) {
	message := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_mobile_list": ["@all"]}}`, msg)
	logmgr.Log.Infof("send to wechat message: %v", message)
	resp, err := http.Post(wxWorkRobotURL, "application/json", strings.NewReader(message))
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
