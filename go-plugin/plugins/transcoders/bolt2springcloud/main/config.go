package main

type Config struct {
	UniqueId    string      `json:"unique_id"`
	Path        string      `json:"path"`
	Method      string      `json:"method"`
	TragetApp   string      `json:"target_app"`
	Class       string      `json:"class"`
	ReqMapping  interface{} `json:"-"`
	RespMapping interface{} `json:"-"`
}
