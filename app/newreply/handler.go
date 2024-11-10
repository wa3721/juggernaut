package newreply

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	logmgr "judgement/log"
)

func NewreplyHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logmgr.Log.Errorf("Error reading request body!")
		return
	}
	err = json.Unmarshal(body, &receive)
	if err != nil {
		logmgr.Log.Errorf("unmarshal json failed!,%v", err)
		return
	}
	Agent, ok := manager.GetAgent(receive.Acceptor)
	if ok {
		Agent.msgChan <- receive
	}

}
