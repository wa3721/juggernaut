package newreply

import (
	"encoding/json"
	"log"
	"os"
)

// 创建一个新的 agentMgr 实例
var manager = newAgentMgr()

type agentMgr struct {
	agentList map[string]*agent
}

// newAgentMgr 创建并返回一个新的 agentMgr 实例
func newAgentMgr() *agentMgr {
	return &agentMgr{
		agentList: make(map[string]*agent),
	}
}

// AddAgent 向 agentMgr 添加一个新的 agent
func (am *agentMgr) AddAgent(agent *agent) {
	if agent != nil {
		am.agentList[agent.agentName] = agent
	}
}

// GetAgent 通过 agentName 获取对应的 agent 实例
func (am *agentMgr) GetAgent(agentName string) (*agent, bool) {
	agent, exists := am.agentList[agentName]
	return agent, exists
}

func init() {
	type tempAgent struct {
		AgentName  string `json:"客服"`
		WebhookUrl string `json:"webhook"`
	}
	data, _ := os.ReadFile("./app/newreply/webhook.json")

	// 创建一个切片来存储解析出的 agent 对象
	var tempAgents []tempAgent
	// 解析 JSON 数据到临时结构体
	err := json.Unmarshal(data, &tempAgents)
	if err != nil {
		log.Fatal(err)
	}

	// 为每个解析出的对象创建一个 agent 实例并添加到 manager
	for _, temp := range tempAgents {
		agent := newAgent(temp.AgentName, temp.WebhookUrl)
		manager.AddAgent(agent)
	}
	//// 通过 agentName 获取对应的 agent 实例
	//agent, exists := manager.GetAgent("王奥")
	//if exists {
	//	fmt.Printf("Found agent: %s, Webhook URL: %s, chan: %v\n", agent.agentName, agent.webhookUrl, agent.msgChan)
	//} else {
	//	log.Println("Agent not found")
	//}
	//for _, agent := range manager.agentList {
	//	fmt.Printf("Found agent: %s, Webhook URL: %s, chan: %v\n", agent.agentName, agent.webhookUrl, agent.msgChan)
	//}
}
