package main

const (
	LogReporter        string = "log"
	GRPCReporter       string = "gRPC"
	DefaultServiceName string = "mosn"
)

type SkyWalkingTraceConfig struct {
	Reporter         string                   `json:"reporter"`
	BackendService   string                   `json:"backend_service"`
	ServiceName      string                   `json:"service_name"`
	MaxSendQueueSize string                   `json:"max_send_queue_size"`
	VmMode           string                   `json:"vmmode"`
	PodName          string                   `json:"pod_name"`
	Authentication   string                   `json:"authentication"`
	TLS              SkyWalkingTraceTLSConfig `json:"tls"`
}

type SkyWalkingTraceTLSConfig struct {
	CertFile           string `json:"cert_file"`
	ServerNameOverride string `json:"server_name_override"`
}
