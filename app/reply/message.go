package reply

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	logmgr "judgement/config/log"
	"strings"
)

//因为不知道udesk触发器发过来的消息是静默还是回复，所以定义一个接口

type Message interface {
	isMessage()
}

// 两个空函数，没什么用，单纯是为了把Reply和Silences 包在接口中
func (r *Reply) isMessage()    {}
func (s *Silences) isMessage() {}

//回复内容

type Reply struct {
	CloudId       string `json:"cloudId"`
	Assignee      string `json:"assignee"`
	TicketUser    string `json:"ticketUser"`
	Subject       string `json:"subject"`
	WebUrl        string `json:"webUrl"`
	LatestComment string `json:"latest_comment"`
	UdeskId       string `json:"udeskId"`
	ctx           context.Context
}

//静默消息

type Silences struct {
	CloudId  string `json:"CloudId"`
	Assignee string `json:"assignee"`
	Silence  string `json:"silence"`
}

// 初始化msg，返回Message对象
func newMsg(msg io.Reader) (Message, error) {
	//定义两个新对象
	var reply Reply
	var silences Silences
	body, err := io.ReadAll(msg)
	if err != nil {
		logmgr.Log.Errorf("init new msg err: %v", err)
		return nil, err
	}
	//如果是回复内容
	if !strings.Contains(string(body), "silence") {
		//解析成回复内容结构体，返回
		if err := json.Unmarshal(body, &reply); err == nil {
			logmgr.Log.Infof("received reply instance: %v", reply)
			return &reply, nil

		}
	} else {
		//如果不是
		//解析成静默内容结构体，返回
		if err := json.Unmarshal(body, &silences); err == nil {
			logmgr.Log.Infof("received silences instance: %v", silences)
			return &silences, nil
		}

	}
	//什么都不是
	return nil, errors.Errorf("unknown message type")
}
