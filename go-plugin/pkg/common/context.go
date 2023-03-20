package common

const (
	VarStartTime                      string = "start_time"
	VarRequestReceivedDuration        string = "request_received_duration"
	VarResponseReceivedDuration       string = "response_received_duration"
	VarRequestFinishedDuration        string = "request_finished_duration"
	VarProcessTimeDuration            string = "process_time_duration"
	VarBytesSent                      string = "bytes_sent"
	VarBytesReceived                  string = "bytes_received"
	VarProtocol                       string = "protocol"
	VarResponseCode                   string = "response_code"
	VarDuration                       string = "duration"
	VarResponseFlag                   string = "response_flag"
	VarResponseFlags                  string = "response_flags"
	VarUpstreamLocalAddress           string = "upstream_local_address"
	VarDownstreamLocalAddress         string = "downstream_local_address"
	VarDownstreamRemoteAddress        string = "downstream_remote_address"
	VarUpstreamHost                   string = "upstream_host"
	VarUpstreamTransportFailureReason string = "upstream_transport_failure_reason"
	VarUpstreamCluster                string = "upstream_cluster"
	VarRequestedServerName            string = "requested_server_name"
	VarRouteName                      string = "route_name"
	VarProtocolConfig                 string = "protocol_config"

	// ReqHeaderPrefix is the prefix of request header's formatter
	VarPrefixReqHeader string = "request_header_"
	// RespHeaderPrefix is the prefix of response header's formatter
	VarPrefixRespHeader string = "response_header_"
)
