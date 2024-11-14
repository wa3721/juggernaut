package reply

import (
	"encoding/json"
	"fmt"
	"io"
	"judgement/config"
	logmgr "judgement/log"
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

func (r *Reply) constructReplayBody(messageContent string) (*strings.Reader, error) {
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
		logmgr.Log.Error("Error construct confluence release post body, json marshal error %v", err)
	}
	return strings.NewReader(string(body)), nil
}

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

func (r *Reply) send(webhookUrl string) {
	logmgr.Log.Info("Current send user %v", r.Assignee)

	replyMutex.Lock()
	messageContent := r.generateMessage()
	replyMutex.Unlock()

	payload, err := r.constructReplayBody(messageContent)
	if err != nil {
		logmgr.Log.Error("Error construct release body %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, webhookUrl, payload)
	if err != nil {
		logmgr.Log.Error("Error construct confluence release post body, json marshal error %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		logmgr.Log.Error("Error construct confluence release post body, json marshal error %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		logmgr.Log.Error("Error respone code not 200 %v", err)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logmgr.Log.Error("Error read respone body %v", err)
			return
		}
		logmgr.Log.Error("Error respone body content", string(respBody))
		return
	}
}

func (r *Reply) assigneeNotExists(name string, parts []config.Users) (string, bool) {
	for _, v := range parts {
		if r.Assignee == v.Name {
			return v.ReplyRobotWebhookUrl, true
		}
	}
	return "", false
}

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

func ReplyHandler(c *gin.Context) {
	// 初始化reply
	re, err := newReply(c.Request)
	if err != nil {
		logmgr.Log.Error(err)
		return
	}

	// 从配置中读取当前拥有的成员
	parts := []config.Users{}
	viper.UnmarshalKey("parts", &parts)

	// 判断工单客服处理人是否存在配置中的成员中，并且返回该成员的机器人webhookUrl
	webhookUrl, exists := re.assigneeNotExists(re.Assignee, parts)
	if !exists {
		logmgr.Log.Error("No assignee %v", re.Assignee)
		return
	}

	// 传递工单客服处理人对应的webhookUrl并发送
	go re.send(webhookUrl)

	c.Writer.Write([]byte(`{"ok":true}`))
}
