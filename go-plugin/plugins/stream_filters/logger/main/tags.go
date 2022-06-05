package main

const (
	TAG_BEGIN = iota - 1
	SERVICE_NAME
	SPAN_TYPE
	METHOD_NAME
	REQUEST_URL
	TRACEID
	SPANID
	RESULT_STATUS
	DOWN_PROTOCOL
	UP_PROTOCOL
	LISTENER_ADDRESS
	UPSTREAM_HOST_ADDRESS
	DOWNSTEAM_HOST_ADDRESS
	DURATION
	START_TIME
	END_TIME
	MOSN_DURATION
	MOSN_REQUEST_DURATION
	MOSN_RSPONSE_DURATION
	APP_NAME
	TAG_END
)

var (
	tagsName = map[int]string{
		SERVICE_NAME:           "serviceName",
		METHOD_NAME:            "method",
		DOWN_PROTOCOL:          "callerProtocol",
		UP_PROTOCOL:            "targetProtocol",
		RESULT_STATUS:          "status",
		UPSTREAM_HOST_ADDRESS:  "targetAddress",
		DOWNSTEAM_HOST_ADDRESS: "callerAddress",
		SPAN_TYPE:              "kind",
		REQUEST_URL:            "requestUrl",
		DURATION:               "duration",
		START_TIME:             "startTime",
		END_TIME:               "endTime",
		LISTENER_ADDRESS:       "listenerAddress",
		MOSN_DURATION:          "mosnDuration",
		MOSN_REQUEST_DURATION:  "mosnRequestDuration",
		MOSN_RSPONSE_DURATION:  "mosnReponseDuration",
		TRACEID:                "traceId",
		SPANID:                 "spanId",
		APP_NAME:               "app_name",
	}
)
