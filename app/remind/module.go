package remind

import (
	"encoding/json"
	logmgr "judgement/log"
	"os"
)

type WorkOrder struct {
	OrderID         string `json:"工单id"`
	Customer        string `json:"客户"`
	CustomerService string `json:"受理客服"`
	Subject         string `json:"主题"`
	Priority        string `json:"优先级"`
	WorkUrl         string `json:"工单链接"`
	CellPhone       string `json:"客户手机号"`
	CreateTime      string `json:"创建时间"`
}

var PhoneList = make(map[string]string, 16)

var Order *WorkOrder

func init() {
	phonelist, err := os.ReadFile("./app/remind/phonelist.json")
	if err != nil {
		logmgr.Log.Fatalf("Error reading file: %v\n", err)
		return
	}
	err = json.Unmarshal(phonelist, &PhoneList)
	if err != nil {
		logmgr.Log.Fatalf("Error unmarshaling JSON: %v\n", err)
		return
	}
}
