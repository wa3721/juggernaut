package reply

import (
	"context"
	"judgement/config"
	logmgr "judgement/config/log"
)

//定义代理管理者，存放客服的名字和他的代理

type assigneeAgentMgr struct {
	assignees map[string]*assigneeAgent
}

// 定义具体代理
type assigneeAgent struct {
	name        string                        //客服名字
	webhookUrl  string                        //客服的通知名字
	ticketMgr   map[string]*Reply             //存放最新回复对象，用于更新新回复 工单id:回复消息结构体
	replyChan   chan *Reply                   //用于接受回复消息的通道
	silenceChan chan *Silences                //用于接受静默消息的通道
	cancelMgr   map[string]context.CancelFunc //存放已经发送回复的对象的发送企业微信的上下文 工单id:对应goroutine的取消函数
}

// 定义一个全局代理管理者
var agmgr *assigneeAgentMgr

//初始化所有客服的代理

func InitAssigneeAgent() {
	agmgr = newAssigneeAgentMgr() //将初始化后的管理者赋值给全局变量 agmgr
	agmgr.initAssigneeAgent()     //初始化所有客服的代理
}

//初始化一个客服的管理者

func newAssigneeAgentMgr() *assigneeAgentMgr {
	return &assigneeAgentMgr{make(map[string]*assigneeAgent, 16)}
}

//启动一个代理

func newAssigneeAgent(name string, webUrl string) (agent *assigneeAgent) {
	agent = &assigneeAgent{
		name, webUrl, make(map[string]*Reply, 8), make(chan *Reply),
		make(chan *Silences), make(map[string]context.CancelFunc, 8),
	}
	go agent.run()
	return agent
}

//管理者方法：启动所有代理

func (m *assigneeAgentMgr) initAssigneeAgent() {
	for _, part := range config.Parts {
		agent := newAssigneeAgent(part.Name, part.ReplyWebhookUrl) //新建一个代理
		m.assignees[part.Name] = agent                             //将客服名字作为代理的键，代理对象作为代理管理者对应的值
		logmgr.Log.Infof("agent init success! %v", m.assignees[part.Name].name)
	}
}

//一个代理的监听方法

func (a *assigneeAgent) run() {
	for {
		select {
		//当从回复通道接受到回复对象的时候
		case reply := <-a.replyChan:
			logmgr.Log.Infof("assigneeAgent %v receive reply %v", reply.Assignee, reply)
			//判断回复对象是否存在在每个对象的工单管理者map中
			_, ok := a.ticketMgr[reply.CloudId]
			//不存在
			if !ok {
				//定义一个context注入发送企业微信的goroutine中，用于后续取消
				ctx, cancel := context.WithCancel(context.Background())
				//存入对应的取消函数管理者
				a.cancelMgr[reply.CloudId] = cancel
				//存入对应的回复管理者
				a.ticketMgr[reply.CloudId] = reply
				//发送企业微信
				go a.sendMsgToWxWorkRobot(ctx, reply)
			} else {
				//已经存在在回复管理者中，说明已经启动了发送goroutine，只需要在发送goroutine中动态更新对应的回复内容
				//后续根据工单id对应读取新的内容
				a.ticketMgr[reply.CloudId] = reply
			}
		//当从静默通道接受到静默对象的时候
		case silence := <-a.silenceChan:
			//在静默管理者中找到对应的工单id
			_, ok := a.cancelMgr[silence.CloudId]
			if ok {
				//找到第一次发送时存放的cancel函数，停止对应已经在执行的goroutine
				cancel := a.cancelMgr[silence.CloudId]
				cancel()
				//清除掉静默管理者中的cancel函数，防止数据无限增加
				delete(a.cancelMgr, silence.CloudId)
			}
		}
	}
}
