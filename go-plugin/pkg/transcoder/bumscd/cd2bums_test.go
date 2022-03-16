package bumscd

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

func TestCd2Bums_HeadInBody(t *testing.T) {
	type fields struct {
		config string
		header api.HeaderMap
		buf    api.IoBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "head",
			fields: fields{
				config: `{"sys_head":[{"cd_key":"RET_SEQ_NO","type":"string","bums_key":"requestSequenceNo"},{"head_iterms":[{"cd_key":"RET_CODE","type":"string","bums_key":"retCode"},{"cd_key":"RET_MSG","type":"string","bums_key":"retMsg"}],"cd_key":"RET","type":"list","bums_key":"ret"},{"cd_key":"RESP_CODE","type":"string","bums_key":"responseCode"},{"cd_key":"RET_SYS_TIME","type":"string","bums_key":"requestSysTime"},{"cd_key":"SERVICE_SCENE","type":"string","bums_key":"serviceScene"},{"cd_key":"MESSAGE_CODE","type":"string","bums_key":"messageCode"},{"cd_key":"SERVICE_CODE","type":"string","bums_key":"serviceCode"},{"cd_key":"MESSAGE_TYPE","type":"string","bums_key":"messageType"},{"cd_key":"RET_BUSI_DATE","type":"string","bums_key":"requestBusiData"},{"cd_key":"RET_STATUS","type":"string","bums_key":"requestStauts"},{"cd_key":"RET_SYS_DATE","type":"string","bums_key":"requestSysData"}],"app_head":[{"cd_key":"USER_ID","type":"string","bums_key":"userId"},{"cd_key":"AGENT_BRANCH_ID","type":"string","bums_key":"agentBranchId"},{"cd_key":"BRANCH_ID","type":"string","bums_key":"branchId"}],"local_head":[],"body":{}}`,
				buf:    buffer.NewIoBufferString(`<?xml version="1.0" encoding="UTF-8"?><service><sys-header><data name="SYS_HEAD"><struct><data name="RET"><array><struct><data name="RET_CODE"><field length="6" scale="0" type="string">999999</field></data><data name="RET_MSG"><field length="9" scale="0" type="string">JDBC调用失败!</field></data></struct></array></data><data name="DEST_BRANCH_NO"><field length="0" scale="0" type="string"/></data><data name="SEQ_NO"><field length="19" scale="0" type="string">ANIU001000000218782</field></data><data name="MESSAGE_CODE"><field length="4" scale="0" type="string">9527</field></data><data name="SERVICE_CODE"><field length="14" scale="0" type="string">BRNC1200952700</field></data><data name="MESSAGE_TYPE"><field length="4" scale="0" type="string">1210</field></data><data name="RET_STATUS"><field length="1" scale="0" type="string">F</field></data><data name="TRAN_TIMESTAMP"><field length="6" scale="0" type="string">135519</field></data><data name="SOURCE_BRANCH_NO"><field length="10" scale="0" type="string">EsbBJFront</field></data><data name="FILE_PATH"><field length="0" scale="0" type="string"/></data><data name="TRAN_DATE"><field length="8" scale="0" type="string">20220113</field></data><data name="BRANCH_ID"><field length="0" scale="0" type="string"/></data></struct></data></sys-header><app-header><data name="APP_HEAD"><struct><data name="AGENT_BRANCH_ID"><field length="9" scale="0" type="string">00301</field></data><data name="USER_ID"><field length="30" scale="0" type="string">BOBQZ2</field></data><data name="BRANCH_ID"><field length="9" scale="0" type="string">00301</field></data></struct></data></app-header><local-header><data name="LOCAL_HEAD"><struct/></data></local-header><body/></service>`),
			},
			want:    `{"agentBranchId":"00301","branchId":"00301","destBranchNo":"","filePath":"","messageCode":"9527","messageType":"1210","requestStauts":"F","retCode":"999999","retMsg":"JDBC调用失败!","seqNo":"ANIU001000000218782","serviceCode":"BRNC1200952700","sourceBranchNo":"EsbBJFront","tranDate":"20220113","tranTimestamp":"135519","userId":"BOBQZ2"}`,
			wantErr: nil,
		},
		{
			name: "head-list",
			fields: fields{
				config: `{"sys_head":[{"cd_key":"RET_SEQ_NO","type":"string","bums_key":"requestSequenceNo"},{"head_iterms":[{"cd_key":"RET_CODE","type":"string","bums_key":"retCode"},{"cd_key":"RET_MSG","type":"string","bums_key":"retMsg"}],"cd_key":"RET_LIST","type":"list","bums_key":"ret_list"},{"cd_key":"RESP_CODE","type":"string","bums_key":"responseCode"},{"cd_key":"RET_SYS_TIME","type":"string","bums_key":"requestSysTime"},{"cd_key":"SERVICE_SCENE","type":"string","bums_key":"serviceScene"},{"cd_key":"MESSAGE_CODE","type":"string","bums_key":"messageCode"},{"cd_key":"SERVICE_CODE","type":"string","bums_key":"serviceCode"},{"cd_key":"MESSAGE_TYPE","type":"string","bums_key":"messageType"},{"cd_key":"RET_BUSI_DATE","type":"string","bums_key":"requestBusiData"},{"cd_key":"RET_STATUS","type":"string","bums_key":"requestStauts"},{"cd_key":"RET_SYS_DATE","type":"string","bums_key":"requestSysData"}],"app_head":[{"cd_key":"USER_ID","type":"string","bums_key":"userId"},{"cd_key":"AGENT_BRANCH_ID","type":"string","bums_key":"agentBranchId"},{"cd_key":"BRANCH_ID","type":"string","bums_key":"branchId"}],"local_head":[],"body":{}}`,
				buf:    buffer.NewIoBufferString(`<?xml version="1.0" encoding="UTF-8"?><service><sys-header><data name="SYS_HEAD"><struct><data name="RET_LIST"><array><struct><data name="RET_CODE"><field length="6" scale="0" type="string">999999</field></data><data name="RET_MSG"><field length="9" scale="0" type="string">JDBC调用失败!</field></data></struct></array></data><data name="DEST_BRANCH_NO"><field length="0" scale="0" type="string"/></data><data name="SEQ_NO"><field length="19" scale="0" type="string">ANIU001000000218782</field></data><data name="MESSAGE_CODE"><field length="4" scale="0" type="string">9527</field></data><data name="SERVICE_CODE"><field length="14" scale="0" type="string">BRNC1200952700</field></data><data name="MESSAGE_TYPE"><field length="4" scale="0" type="string">1210</field></data><data name="RET_STATUS"><field length="1" scale="0" type="string">F</field></data><data name="TRAN_TIMESTAMP"><field length="6" scale="0" type="string">135519</field></data><data name="SOURCE_BRANCH_NO"><field length="10" scale="0" type="string">EsbBJFront</field></data><data name="FILE_PATH"><field length="0" scale="0" type="string"/></data><data name="TRAN_DATE"><field length="8" scale="0" type="string">20220113</field></data><data name="BRANCH_ID"><field length="0" scale="0" type="string"/></data></struct></data></sys-header><app-header><data name="APP_HEAD"><struct><data name="AGENT_BRANCH_ID"><field length="9" scale="0" type="string">00301</field></data><data name="USER_ID"><field length="30" scale="0" type="string">BOBQZ2</field></data><data name="BRANCH_ID"><field length="9" scale="0" type="string">00301</field></data></struct></data></app-header><local-header><data name="LOCAL_HEAD"><struct/></data></local-header><body/></service>`),
			},
			want:    `{"ret_list":[{"retCode":"999999","retMsg":"JDBC调用失败!"}],"destBranchNo":"","seqNo":"ANIU001000000218782","messageCode":"9527","serviceCode":"BRNC1200952700","messageType":"1210","requestStauts":"F","tranTimestamp":"135519","sourceBranchNo":"EsbBJFront","filePath":"","tranDate":"20220113","branchId":"00301","agentBranchId":"00301","userId":"BOBQZ2"}`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var relation Relation
			json.Unmarshal([]byte(tt.fields.config), &relation)
			c2b, err := NewCd2Bums(tt.fields.header, tt.fields.buf, &relation)
			assert.NoError(t, err)
			got, err := c2b.HeadInBody()
			assert.Equal(t, err, tt.wantErr)
			assert.Len(t, got.String(), len(tt.want))
		})
	}
}

func TestCd2Bums_BodyInBody(t *testing.T) {
	type fields struct {
		config string
		header api.HeaderMap
		buf    api.IoBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "head-array",
			fields: fields{
				config: `{"sys_head":[],"app_head":[],"local_head":[],"body":{"I_CDNO":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"I_CHAF1":{"cd_key":"I_CHAF1","type":"int","length":"3","scale":"0","describes":"渠道标志","bums_key":"flagChannel"},"I_CHAF2":{"cd_key":"I_CHAF2","type":"int","length":"3","scale":"0","describes":"渠道标志2","bums_key":"flagChannel2"},"I_CRNO":{"list_iterms":{"O@ENDDT":{"cd_key":"O@ENDDT","type":"string","length":"8","scale":"0","describes":"有效截止日期","bums_key":"dateDue"},"O@IDNO":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"O@NAME":{"cd_key":"O@NAME","type":"string","length":"200","scale":"0","describes":"姓名","bums_key":"nameClient"},"TEST_ARRAY":{"list_iterms":{"TEST@MDMLX":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"},"TEST@PGCP":{"cd_key":"TEST@PGCP","type":"double","length":"2","scale":"1","describes":"评估产品","bums_key":"testCodeProduct"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"},"I_PSWD":{"cd_key":"I_PSWD","type":"string","length":"16","scale":"0","describes":"密码","bums_key":"password"},"I_YMFG1":{"cd_key":"I_YMFG1","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk"},"I_YMFG2":{"cd_key":"I_YMFG2","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk2"},"I_YMFG3":{"cd_key":"I_YMFG3","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk3"}}}`,
				buf:    buffer.NewIoBufferString(`<?xml version="1.0" encoding="UTF-8"?><service><body><data name="I_CRNO"><array><struct><data name="O@NAME"><field length="200" scale="0" type="string">SnameClient</field></data><data name="O@ENDDT"><field length="8" scale="0" type="string">SdateDue</field></data></struct><struct><data name="O@IDNO"><field length="40" scale="0" type="string">SidGlobal</field></data></struct><struct><data name="TEST_ARRAY"><array><struct><data name="TEST@MDMLX"><field length="1" scale="0" type="double">6.6</field></data></struct><struct><data name="TEST@PGCP"><field length="2" scale="1" type="double">7.8</field></data></struct></array></data></struct></array></data></body></service>`),
			},
			want:    `{"ccy":[{"dateDue":"SdateDue","nameClient":"SnameClient"},{"idGlobal":"SidGlobal"},{"testArray":[{"testTypeBusi":6.6},{"testCodeProduct":7.8}]}]}`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var relation Relation
			json.Unmarshal([]byte(tt.fields.config), &relation)
			c2b, err := NewCd2Bums(tt.fields.header, tt.fields.buf, &relation)
			assert.NoError(t, err)
			got, err := c2b.BodyInBody()
			assert.Equal(t, err, tt.wantErr)
			assert.Len(t, got.String(), len(tt.want))
		})
	}
}

func TestCd2Bums_Body(t *testing.T) {
	type fields struct {
		config string
		header api.HeaderMap
		buf    api.IoBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "head-array",
			fields: fields{
				config: `{"sys_head":[{"cd_key":"RET_SEQ_NO","type":"string","bums_key":"requestSequenceNo"},{"head_iterms":[{"cd_key":"RET_CODE","type":"string","bums_key":"retCode"},{"cd_key":"RET_MSG","type":"string","bums_key":"retMsg"}],"cd_key":"RET","type":"list","bums_key":"ret"},{"cd_key":"RESP_CODE","type":"string","bums_key":"responseCode"},{"cd_key":"RET_SYS_TIME","type":"string","bums_key":"requestSysTime"},{"cd_key":"SERVICE_SCENE","type":"string","bums_key":"serviceScene"},{"cd_key":"MESSAGE_CODE","type":"string","bums_key":"messageCode"},{"cd_key":"SERVICE_CODE","type":"string","bums_key":"serviceCode"},{"cd_key":"MESSAGE_TYPE","type":"string","bums_key":"messageType"},{"cd_key":"RET_BUSI_DATE","type":"string","bums_key":"requestBusiData"},{"cd_key":"RET_STATUS","type":"string","bums_key":"requestStauts"},{"cd_key":"RET_SYS_DATE","type":"string","bums_key":"requestSysData"}],"app_head":[{"cd_key":"USER_ID","type":"string","bums_key":"userId"},{"cd_key":"AGENT_BRANCH_ID","type":"string","bums_key":"agentBranchId"},{"cd_key":"BRANCH_ID","type":"string","bums_key":"branchId"}],"local_head":[],"body":{"I_CDNO":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"I_CHAF1":{"cd_key":"I_CHAF1","type":"int","length":"3","scale":"0","describes":"渠道标志","bums_key":"flagChannel"},"I_CHAF2":{"cd_key":"I_CHAF2","type":"int","length":"3","scale":"0","describes":"渠道标志2","bums_key":"flagChannel2"},"I_CRNO":{"list_iterms":{"O@ENDDT":{"cd_key":"O@ENDDT","type":"string","length":"8","scale":"0","describes":"有效截止日期","bums_key":"dateDue"},"O@IDNO":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"O@NAME":{"cd_key":"O@NAME","type":"string","length":"200","scale":"0","describes":"姓名","bums_key":"nameClient"},"TEST_ARRAY":{"list_iterms":{"TEST@MDMLX":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"},"TEST@PGCP":{"cd_key":"TEST@PGCP","type":"double","length":"2","scale":"1","describes":"评估产品","bums_key":"testCodeProduct"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"},"I_PSWD":{"cd_key":"I_PSWD","type":"string","length":"16","scale":"0","describes":"密码","bums_key":"password"},"I_YMFG1":{"cd_key":"I_YMFG1","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk"},"I_YMFG2":{"cd_key":"I_YMFG2","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk2"},"I_YMFG3":{"cd_key":"I_YMFG3","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk3"}}}`,
				buf:    buffer.NewIoBufferString(`<?xml version="1.0" encoding="UTF-8"?><service><sys-header><data name="SYS_HEAD"><struct><data name="RET"><array><struct><data name="RET_CODE"><field length="6" scale="0" type="string">999999</field></data><data name="RET_MSG"><field length="9" scale="0" type="string">JDBC调用失败!</field></data></struct></array></data><data name="DEST_BRANCH_NO"><field length="0" scale="0" type="string"/></data><data name="SEQ_NO"><field length="19" scale="0" type="string">ANIU001000000218782</field></data><data name="MESSAGE_CODE"><field length="4" scale="0" type="string">9527</field></data><data name="SERVICE_CODE"><field length="14" scale="0" type="string">BRNC1200952700</field></data><data name="MESSAGE_TYPE"><field length="4" scale="0" type="string">1210</field></data><data name="RET_STATUS"><field length="1" scale="0" type="string">F</field></data><data name="TRAN_TIMESTAMP"><field length="6" scale="0" type="string">135519</field></data><data name="SOURCE_BRANCH_NO"><field length="10" scale="0" type="string">EsbBJFront</field></data><data name="FILE_PATH"><field length="0" scale="0" type="string"/></data><data name="TRAN_DATE"><field length="8" scale="0" type="string">20220113</field></data><data name="BRANCH_ID"><field length="0" scale="0" type="string"/></data></struct></data></sys-header><app-header><data name="APP_HEAD"><struct><data name="AGENT_BRANCH_ID"><field length="9" scale="0" type="string">00301</field></data><data name="USER_ID"><field length="30" scale="0" type="string">BOBQZ2</field></data><data name="BRANCH_ID"><field length="9" scale="0" type="string">00301</field></data></struct></data></app-header><local-header><data name="LOCAL_HEAD"><struct/></data></local-header><body><data name="I_CRNO"><array><struct><data name="O@NAME"><field length="200" scale="0" type="string">SnameClient</field></data><data name="O@ENDDT"><field length="8" scale="0" type="string">SdateDue</field></data></struct><struct><data name="O@IDNO"><field length="40" scale="0" type="string">SidGlobal</field></data></struct><struct><data name="TEST_ARRAY"><array><struct><data name="TEST@MDMLX"><field length="1" scale="0" type="double">6.6</field></data></struct><struct><data name="TEST@PGCP"><field length="2" scale="1" type="double">7.8</field></data></struct></array></data></struct></array></data></body></service>`),
			},
			want:    `{"body":{"ccy":[{"dateDue":"SdateDue","nameClient":"SnameClient"},{"idGlobal":"SidGlobal"},{"testArray":[{"testTypeBusi":6.6},{"testCodeProduct":7.8}]}]},"head":{"agentBranchId":"00301","branchId":"00301","destBranchNo":"","filePath":"","messageCode":"9527","messageType":"1210","requestStauts":"F","retCode":"999999","retMsg":"JDBC调用失败!","seqNo":"ANIU001000000218782","serviceCode":"BRNC1200952700","sourceBranchNo":"EsbBJFront","tranDate":"20220113","tranTimestamp":"135519","userId":"BOBQZ2"}}`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var relation Relation
			json.Unmarshal([]byte(tt.fields.config), &relation)
			c2b, err := NewCd2Bums(tt.fields.header, tt.fields.buf, &relation)
			assert.NoError(t, err)
			got, err := c2b.Body()
			assert.Equal(t, err, tt.wantErr)
			assert.Len(t, got.String(), len(tt.want))
			t.Log(got.String())
		})
	}
}
