package reply

import (
	"encoding/json"
	"fmt"
	"io"
	"judgement/config"
	"judgement/config/log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Reply struct {
	CloudId       string       `json:"cloudId"`
	Assignee      string       `json:"assignee"`
	TicketUser    string       `json:"ticketUser"`
	Subject       string       `json:"subject"`
	WebUrl        string       `json:"webUrl"`
	LatestComment string       `json:"latest_comment"`
	HttpClient    *http.Client `json:",omitempty"`
}

var replyMutex sync.Mutex

// 构造通知消息内容结构体，返回body Reader
func (r *Reply) constructMessageRequestBody(messageContent string) (*strings.Reader, error) {
	data := &struct {
		Msgtype string `json:"msgtype"`
		Text    struct {
			Content       string   `json:"content"`
			MentionedList []string `json:"mentioned_list"`
		} `json:"text"`
	}{
		Msgtype: "text",
		Text: struct {
			Content       string   "json:\"content\""
			MentionedList []string "json:\"mentioned_list\""
		}{Content: messageContent, MentionedList: []string{"@all"}},
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(string(body)), nil
}

// 生成消息内容
func (r *Reply) generateMessage() string {
	return fmt.Sprintf(
		"当前客户有新的回复，请及时查看\ncloud_id：%v\n工单标题：%v\n客户：%v\n回复内容：\n%v\n工单链接：\n%v",
		r.CloudId,
		r.Subject,
		r.TicketUser,
		r.LatestComment,
		r.WebUrl,
	)
}

// 传入客服处理人对应的webhookUrl，发送通知
func (r *Reply) send(webhookUrl string) {
	logmgr.Log.Info("Current send user %v", r.Assignee)

	// 生成消息内容（共享资源加锁）
	replyMutex.Lock()
	messageContent := r.generateMessage()
	replyMutex.Unlock()

	// 构造通知消息body
	payload, err := r.constructMessageRequestBody(messageContent)
	if err != nil {
		logmgr.Log.Error("Error construct message post body, json marshal error %v", err)
		return
	}

	// 创建请求，传入webhookUrl，通知消息body
	req, err := http.NewRequest(http.MethodPost, webhookUrl, payload)
	if err != nil {
		logmgr.Log.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := r.HttpClient.Do(req)
	if err != nil {
		logmgr.Log.Error(err)
		return
	}

	// 判断返回是否成功，否则打印返回
	if resp.StatusCode != http.StatusOK {
		logmgr.Log.Error("Error respone code not 200 %v", err)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logmgr.Log.Error(err)
			return
		}
		logmgr.Log.Error("Error respone body content", string(respBody))
		return
	}
}

// 判断客服处理人是否在成员列表中并且返回客服处理人的webhookUrl
func (r *Reply) assigneeExists(name string, parts []config.Part) (string, bool) {
	for _, v := range parts {
		if r.Assignee == v.Name && v.ReplyWebhookUrl != "" {
			return v.ReplyWebhookUrl, true
		}
	}
	return "", false
}

// 初始化Reply
func newReply(r *http.Request) (*Reply, error) {
	re := &Reply{
		HttpClient: &http.Client{},
	}

	err := json.NewDecoder(r.Body).Decode(&re)
	if err != nil {
		return nil, err
	}

	return re, nil
}

// 主逻辑
func ReplyHandler(c *gin.Context) {
	// 初始化reply
	re, err := newReply(c.Request)
	if err != nil {
		logmgr.Log.Error(err)
		return
	}

	// 从配置中读取当前拥有的成员
	parts := []config.Part{}
	viper.UnmarshalKey("parts", &parts)

	// 判断客服处理人在配置成员中是否存在
	// 并且返回客服处理人的webhookUrl
	webhookUrl, exists := re.assigneeExists(re.Assignee, parts)
	if !exists {
		logmgr.Log.Error("No assignee or webhhookUrl is empty - %v", re.Assignee)
		return
	}

	// 将客户回复的通知发送给对应客服处理人的webhookUrl
	go re.send(webhookUrl)

	c.Writer.Write([]byte(`{"ok":true}`))
}
