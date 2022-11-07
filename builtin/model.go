package builtin

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Status struct {
	Nodes []*NodeInfo `json:"nodes"`
}

type Leader struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type NodeInfo struct {
	ID       string  `json:"id"`
	IsLeader bool    `json:"is_leader"`
	Leader   *Leader `json:"leader"`
}

type JoinRequest struct {
	GroupID string `json:"group_id"`
	NodeID  string `json:"node_id"`
	Address string `json:"address"`
}
