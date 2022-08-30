package main

import (
	"strconv"
	"strings"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/plugins/trace/sofa/main/generator"
	"mosn.io/pkg/buffer"
)

/*
rpc-client-digest.log
2021-02-05 11:57:28.284,reservation-service,1e49a2be1612497448281145168886,0,com.alipay.ms.service.EchoService,echo,Dubbo,SYNC,30.73.162.190:30800,reservation-service,GZ00B,,,,00,242B,118B,2ms,0ms,0ms,0ms,,,,,,30.73.162.190,57510,,F,,mosn_cluster=cb7aaa7a31fd0490299c530e83a868231&mosn_namespace=default&mosn_log=true&mosn_cluster=cb7aaa7a31fd0490299c530e83a868231&mosn_namespace=default&mosn_log=true&mosn_data_id=com.alipay.sofa.ms.service.EchoService:aaa:bbb@dubbo&mosn_data_ver=&,325.799µs
*/

func (s *SofaRPCSpan) isPod() string {
	if !s.pod {
		return "T"
	}
	return "F"
}

func (s *SofaRPCSpan) status(printData buffer.IoBuffer) {
	statusCode, _ := strconv.Atoi(s.tags[generator.RESULT_STATUS])
	var code = "02"
	if statusCode == api.SuccessCode {
		code = "00"
	} else if statusCode == api.TimeoutExceptionCode {
		code = "03"
	} else if statusCode == api.RouterUnavailableCode || statusCode == api.NoHealthUpstreamCode {
		code = "04"
	} else {
		code = "02"
	}
	printData.WriteString(code + ",")
}

func (s *SofaRPCSpan) clientRpcLogger() error {
	printData := buffer.NewIoBuffer(512)
	//1 日志打印时间:2021-02-05 13:46:44.264
	date := s.endTime.Format("2006-01-02 15:04:05.000")
	printData.WriteString(date + ",")
	//2 当前应用名:echo-server
	printData.WriteString(s.appName + ",")
	//3 TraceId:0ad55acc1512711661769715320278
	printData.WriteString(s.tags[generator.TRACE_ID] + ",")
	//4 RpcId:0
	printData.WriteString(s.tags[generator.SPAN_ID] + ",")
	//5 服务名:com.alipay.ms.service.SofaEchoService:1.0
	printData.WriteString(s.tags[generator.SERVICE_NAME] + ",")
	//6 方法名:echo
	printData.WriteString(s.tags[generator.METHOD_NAME] + ",")
	//7 协议:bolt
	printData.WriteString(s.tags[generator.PROTOCOL] + ",")
	//8 调用方式:SYNC
	printData.WriteString("SYNC" + ",")
	//9 目标地址:30.73.162.190:12200
	printData.WriteString(s.tags[generator.UPSTREAM_HOST_ADDRESS] + ",")
	//10 目标系统名:echo-server
	printData.WriteString(s.tags[generator.TARGET_APP_NAME] + ",")
	//11 目标Zone:GZ00B
	printData.WriteString(s.tags[generator.TARGET_CELL] + ",")
	//12 目标IDC:
	printData.WriteString(s.tags[generator.TARGET_IDC] + ",")
	//13 目标City:
	printData.WriteString(s.tags[generator.TARGET_CITY] + ",")
	//14 uid:
	printData.WriteString(s.tags[generator.UID] + ",")
	// status
	s.status(printData)
	//16 请求大小:388B
	printData.WriteString(s.tags[generator.REQUEST_SIZE] + ",")
	//17 响应大小:105B
	printData.WriteString(s.tags[generator.RESPONSE_SIZE] + ",")
	//18 调用耗时:2ms
	printData.WriteString(s.tags[generator.DURATION] + ",")
	//19 链接建立耗时:0ms
	printData.WriteString("0ms,")
	//20 请求序列化耗时:0ms
	printData.WriteString("0ms,")
	//21 超时参考耗时:0ms
	printData.WriteString("0ms,")
	//22 当前线程名:
	printData.WriteString(",")
	//23 路由记录:
	printData.WriteString(s.tags[generator.ROUTE_RECORD] + ",")
	//24 弹性数据位:
	printData.WriteString(",")
	//25 是否需要弹性:
	printData.WriteString(",")
	//26 转发的服务名称:
	printData.WriteString(s.tags[generator.CALLER_APP_NAME] + ",")
	//27 clientip: 127.0.0.1
	downStreamHostAddress := strings.Split(s.tags[generator.DOWNSTEAM_HOST_ADDRESS], ":")
	if len(downStreamHostAddress) > 0 {
		localIp := strings.Split(s.tags[generator.DOWNSTEAM_HOST_ADDRESS], ":")[0]
		printData.WriteString(localIp + ",")
	} else {
		printData.WriteString(",")
	}
	//28 clientport: 61181
	if len(downStreamHostAddress) > 1 {
		localPort := strings.Split(s.tags[generator.DOWNSTEAM_HOST_ADDRESS], ":")[1]
		printData.WriteString(localPort + ",")
	} else {
		printData.WriteString(",")
	}
	//29 当前zone:
	printData.WriteString(s.tags[generator.CALLER_CELL] + ",")
	//30 是否物理机器: F
	printData.WriteString(s.isPod() + ",")
	//31 systemMap:
	printData.WriteString(",")
	//32 系统穿透数据
	printData.WriteString(s.tags[generator.BAGGAGE_DATA] + ",")
	//33 mosn处理时间
	printData.WriteString(s.tags[generator.MOSN_PROCESS_TIME])
	printData.WriteString("\n")
	return s.egressLogger.Print(printData, true)
}

/*
//rpc-server-digest
2021-02-05 11:57:53.431,reservation-service,1e49a2be1612497473430148068886,0,com.alipay.ms.service.EchoService,echo,Dubbo,,30.73.162.190:57523,,,,1ms,0ms,,00,,,0,,mosn_tls_state=off&mosn_cluster=cb7aaa7a31fd0490299c530e83a868231&mosn_namespace=default&mosn_log=true&mosn_tls_state=off&mosn_cluster=cb7aaa7a31fd0490299c530e83a868231&mosn_namespace=default&mosn_log=true&mosn_data_id=com.alipay.sofa.ms.service.EchoService:aaa:bbb@dubbo&mosn_data_ver=&,285.268µs,242B,118B
*/

func (s *SofaRPCSpan) serverRpcLogger() error {
	printData := buffer.NewIoBuffer(512)
	//1.日志打印时间
	date := s.endTime.Format("2006-01-02 15:04:05.000")
	printData.WriteString(date + ",")
	//2当前应用名
	printData.WriteString(s.appName + ",")
	//3TraceId
	printData.WriteString(s.tags[generator.TRACE_ID] + ",")
	//4RpcId
	printData.WriteString(s.tags[generator.SPAN_ID] + ",")
	//5服务名
	printData.WriteString(s.tags[generator.SERVICE_NAME] + ",")
	//6方法名
	printData.WriteString(s.tags[generator.METHOD_NAME] + ",")
	//7协议
	printData.WriteString(s.tags[generator.PROTOCOL] + ",")
	//8调用方式
	printData.WriteString("SYNC" + ",")
	//9调用者 URL
	printData.WriteString(s.tags[generator.DOWNSTEAM_HOST_ADDRESS] + ",")
	//10调用者应用名
	printData.WriteString(s.tags[generator.CALLER_APP_NAME] + ",")
	//11 目标Zone:GZ00B
	printData.WriteString(s.tags[generator.TARGET_CELL] + ",")
	//12 目标IDC:
	printData.WriteString(s.tags[generator.TARGET_IDC] + ",")
	//13请求处理耗时（ms）
	printData.WriteString(s.tags[generator.DURATION] + ",")
	//14服务端响应序列化耗时（ms）
	printData.WriteString(s.tags[generator.UPSTREAM_DURATION] + ",")
	//15当前线程名
	printData.WriteString(",")
	//
	s.status(printData)
	//17 beElasticServiceName（表明这次调用是转发调用,转发的服务名称和方法名称是啥值如：“com.test.service.testservice.TestService:1.0:biztest---doProcess”）
	printData.WriteString(",")
	//18 beElastic（表示没有被转发的处理）
	printData.WriteString("0ms,")
	//19 rpc线程池等待时间
	printData.WriteString(",")
	//20 系统穿透数据（kv 格式，用于传送系统灾备信息等）
	printData.WriteString(",")
	//21 穿透数据（kv格式）
	printData.WriteString(s.tags[generator.BAGGAGE_DATA] + ",")
	//22 mosn处理时间
	printData.WriteString(s.tags[generator.MOSN_PROCESS_TIME] + ",")
	//23 请求大小:388B
	printData.WriteString(s.tags[generator.REQUEST_SIZE] + ",")
	//24 响应大小:105B
	printData.WriteString(s.tags[generator.RESPONSE_SIZE])
	printData.WriteString("\n")
	return s.ingressLogger.Print(printData, true)
}

/*
springcloud-client-digest.log
2021-02-05 11:55:33.61,reservation-service,1e49a2be1612497333607134168886,0,127.0.0.1:10088,/echo/name,HTTP,SYNC,30.73.162.190:10080,reservation-service,GZ00B,,,,00,18B,20B,2ms,0ms,0ms,0ms,,,,,,127.0.0.1,57817,,F,,mosn_cluster=cb7aaa7a31fd0490299c530e83a868231&mosn_namespace=default&mosn_log=true&,426.311µs
*/
func (s *SofaRPCSpan) clientHttpLogger() error {
	printData := buffer.NewIoBuffer(512)
	//1 日志打印时间:2021-02-05 13:46:44.264
	date := s.endTime.Format("2006-01-02 15:04:05.000")
	printData.WriteString(date + ",")
	//2 当前应用名:echo-server
	printData.WriteString(s.appName + ",")
	//3 TraceId:0ad55acc1512711661769715320278
	printData.WriteString(s.tags[generator.TRACE_ID] + ",")
	//4 RpcId:0
	printData.WriteString(s.tags[generator.SPAN_ID] + ",")
	//5 host
	printData.WriteString(s.tags[generator.DOWNSTEAM_HOST_ADDRESS] + ",")
	//6 uri
	printData.WriteString(s.tags[generator.REQUEST_URL] + ",")
	//7 协议:bolt
	printData.WriteString(s.tags[generator.PROTOCOL] + ",")
	//8 调用方式:SYNC
	printData.WriteString("SYNC" + ",")
	//9 目标地址:30.73.162.190:12200
	printData.WriteString(s.tags[generator.UPSTREAM_HOST_ADDRESS] + ",")
	//10 目标系统名:echo-server
	printData.WriteString(s.tags[generator.TARGET_APP_NAME] + ",")
	//11 目标Zone:GZ00B
	printData.WriteString(s.tags[generator.TARGET_CELL] + ",")
	//12 目标IDC:
	printData.WriteString(s.tags[generator.TARGET_IDC] + ",")
	//13 目标City:
	printData.WriteString(s.tags[generator.TARGET_CITY] + ",")
	//14 uid:
	printData.WriteString(s.tags[generator.UID] + ",")
	//15 status
	s.status(printData)
	//16 请求大小:388B
	printData.WriteString(s.tags[generator.REQUEST_SIZE] + ",")
	//17 响应大小:105B
	printData.WriteString(s.tags[generator.RESPONSE_SIZE] + ",")
	//18 调用耗时:2ms
	printData.WriteString(s.tags[generator.DURATION] + ",")
	//19 链接建立耗时:0ms
	printData.WriteString("0ms,")
	//20 请求序列化耗时:0ms
	printData.WriteString("0ms,")
	//21 超时参考耗时:0ms
	printData.WriteString("0ms,")
	//22 当前线程名:
	printData.WriteString(",")
	//23 路由记录:
	printData.WriteString(s.tags[generator.ROUTE_RECORD] + ",")
	//24 弹性数据位:
	printData.WriteString(",")
	//25 是否需要弹性:
	printData.WriteString(",")
	//26 转发的服务名称:
	printData.WriteString(s.tags[generator.CALLER_APP_NAME] + ",")
	//27 clientip: 127.0.0.1
	downStreamHostAddress := strings.Split(s.tags[generator.DOWNSTEAM_HOST_ADDRESS], ":")
	if len(downStreamHostAddress) > 0 {
		localIp := strings.Split(s.tags[generator.DOWNSTEAM_HOST_ADDRESS], ":")[0]
		printData.WriteString(localIp + ",")
	} else {
		printData.WriteString(",")
	}
	//28 clientport: 61181
	if len(downStreamHostAddress) > 1 {
		localPort := strings.Split(s.tags[generator.DOWNSTEAM_HOST_ADDRESS], ":")[1]
		printData.WriteString(localPort + ",")
	} else {
		printData.WriteString(",")
	}
	//29 当前zone:
	printData.WriteString(",")
	//30 是否物理机器: F
	printData.WriteString(s.isPod() + ",")
	//31 systemMap:
	printData.WriteString(",")
	//32 系统穿透数据
	printData.WriteString(s.tags[generator.BAGGAGE_DATA] + ",")
	//33 mosn处理时间
	printData.WriteString(s.tags[generator.MOSN_PROCESS_TIME])
	printData.WriteString("\n")
	return s.egressLogger.Print(printData, true)
}

/*
springcloud-server-digest.log
2021-02-05 11:56:30.962,reservation-service,1e49a2be1612497390960139868886,0,127.0.0.1:10088,/echo/name/aaa,HTTP,,30.73.162.190:57818,reservation-service,,,1ms,0ms,,00,,,0,,mosn_tls_state=off&mosn_cluster=cb7aaa7a31fd0490299c530e83a868231&mosn_namespace=default&mosn_log=true&,262.787µs,0B,16B
*/
func (s *SofaRPCSpan) serverHttpLogger() error {
	printData := buffer.NewIoBuffer(512)
	//1.日志打印时间
	date := s.endTime.Format("2006-01-02 15:04:05.000")
	printData.WriteString(date + ",")
	//2当前应用名
	printData.WriteString(s.appName + ",")
	//3TraceId
	printData.WriteString(s.tags[generator.TRACE_ID] + ",")
	//4RpcId
	printData.WriteString(s.tags[generator.SPAN_ID] + ",")
	//5host
	printData.WriteString(s.tags[generator.DOWNSTEAM_HOST_ADDRESS] + ",")
	//6方法名
	printData.WriteString(s.tags[generator.METHOD_NAME] + ",")
	//7协议
	printData.WriteString(s.tags[generator.PROTOCOL] + ",")
	//8调用方式
	printData.WriteString("SYNC" + ",")
	//9调用者 URL
	printData.WriteString(s.tags[generator.REQUEST_URL] + ",")
	//10调用者应用名
	printData.WriteString(s.tags[generator.CALLER_APP_NAME] + ",")
	//11 目标Zone:GZ00B
	printData.WriteString(s.tags[generator.TARGET_CELL] + ",")
	//12 目标IDC:
	printData.WriteString(s.tags[generator.TARGET_IDC] + ",")
	//13请求处理耗时（ms）
	printData.WriteString(s.tags[generator.DURATION] + ",")
	//14服务端响应序列化耗时（ms）
	printData.WriteString(s.tags[generator.UPSTREAM_DURATION] + ",")
	//15当前线程名
	printData.WriteString(",")
	s.status(printData)
	//17 beElasticServiceName（表明这次调用是转发调用,转发的服务名称和方法名称是啥值如：“com.test.service.testservice.TestService:1.0:biztest---doProcess”）
	printData.WriteString(",")
	//18 beElastic（表示没有被转发的处理）
	printData.WriteString(",")
	//19 rpc线程池等待时间
	printData.WriteString(",")
	//20 系统穿透数据（kv 格式，用于传送系统灾备信息等）
	printData.WriteString(",")
	//21 穿透数据（kv格式）
	printData.WriteString(s.tags[generator.BAGGAGE_DATA] + ",")
	//22 mosn处理时间
	printData.WriteString(s.tags[generator.MOSN_PROCESS_TIME] + ",")
	//23 请求大小:388B
	printData.WriteString(s.tags[generator.REQUEST_SIZE] + ",")
	//24 响应大小:105B
	printData.WriteString(s.tags[generator.RESPONSE_SIZE])

	printData.WriteString("\n")
	return s.ingressLogger.Print(printData, true)
}
