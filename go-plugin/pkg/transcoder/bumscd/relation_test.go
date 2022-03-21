package bumscd

import (
	"encoding/json"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestNewRelation(t *testing.T) {
	tests := []struct {
		name    string
		head    string
		resp    bool
		body    string
		want    error
		wantstr string
	}{
		{
			name:    "normal2",
			head:    `<?xml version="1.0" encoding="UTF-8" ?><root><SYS_HEAD><TRAN_TIMESTAMP type="string" length="9" positionOrDefault="head@tranTimestamp@" scale="0" /><MODULE_ID  type="string" length="4" positionOrDefault="@@RB" scale="0" /><SERVICE_SCENE type="string" length="2" positionOrDefault="@@03" scale="0" /><SOURCE_BRANCH_NO type="string" length="10" positionOrDefault="@@EsbBJFront" scale="0" /><TRAN_DATE type="string" length="8" positionOrDefault="head@tranDate@" scale="0" /><CONSUMER_ID type="string" length="6" positionOrDefault="headers@OrigSender@QDT001" scale="0" /><CONSUMER_SEQ_NO type="string" length="52" positionOrDefault="head@consumerSeqNo@" scale="0" /><CONSUMER_SVR_ID type="string" length="8" positionOrDefault="@@SmartUNFront" scale="0" /><TRAN_CODE type="string" length="4" positionOrDefault="head@tranCode@" scale="0" /><SERVICE_CODE type="string" length="30" positionOrDefault="@@ECIF1200003000" scale="0" /><MESSAGE_CODE type="string" length="6" positionOrDefault="@@0008" scale="0" /></SYS_HEAD><APP_HEAD><USER_ID type="string" length="30" positionOrDefault="@@BOBQZ2" scale="0" /><AGENT_BRANCH_ID type="string" length="9" positionOrDefault="head@branchId@" scale="0" /><BRANCH_ID type="string" length="9" positionOrDefault="head@branchId@" scale="0" /></APP_HEAD><LOCAL_HEAD></LOCAL_HEAD></root>`,
			body:    `<?xml version="1.0" encoding="UTF-8" ?><root><cardOrAcctNo type="string" length="19" convert="I_CDNO" cnName="卡号/账号"/><flagChannel type="int" length="3" convert="I_CHAF1" cnName="渠道标志"/><flagChannel2 type="int" length="3" convert="I_CHAF2" cnName="渠道标志2"/><password type="string" length="16" convert="I_PSWD" cnName="密码"/><flagPwdChk type="double" length="5" scale="3" convert="I_YMFG1" cnName="验密标志"/><flagPwdChk2 type="double" length="5" scale="3" convert="I_YMFG2" cnName="验密标志"/><flagPwdChk3 type="double" length="5" scale="3" convert="I_YMFG3" cnName="验密标志"/><ccy type="list" convert="I_CRNO" cnName="币种"><dateDue type="string" length="8" scale="0" convert="O___ENDDT" cnName="有效截止日期"/><nameClient type="string" length="200" convert="O___NAME" cnName="姓名"/><testArray type= "list" convert="TEST_ARRAY"><testTypeBusi type="double" length="1" scale="0" convert="TEST___MDMLX" cnName="面对面类型"/><testCodeProduct type="double" length="2" scale="1" convert="TEST___PGCP" cnName="评估产品"/></testArray><idGlobal type="string" length="40" convert="O___IDNO" cnName="证件号码"/></ccy></root>`,
			wantstr: `{"sys_head":[{"cd_key":"TRAN_TIMESTAMP","type":"string","length":"9","scale":"0","source":"head","bums_key":"tranTimestamp"},{"cd_key":"MODULE_ID","type":"string","length":"4","scale":"0","default":"RB"},{"cd_key":"SERVICE_SCENE","type":"string","length":"2","scale":"0","default":"03"},{"cd_key":"SOURCE_BRANCH_NO","type":"string","length":"10","scale":"0","default":"EsbBJFront"},{"cd_key":"TRAN_DATE","type":"string","length":"8","scale":"0","source":"head","bums_key":"tranDate"},{"cd_key":"CONSUMER_ID","type":"string","length":"6","scale":"0","source":"headers","bums_key":"OrigSender","default":"QDT001"},{"cd_key":"CONSUMER_SEQ_NO","type":"string","length":"52","scale":"0","source":"head","bums_key":"consumerSeqNo"},{"cd_key":"CONSUMER_SVR_ID","type":"string","length":"8","scale":"0","default":"SmartUNFront"},{"cd_key":"TRAN_CODE","type":"string","length":"4","scale":"0","source":"head","bums_key":"tranCode"},{"cd_key":"SERVICE_CODE","type":"string","length":"30","scale":"0","default":"ECIF1200003000"},{"cd_key":"MESSAGE_CODE","type":"string","length":"6","scale":"0","default":"0008"}],"app_head":[{"cd_key":"USER_ID","type":"string","length":"30","scale":"0","default":"BOBQZ2"},{"cd_key":"AGENT_BRANCH_ID","type":"string","length":"9","scale":"0","source":"head","bums_key":"branchId"},{"cd_key":"BRANCH_ID","type":"string","length":"9","scale":"0","source":"head","bums_key":"branchId"}],"local_head":[],"body":{"cardOrAcctNo":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"ccy":{"list_iterms":{"dateDue":{"cd_key":"O@ENDDT","type":"string","length":"8","scale":"0","describes":"有效截止日期","bums_key":"dateDue"},"idGlobal":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"nameClient":{"cd_key":"O@NAME","type":"string","length":"200","scale":"0","describes":"姓名","bums_key":"nameClient"},"testArray":{"list_iterms":{"testCodeProduct":{"cd_key":"TEST@PGCP","type":"double","length":"2","scale":"1","describes":"评估产品","bums_key":"testCodeProduct"},"testTypeBusi":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"},"flagChannel":{"cd_key":"I_CHAF1","type":"int","length":"3","scale":"0","describes":"渠道标志","bums_key":"flagChannel"},"flagChannel2":{"cd_key":"I_CHAF2","type":"int","length":"3","scale":"0","describes":"渠道标志2","bums_key":"flagChannel2"},"flagPwdChk":{"cd_key":"I_YMFG1","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk"},"flagPwdChk2":{"cd_key":"I_YMFG2","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk2"},"flagPwdChk3":{"cd_key":"I_YMFG3","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk3"},"password":{"cd_key":"I_PSWD","type":"string","length":"16","scale":"0","describes":"密码","bums_key":"password"}}}`,
		},
		{
			resp:    true,
			name:    "normal-respone",
			head:    `<?xml version="1.0" encoding="UTF-8" ?><root><SYS_HEAD><RET_SEQ_NO convert="requestSequenceNo"/><RET convert="ret" type="list"><RET_CODE convert="retCode" /><RET_MSG convert="retMsg" /></RET><RESP_CODE convert="responseCode" /><RET_SYS_TIME convert="requestSysTime" /><SERVICE_SCENE convert="serviceScene" /><MESSAGE_CODE convert="messageCode" /><SERVICE_CODE convert="serviceCode" /><MESSAGE_TYPE convert="messageType" /><RET_BUSI_DATE convert="requestBusiData" /><RET_STATUS convert="requestStauts" /><RET_SYS_DATE convert="requestSysData" /></SYS_HEAD><APP_HEAD><USER_ID convert="userId" /><AGENT_BRANCH_ID convert="" /><BRANCH_ID /></APP_HEAD><LOCAL_HEAD></LOCAL_HEAD></root>`,
			body:    `<?xml version="1.0" encoding="UTF-8" ?><root><I_CDNO type="string" length="19" convert="cardOrAcctNo" cnName="卡号/账号"/><I_CHAF1 type="int" length="3" convert="flagChannel" cnName="渠道标志"/><I_CHAF2 type="int" length="3" convert="flagChannel2" cnName="渠道标志2"/><I_PSWD type="string" length="16" convert="password" cnName="密码"/><I_YMFG1 type="double" length="5" scale="3" convert="flagPwdChk" cnName="验密标志"/><I_YMFG2 type="double" length="5" scale="3" convert="flagPwdChk2" cnName="验密标志"/><I_YMFG3 type="double" length="5" scale="3" convert="flagPwdChk3" cnName="验密标志"/><I_CRNO type="list" convert="ccy" cnName="币种"><O___ENDDT type="string" length="8" scale="0" convert="dateDue" cnName="有效截止日期"/><O___NAME type="string" length="200" convert="nameClient" cnName="姓名"/><TEST_ARRAY type= "list" convert="testArray"><TEST___MDMLX type="double" length="1" scale="0" convert="testTypeBusi" cnName="面对面类型"/><TEST___PGCP type="double" length="2" scale="1" convert="testCodeProduct" cnName="评估产品"/></TEST_ARRAY><O___IDNO type="string" length="40" convert="idGlobal" cnName="证件号码"/></I_CRNO></root>`,
			wantstr: `{"sys_head":[{"cd_key":"RET_SEQ_NO","type":"string","bums_key":"requestSequenceNo"},{"head_iterms":[{"cd_key":"RET_CODE","type":"string","bums_key":"retCode"},{"cd_key":"RET_MSG","type":"string","bums_key":"retMsg"}],"cd_key":"RET","type":"list","bums_key":"ret"},{"cd_key":"RESP_CODE","type":"string","bums_key":"responseCode"},{"cd_key":"RET_SYS_TIME","type":"string","bums_key":"requestSysTime"},{"cd_key":"SERVICE_SCENE","type":"string","bums_key":"serviceScene"},{"cd_key":"MESSAGE_CODE","type":"string","bums_key":"messageCode"},{"cd_key":"SERVICE_CODE","type":"string","bums_key":"serviceCode"},{"cd_key":"MESSAGE_TYPE","type":"string","bums_key":"messageType"},{"cd_key":"RET_BUSI_DATE","type":"string","bums_key":"requestBusiData"},{"cd_key":"RET_STATUS","type":"string","bums_key":"requestStauts"},{"cd_key":"RET_SYS_DATE","type":"string","bums_key":"requestSysData"}],"app_head":[{"cd_key":"USER_ID","type":"string","bums_key":"userId"},{"cd_key":"AGENT_BRANCH_ID","type":"string","bums_key":"agentBranchId"},{"cd_key":"BRANCH_ID","type":"string","bums_key":"branchId"}],"local_head":[],"body":{"I_CDNO":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"I_CHAF1":{"cd_key":"I_CHAF1","type":"int","length":"3","scale":"0","describes":"渠道标志","bums_key":"flagChannel"},"I_CHAF2":{"cd_key":"I_CHAF2","type":"int","length":"3","scale":"0","describes":"渠道标志2","bums_key":"flagChannel2"},"I_CRNO":{"list_iterms":{"O@ENDDT":{"cd_key":"O@ENDDT","type":"string","length":"8","scale":"0","describes":"有效截止日期","bums_key":"dateDue"},"O@IDNO":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"O@NAME":{"cd_key":"O@NAME","type":"string","length":"200","scale":"0","describes":"姓名","bums_key":"nameClient"},"TEST_ARRAY":{"list_iterms":{"TEST@MDMLX":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"},"TEST@PGCP":{"cd_key":"TEST@PGCP","type":"double","length":"2","scale":"1","describes":"评估产品","bums_key":"testCodeProduct"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"},"I_PSWD":{"cd_key":"I_PSWD","type":"string","length":"16","scale":"0","describes":"密码","bums_key":"password"},"I_YMFG1":{"cd_key":"I_YMFG1","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk"},"I_YMFG2":{"cd_key":"I_YMFG2","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk2"},"I_YMFG3":{"cd_key":"I_YMFG3","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk3"}}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRelation(tt.resp)
			err := got.ParseString(tt.head, tt.body)
			assert.Equal(t, err, tt.want)
			str, _ := json.Marshal(got)
			assert.Len(t, str, len(tt.wantstr))
			t.Logf("%s", str)
		})
	}
}

func TestParseBody(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
		resp bool
	}{
		{
			name: "normal-field",
			body: `<?xml version="1.0" encoding="UTF-8" ?><root><cardOrAcctNo type="string" length="19" convert="I_CDNO" cnName="卡号/账号"/><idGlobal type="string" length="40" convert="O___IDNO" cnName="证件号码"/>></root>`,
			want: `{"cardOrAcctNo":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"idGlobal":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"}}`,
		},
		{
			name: "normal-array",
			body: `<?xml version="1.0" encoding="UTF-8" ?><root><ccy type="list" convert="I_CRNO" cnName="币种"><testArray type= "list" convert="TEST_ARRAY"><testTypeBusi type="double" length="1" scale="0" convert="TEST___MDMLX" cnName="面对面类型"/></testArray><idGlobal type="string" length="40" convert="O___IDNO" cnName="证件号码"/></ccy></root>`,
			want: `{"ccy":{"list_iterms":{"idGlobal":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"testArray":{"list_iterms":{"testTypeBusi":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"}}`,
		},
		{
			name: "normal-respone",
			body: `<?xml version="1.0" encoding="UTF-8" ?><root><I_CDNO type="string" length="19" convert="cardOrAcctNo" cnName="卡号/账号"/><I_CHAF1 type="int" length="3" convert="flagChannel" cnName="渠道标志"/><I_CHAF2 type="int" length="3" convert="flagChannel2" cnName="渠道标志2"/><I_PSWD type="string" length="16" convert="password" cnName="密码"/><I_YMFG1 type="double" length="5" scale="3" convert="flagPwdChk" cnName="验密标志"/><I_YMFG2 type="double" length="5" scale="3" convert="flagPwdChk2" cnName="验密标志"/><I_YMFG3 type="double" length="5" scale="3" convert="flagPwdChk3" cnName="验密标志"/><I_CRNO type="list" convert="ccy" cnName="币种"><O___ENDDT type="string" length="8" scale="0" convert="dateDue" cnName="有效截止日期"/><O___NAME type="string" length="200" convert="nameClient" cnName="姓名"/><TEST_ARRAY type= "list" convert="testArray"><TEST___MDMLX type="double" length="1" scale="0" convert="testTypeBusi" cnName="面对面类型"/><TEST___PGCP type="double" length="2" scale="1" convert="testCodeProduct" cnName="评估产品"/></TEST_ARRAY><O___IDNO type="string" length="40" convert="idGlobal" cnName="证件号码"/></I_CRNO></root>`,
			want: `{"I_CDNO":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"I_CHAF1":{"cd_key":"I_CHAF1","type":"int","length":"3","scale":"0","describes":"渠道标志","bums_key":"flagChannel"},"I_CHAF2":{"cd_key":"I_CHAF2","type":"int","length":"3","scale":"0","describes":"渠道标志2","bums_key":"flagChannel2"},"I_CRNO":{"list_iterms":{"O@ENDDT":{"cd_key":"O@ENDDT","type":"string","length":"8","scale":"0","describes":"有效截止日期","bums_key":"dateDue"},"O@IDNO":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"O@NAME":{"cd_key":"O@NAME","type":"string","length":"200","scale":"0","describes":"姓名","bums_key":"nameClient"},"TEST_ARRAY":{"list_iterms":{"TEST@MDMLX":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"},"TEST@PGCP":{"cd_key":"TEST@PGCP","type":"double","length":"2","scale":"1","describes":"评估产品","bums_key":"testCodeProduct"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"},"I_PSWD":{"cd_key":"I_PSWD","type":"string","length":"16","scale":"0","describes":"密码","bums_key":"password"},"I_YMFG1":{"cd_key":"I_YMFG1","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk"},"I_YMFG2":{"cd_key":"I_YMFG2","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk2"},"I_YMFG3":{"cd_key":"I_YMFG3","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk3"}}`,
			resp: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRelation(tt.resp)
			root := etree.NewDocument()
			root.ReadFromString(tt.body)
			got.ParseBody(root)
			str, _ := json.Marshal(got.Body)
			assert.Equal(t, string(str), tt.want)
		})
	}
}

func TestParseHead(t *testing.T) {
	tests := []struct {
		name  string
		head  string
		sys   string
		app   string
		resp  bool
		local string
	}{
		{
			name:  "normal-field",
			head:  `<?xml version="1.0" encoding="UTF-8" ?><root><SYS_HEAD><TRAN_TIMESTAMP type="string" length="9" positionOrDefault="head@tranTimestamp@" scale="0" /><MODULE_ID  type="string" length="4" positionOrDefault="@@RB" scale="0" /><SERVICE_SCENE type="string" length="2" positionOrDefault="@@03" scale="0" /><SOURCE_BRANCH_NO type="string" length="10" positionOrDefault="@@EsbBJFront" scale="0" /><TRAN_DATE type="string" length="8" positionOrDefault="head@tranDate@" scale="0" /><CONSUMER_ID type="string" length="6" positionOrDefault="headers@OrigSender@QDT001" scale="0" /><CONSUMER_SEQ_NO type="string" length="52" positionOrDefault="head@consumerSeqNo@" scale="0" /><CONSUMER_SVR_ID type="string" length="8" positionOrDefault="@@SmartUNFront" scale="0" /><TRAN_CODE type="string" length="4" positionOrDefault="head@tranCode@" scale="0" /><SERVICE_CODE type="string" length="30" positionOrDefault="@@ECIF1200003000" scale="0" /><MESSAGE_CODE type="string" length="6" positionOrDefault="@@0008" scale="0" /></SYS_HEAD><APP_HEAD><USER_ID type="string" length="30" positionOrDefault="@@BOBQZ2" scale="0" /><AGENT_BRANCH_ID type="string" length="9" positionOrDefault="head@branchId@" scale="0" /><BRANCH_ID type="string" length="9" positionOrDefault="head@branchId@" scale="0" /></APP_HEAD><LOCAL_HEAD></LOCAL_HEAD></root>`,
			sys:   `[{"cd_key":"TRAN_TIMESTAMP","type":"string","length":"9","scale":"0","source":"head","bums_key":"tranTimestamp"},{"cd_key":"MODULE_ID","type":"string","length":"4","scale":"0","default":"RB"},{"cd_key":"SERVICE_SCENE","type":"string","length":"2","scale":"0","default":"03"},{"cd_key":"SOURCE_BRANCH_NO","type":"string","length":"10","scale":"0","default":"EsbBJFront"},{"cd_key":"TRAN_DATE","type":"string","length":"8","scale":"0","source":"head","bums_key":"tranDate"},{"cd_key":"CONSUMER_ID","type":"string","length":"6","scale":"0","source":"headers","bums_key":"OrigSender","default":"QDT001"},{"cd_key":"CONSUMER_SEQ_NO","type":"string","length":"52","scale":"0","source":"head","bums_key":"consumerSeqNo"},{"cd_key":"CONSUMER_SVR_ID","type":"string","length":"8","scale":"0","default":"SmartUNFront"},{"cd_key":"TRAN_CODE","type":"string","length":"4","scale":"0","source":"head","bums_key":"tranCode"},{"cd_key":"SERVICE_CODE","type":"string","length":"30","scale":"0","default":"ECIF1200003000"},{"cd_key":"MESSAGE_CODE","type":"string","length":"6","scale":"0","default":"0008"}]`,
			app:   `[{"cd_key":"USER_ID","type":"string","length":"30","scale":"0","default":"BOBQZ2"},{"cd_key":"AGENT_BRANCH_ID","type":"string","length":"9","scale":"0","source":"head","bums_key":"branchId"},{"cd_key":"BRANCH_ID","type":"string","length":"9","scale":"0","source":"head","bums_key":"branchId"}]`,
			resp:  false,
			local: `[]`,
		},
		{
			name:  "normal-list-request",
			resp:  false,
			local: `[]`,
			app:   `[]`,
			sys:   `[{"cd_key":"TRAN_DATE","type":"string","length":"8","scale":"0","source":"head","bums_key":"tranDate"},{"head_iterms":[{"cd_key":"SUBH","type":"string","length":"6","scale":"0","source":"head","bums_key":"get#sub"}],"cd_key":"LISTH","type":"list","source":"head","bums_key":"listh"}]`,
			head:  `<?xml version="1.0" encoding="UTF-8" ?><root><SYS_HEAD><TRAN_DATE type="string" length="8" positionOrDefault="head@tranDate@" scale="0" /><LISTH type="list" positionOrDefault="head@listh@"><SUBH type="string" length="6" positionOrDefault="head@get#sub@" scale="0" /></LISTH></SYS_HEAD><APP_HEAD></APP_HEAD><LOCAL_HEAD></LOCAL_HEAD></root>`,
		},
		{
			name:  "normal-list",
			head:  `<?xml version="1.0" encoding="UTF-8" ?><root><SYS_HEAD><RET_SEQ_NO convert="requestSequenceNo"/><RET convert="ret" type="list"><RET_CODE convert="retCode" /><RET_MSG convert="retMsg" /></RET><RESP_CODE convert="responseCode" /><RET_SYS_TIME convert="requestSysTime" /><SERVICE_SCENE convert="serviceScene" /><MESSAGE_CODE convert="messageCode" /><SERVICE_CODE convert="serviceCode" /><MESSAGE_TYPE convert="messageType" /><RET_BUSI_DATE convert="requestBusiData" /><RET_STATUS convert="requestStauts" /><RET_SYS_DATE convert="requestSysData" /></SYS_HEAD><APP_HEAD><USER_ID convert="userId" /><AGENT_BRANCH_ID convert="" /><BRANCH_ID /></APP_HEAD><LOCAL_HEAD></LOCAL_HEAD></root>`,
			resp:  true,
			sys:   `[{"cd_key":"RET_SEQ_NO","type":"string","bums_key":"requestSequenceNo"},{"head_iterms":[{"cd_key":"RET_CODE","type":"string","bums_key":"retCode"},{"cd_key":"RET_MSG","type":"string","bums_key":"retMsg"}],"cd_key":"RET","type":"list","bums_key":"ret"},{"cd_key":"RESP_CODE","type":"string","bums_key":"responseCode"},{"cd_key":"RET_SYS_TIME","type":"string","bums_key":"requestSysTime"},{"cd_key":"SERVICE_SCENE","type":"string","bums_key":"serviceScene"},{"cd_key":"MESSAGE_CODE","type":"string","bums_key":"messageCode"},{"cd_key":"SERVICE_CODE","type":"string","bums_key":"serviceCode"},{"cd_key":"MESSAGE_TYPE","type":"string","bums_key":"messageType"},{"cd_key":"RET_BUSI_DATE","type":"string","bums_key":"requestBusiData"},{"cd_key":"RET_STATUS","type":"string","bums_key":"requestStauts"},{"cd_key":"RET_SYS_DATE","type":"string","bums_key":"requestSysData"}]`,
			app:   `[{"cd_key":"USER_ID","type":"string","bums_key":"userId"},{"cd_key":"AGENT_BRANCH_ID","type":"string","bums_key":"agentBranchId"},{"cd_key":"BRANCH_ID","type":"string","bums_key":"branchId"}]`,
			local: `[]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRelation(tt.resp)
			root := etree.NewDocument()
			root.ReadFromString(tt.head)
			got.ParseHead(root)

			str, _ := json.Marshal(got.SysHead)
			assert.Equal(t, string(str), tt.sys)

			str, _ = json.Marshal(got.AppHead)
			assert.Equal(t, string(str), tt.app)

			str, _ = json.Marshal(got.LocalHead)
			assert.Equal(t, string(str), tt.local)
		})
	}
}
