package remind

import (
	"encoding/json"
	"fmt"
	"io"
	logmgr "judgement/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RemindHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logmgr.Log.Errorf("Error reading request body!")
		return
	}
	json.Unmarshal(body, &Order)
	var builder strings.Builder //定义的企业微信机器人发送的字符串
	message := fmt.Sprintf("工单已分配!\n客户: %s\n客户手机号: %s\n受理客服: %s\n工单id: %s\n工单链接: %s\n工单主题: %s",
		Order.Customer,
		Order.CellPhone,
		Order.CustomerService,
		Order.OrderID,
		Order.WorkUrl,
		Order.Subject)
	builder.WriteString(message)
	msg := builder.String()
	for _, v := range Parts {
		if Order.CustomerService == v.Name {
			sendMsgToWxWorkRobot(msg, v.Name)
			return
		}
	}
	sendMsgToWxWorkRobot(msg, "")

}

func sendMsgToWxWorkRobot(msg, phone string) {
	wxWorkRobotURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=7b50e0c4-8a35-4f29-b652-25e1d6142c2b"
	message := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_mobile_list": ["%s"]}}`, msg, phone)
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
