package newreply

import (
	"encoding/json"
	logmgr "judgement/log"
	"os"
)

var WebhookList = make(map[string]string, 16)

func init() {
	webhooklist, err := os.ReadFile("./app/newreply/webhook.json")
	if err != nil {
		logmgr.Log.Fatalf("Error reading file: %v\n", err)
		return
	}
	err = json.Unmarshal(webhooklist, &WebhookList)
	if err != nil {
		logmgr.Log.Fatalf("Error unmarshaling JSON: %v\n", err)
		return
	}
}
