package captcha

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	logmgr "judgement/log"
	"net/http"
	"net/url"
	"strings"
)

//这里处理业务逻辑

func CaptchaHandler(c *gin.Context) {
	// 读取请求体数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logmgr.Log.Error("Error reading request body!")
		return
	}
	func([]byte) {
		content := string(body)
		logmgr.Log.Debug(content)
		spiltDecodedValue1 := strings.Split(content, "&")
		spiltDecodedValue2 := strings.Split(spiltDecodedValue1[1], "=")
		decodedValue, error := url.QueryUnescape(spiltDecodedValue2[1])
		if error != nil {
			logmgr.Log.Errorf("message string convert error: %s", error)
		}
		content = decodedValue
		logmgr.Log.Info(content)
		keyword := []string{"验证码", "华为", "动态口令", "短信口令"}
		if containAll(content, keyword) {
			sendToWechatBot(content)
		}
	}(body)
}

func containAll(mainString string, subString []string) bool {
	for i := range subString {
		if !strings.Contains(mainString, subString[i]) {
			continue
		} else {
			return true
		}
	}
	return false
}

func sendToWechatBot(content string) {

	// 这里实现发送信息到企业微信机器人的逻辑
	// 可以使用企业微信机器人的 Webhook API 发送消息
	webhookURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=7b50e0c4-8a35-4f29-b652-25e1d6142c2b"
	var message string

	if strings.Contains(content, "华为") {
		// 如果 content 中包含 "华为"
		message = fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s","mentioned_mobile_list":["13241102589"]}}`, content)
	} else {
		// 如果 content 中不包含 "华为"
		message = fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s"}}`, content)
	}
	// 发送 HTTP POST 请求到企业微信机器人的 Webhook
	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(message))
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
