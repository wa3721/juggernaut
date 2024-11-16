package remind

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

//接受触发器请求

var Order *WorkOrder
