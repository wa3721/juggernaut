package newticket

import (
	"encoding/json"
	"fmt"
	"io"
	"judgement/config"
	logmgr "judgement/config/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type ticket struct {
	Title            string `json:"标题"`
	Level            string `json:"级别"`
	Environment      string //这里通过gjson传string
	Creator          string `json:"提单人"`
	TicketCreateTime string `json:"提单时间"`
	TicketUrl        string `json:"工单地址"`
}

func NewTicketHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logmgr.Log.Error("Error reading request body!")
	}
	var newTic *ticket
	json.Unmarshal(body, &newTic)
	bodyjson := string(body)
	envs := gjson.Get(bodyjson, "环境").Array()
	for _, v := range envs {
		if v.String() == "<空>" {
			continue
		} else {
			newTic.Environment = v.String()
			break
		}
	}
	if newTic.Environment == "" {
		newTic.Environment = "环境信息暂无┭┮﹏┭┮"
	}
	var builder strings.Builder
	message := fmt.Sprintf("当前有新的工单，请及时处理!\n标题: %s\n级别: %s\n环境: %s\n提单人: %s\n提单时间: %s\n工单地址: %s",
		newTic.Title,
		newTic.Level,
		newTic.Environment,
		newTic.Creator,
		newTic.TicketCreateTime,
		newTic.TicketUrl)
	builder.WriteString(message)
	logmgr.Log.Info(message)
	sendTowechatBot(message, config.NewTicketWebhookUrl)
}

func sendTowechatBot(msg, wxWorkRobotURL string) error {
	message := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_mobile_list": ["@all"]}}`, msg)
	logmgr.Log.Infof("send to wechat message: %v", message)

	var resp *http.Response
	var err error

	//重试循环new
	for i := 0; i < 3; i++ {
		resp, err = http.Post(wxWorkRobotURL, "application/json", strings.NewReader(message))
		if err != nil {
			logmgr.Log.Errorf("Error sending message to Wechat Bot (Attempt %d/%d): %v", i+1, 3, err)
			continue // 如果请求错误，继续循环
		}

		// 读取响应体并解析 errcode
		defer resp.Body.Close()
		respBody, err1 := io.ReadAll(resp.Body)
		if err1 != nil {
			logmgr.Log.Errorf("Error reading response body (Attempt %d/%d): %v", i+1, 3, err1)
			continue // 如果读取响应错误，继续循环
		}

		// 解析响应的 JSON 数据并检查 errcode 字段
		var result map[string]interface{}
		if json.Unmarshal(respBody, &result) != nil {
			logmgr.Log.Errorf("Error unmarshaling response JSON (Attempt %d/%d): %v", i+1, 3, err1)
			continue // 如果 JSON 解析错误，继续循环
		}

		// 判断 errcode 是否为 0
		if errcode, ok := result["errcode"].(float64); ok && errcode == 0 {
			logmgr.Log.Info("Message send successfully. Response:", string(respBody))
			break // 成功发送并且 errcode 为 0，跳出循环
		} else {
			logmgr.Log.Errorf("Received non-zero errcode from Wechat Bot (Attempt %d/%d): %v", i+1, 3, result["errcode"])
		}

		if err != nil || err1 != nil {
			return fmt.Errorf("Failed to send message to Wechat Bot after 3 attempts: %v, %v", err, err1)
		}
	}
	return nil
}
