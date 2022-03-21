package bumscd

type Config struct {
	Path         string    `json:"path,omitempty"`
	Method       string    `json:"method,omitempty"`
	GWName       string    `json:"gw,omitempty"`
	ServiceCode  string    `json:"service_code,omitempty"`
	ServiceScene string    `json:"service_scene,omitempty"`
	ReqMapping   *Relation `json:"req_mapping,omitempty"`
	RespMapping  *Relation `json:"resp_mapping,omitempty"`
}
