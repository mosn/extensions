package bumscd

import (
	"encoding/json"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/protocol/bolt"
)

func TestBumsReq2CdReq_GetXmlBytes(t *testing.T) {
	header := &bolt.RequestHeader{}
	header.Set("OrigSender", "QDT001")
	type fields struct {
		header api.HeaderMap
		config string
		value  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "normal",
			fields: fields{
				config: `{"sys_head":[{"cd_key":"TRAN_TIMESTAMP","type":"string","length":"9","scale":"0","source":"head","bums_key":"tranTimestamp"},{"cd_key":"MODULE_ID","type":"string","length":"4","scale":"0","default":"RB"},{"cd_key":"SERVICE_SCENE","type":"string","length":"2","scale":"0","default":"03"},{"cd_key":"SOURCE_BRANCH_NO","type":"string","length":"10","scale":"0","default":"EsbBJFront"},{"cd_key":"TRAN_DATE","type":"string","length":"8","scale":"0","source":"head","bums_key":"tranDate"},{"cd_key":"CONSUMER_ID","type":"string","length":"6","scale":"0","source":"headers","bums_key":"OrigSender","default":"QDT001"},{"cd_key":"CONSUMER_SEQ_NO","type":"string","length":"52","scale":"0","source":"head","bums_key":"consumerSeqNo"},{"cd_key":"CONSUMER_SVR_ID","type":"string","length":"8","scale":"0","default":"SmartUNFront"},{"cd_key":"TRAN_CODE","type":"string","length":"4","scale":"0","source":"head","bums_key":"tranCode"},{"cd_key":"SERVICE_CODE","type":"string","length":"30","scale":"0","default":"ECIF1200003000"},{"cd_key":"MESSAGE_CODE","type":"string","length":"6","scale":"0","default":"0008"}],"app_head":[{"cd_key":"USER_ID","type":"string","length":"30","scale":"0","default":"BOBQZ2"},{"cd_key":"AGENT_BRANCH_ID","type":"string","length":"9","scale":"0","source":"head","bums_key":"branchId"},{"cd_key":"BRANCH_ID","type":"string","length":"9","scale":"0","source":"head","bums_key":"branchId"}],"local_head":[],"body":{"cardOrAcctNo":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"ccy":{"list_iterms":{"dateDue":{"cd_key":"O@ENDDT","type":"string","length":"8","scale":"0","describes":"有效截止日期","bums_key":"dateDue"},"idGlobal":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"nameClient":{"cd_key":"O@NAME","type":"string","length":"200","scale":"0","describes":"姓名","bums_key":"nameClient"},"testArray":{"list_iterms":{"testCodeProduct":{"cd_key":"TEST@PGCP","type":"double","length":"2","scale":"1","describes":"评估产品","bums_key":"testCodeProduct"},"testTypeBusi":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"},"flagChannel":{"cd_key":"I_CHAF1","type":"int","length":"3","scale":"0","describes":"渠道标志","bums_key":"flagChannel"},"flagChannel2":{"cd_key":"I_CHAF2","type":"int","length":"3","scale":"0","describes":"渠道标志2","bums_key":"flagChannel2"},"flagPwdChk":{"cd_key":"I_YMFG1","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk"},"flagPwdChk2":{"cd_key":"I_YMFG2","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk2"},"flagPwdChk3":{"cd_key":"I_YMFG3","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk3"},"password":{"cd_key":"I_PSWD","type":"string","length":"16","scale":"0","describes":"密码","bums_key":"password"}}}`,
				value:  `{"body":{"cardOrAcctNo":"WcardOrAcctNo","ccy":[{"dateDue":"SdateDue","nameClient":"SnameClient"},{"idGlobal":"SidGlobal"},{"testArray":[{"testTypeBusi":6.6},{"testCodeProduct":7.8}]}],"flagChannel":"1","flagChannel2":"2","flagPwdChk":"1.1","flagPwdChk2":"2.1","flagPwdChk3":"3.1","password":"Wpassword"},"head":{"areaCode":"0000","branchId":"90012","consumerId":"QDT001","consumerSeqNo":"AQDT202108110670585","tranCode":"ACC582300","tranDate":"20210811","tranTimestamp":"123456","versionId":"0001"}}`,
				header: header,
			},
			want:    `<?xml version="1.0" encoding="UTF-8"?><service><sys-header><data name="SYS_HEAD"><struct><data name="CONSUMER_SVR_ID"><field length="8" scale="0" type="string">SmartUNFront</field></data><data name="SERVICE_SCENE"><field length="2" scale="0" type="string">03</field></data><data name="SERVICE_CODE"><field length="30" scale="0" type="string">ECIF1200003000</field></data><data name="MESSAGE_CODE"><field length="6" scale="0" type="string">0008</field></data><data name="CONSUMER_SEQ_NO"><field length="52" scale="0" type="string">AQDT202108110670585</field></data><data name="TRAN_TIMESTAMP"><field length="9" scale="0" type="string">123456</field></data><data name="MODULE_ID"><field length="4" scale="0" type="string">RB</field></data><data name="SOURCE_BRANCH_NO"><field length="10" scale="0" type="string">EsbBJFront</field></data><data name="TRAN_DATE"><field length="8" scale="0" type="string">20210811</field></data><data name="CONSUMER_ID"><field length="6" scale="0" type="string">QDT001</field></data><data name="TRAN_CODE"><field length="4" scale="0" type="string">ACC582300</field></data></struct></data></sys-header><app-header><data name="APP_HEAD"><struct><data name="AGENT_BRANCH_ID"><field length="9" scale="0" type="string">90012</field></data><data name="USER_ID"><field length="30" scale="0" type="string">BOBQZ2</field></data><data name="BRANCH_ID"><field length="9" scale="0" type="string">90012</field></data></struct></data></app-header><local-header><data name="LOCAL_HEAD"><struct/></data></local-header><body><data name="I_CDNO"><field length="19" scale="0" type="string">WcardOrAcctNo</field></data><data name="I_CHAF1"><field length="3" scale="0" type="int">1</field></data><data name="I_PSWD"><field length="16" scale="0" type="string">Wpassword</field></data><data name="I_YMFG3"><field length="5" scale="3" type="double">3.1</field></data><data name="I_CHAF2"><field length="3" scale="0" type="int">2</field></data><data name="I_CRNO"><array><struct><data name="O@NAME"><field length="200" scale="0" type="string">SnameClient</field></data><data name="O@ENDDT"><field length="8" scale="0" type="string">SdateDue</field></data></struct><struct><data name="O@IDNO"><field length="40" scale="0" type="string">SidGlobal</field></data></struct><struct><data name="TEST_ARRAY"><array><struct><data name="TEST@MDMLX"><field length="1" scale="0" type="double">6.6</field></data></struct><struct><data name="TEST@PGCP"><field length="2" scale="1" type="double">7.8</field></data></struct></array></data></struct></array></data><data name="I_YMFG2"><field length="5" scale="3" type="double">2.1</field></data><data name="I_YMFG1"><field length="5" scale="3" type="double">1.1</field></data></body></service>`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var relation Relation
			json.Unmarshal([]byte(tt.fields.config), &relation)
			br2cd, err := NewBums2Cd(tt.fields.header, tt.fields.value, &relation)
			assert.NoError(t, err)
			got, err := br2cd.GetXmlBytes()
			assert.Len(t, B2S(got), len(tt.want))
			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestBumsReq2CdReq_ParseBody(t *testing.T) {
	header := &bolt.RequestHeader{}
	header.Set("OrigSender", "QDT001")
	type fields struct {
		header api.HeaderMap
		config string
		value  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "normal-head",
			fields: fields{
				config: `{"sys_head":[],"app_head":[],"local_head":[],"body":{"cardOrAcctNo":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"ccy":{"list_iterms":{"dateDue":{"cd_key":"O@ENDDT","type":"string","length":"8","scale":"0","describes":"有效截止日期","bums_key":"dateDue"},"idGlobal":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"nameClient":{"cd_key":"O@NAME","type":"string","length":"200","scale":"0","describes":"姓名","bums_key":"nameClient"},"testArray":{"list_iterms":{"testCodeProduct":{"cd_key":"TEST@PGCP","type":"double","length":"2","scale":"1","describes":"评估产品","bums_key":"testCodeProduct"},"testTypeBusi":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"},"flagChannel":{"cd_key":"I_CHAF1","type":"int","length":"3","scale":"0","describes":"渠道标志","bums_key":"flagChannel"},"flagChannel2":{"cd_key":"I_CHAF2","type":"int","length":"3","scale":"0","describes":"渠道标志2","bums_key":"flagChannel2"},"flagPwdChk":{"cd_key":"I_YMFG1","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk"},"flagPwdChk2":{"cd_key":"I_YMFG2","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk2"},"flagPwdChk3":{"cd_key":"I_YMFG3","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk3"},"password":{"cd_key":"I_PSWD","type":"string","length":"16","scale":"0","describes":"密码","bums_key":"password"}}}`,
				value:  `{"body":{"flagChannel":"1","flagPwdChk":"1.1","password":"Wpassword"},"head":{}}`,
				header: header,
			},
			want:    `<?xml version="1.0" encoding="UTF-8"?><service><body><data name="I_CHAF1"><field length="3" scale="0" type="int">1</field></data><data name="I_PSWD"><field length="16" scale="0" type="string">Wpassword</field></data><data name="I_YMFG1"><field length="5" scale="3" type="double">1.1</field></data></body></service>`,
			wantErr: nil,
		},
		{
			name: "normal-list",
			fields: fields{
				config: `{"sys_head":[],"app_head":[],"local_head":[],"body":{"cardOrAcctNo":{"cd_key":"I_CDNO","type":"string","length":"19","scale":"0","describes":"卡号/账号","bums_key":"cardOrAcctNo"},"ccy":{"list_iterms":{"dateDue":{"cd_key":"O@ENDDT","type":"string","length":"8","scale":"0","describes":"有效截止日期","bums_key":"dateDue"},"idGlobal":{"cd_key":"O@IDNO","type":"string","length":"40","scale":"0","describes":"证件号码","bums_key":"idGlobal"},"nameClient":{"cd_key":"O@NAME","type":"string","length":"200","scale":"0","describes":"姓名","bums_key":"nameClient"},"testArray":{"list_iterms":{"testCodeProduct":{"cd_key":"TEST@PGCP","type":"double","length":"2","scale":"1","describes":"评估产品","bums_key":"testCodeProduct"},"testTypeBusi":{"cd_key":"TEST@MDMLX","type":"double","length":"1","scale":"0","describes":"面对面类型","bums_key":"testTypeBusi"}},"cd_key":"TEST_ARRAY","type":"list","bums_key":"testArray"}},"cd_key":"I_CRNO","type":"list","bums_key":"ccy"},"flagChannel":{"cd_key":"I_CHAF1","type":"int","length":"3","scale":"0","describes":"渠道标志","bums_key":"flagChannel"},"flagChannel2":{"cd_key":"I_CHAF2","type":"int","length":"3","scale":"0","describes":"渠道标志2","bums_key":"flagChannel2"},"flagPwdChk":{"cd_key":"I_YMFG1","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk"},"flagPwdChk2":{"cd_key":"I_YMFG2","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk2"},"flagPwdChk3":{"cd_key":"I_YMFG3","type":"double","length":"5","scale":"3","describes":"验密标志","bums_key":"flagPwdChk3"},"password":{"cd_key":"I_PSWD","type":"string","length":"16","scale":"0","describes":"密码","bums_key":"password"}}}`,
				value:  `{"body":{"ccy":[{"dateDue":"SdateDue","nameClient":"SnameClient"},{"idGlobal":"SidGlobal"},{"testArray":[{"testTypeBusi":6.6},{"testCodeProduct":7.8}]}]},"head":{}}`,
				header: header,
			},
			want:    `<?xml version="1.0" encoding="UTF-8"?><service><body><data name="I_CRNO"><array><struct><data name="O@NAME"><field length="200" scale="0" type="string">SnameClient</field></data><data name="O@ENDDT"><field length="8" scale="0" type="string">SdateDue</field></data></struct><struct><data name="O@IDNO"><field length="40" scale="0" type="string">SidGlobal</field></data></struct><struct><data name="TEST_ARRAY"><array><struct><data name="TEST@MDMLX"><field length="1" scale="0" type="double">6.6</field></data></struct><struct><data name="TEST@PGCP"><field length="2" scale="1" type="double">7.8</field></data></struct></array></data></struct></array></data></body></service>`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var relation Relation
			json.Unmarshal([]byte(tt.fields.config), &relation)
			br2cd, err := NewBums2Cd(tt.fields.header, tt.fields.value, &relation)
			assert.NoError(t, err)

			cd := etree.NewDocument()
			cd.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
			element := cd.CreateElement("service")
			err = br2cd.Body(element)
			assert.Equal(t, err, tt.wantErr)

			str, _ := cd.WriteToString()
			assert.Len(t, str, len(tt.want))
		})
	}
}

func TestBumsReq2CdReq_BodyHead(t *testing.T) {
	header := &bolt.RequestHeader{}
	header.Set("OrigSender", "QDT001")
	type fields struct {
		header api.HeaderMap
		config string
		value  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "normal-head",
			fields: fields{
				config: `{"sys_head":[{"cd_key":"TRAN_TIMESTAMP","type":"string","length":"9","scale":"0","source":"head","bums_key":"tranTimestamp"},{"cd_key":"MODULE_ID","type":"string","length":"4","scale":"0","default":"RB"},{"cd_key":"SERVICE_SCENE","type":"string","length":"2","scale":"0","default":"03"},{"cd_key":"SOURCE_BRANCH_NO","type":"string","length":"10","scale":"0","default":"EsbBJFront"},{"cd_key":"TRAN_DATE","type":"string","length":"8","scale":"0","source":"head","bums_key":"tranDate"},{"cd_key":"CONSUMER_ID","type":"string","length":"6","scale":"0","source":"headers","bums_key":"OrigSender","default":"QDT001"},{"cd_key":"CONSUMER_SEQ_NO","type":"string","length":"52","scale":"0","source":"head","bums_key":"consumerSeqNo"},{"cd_key":"CONSUMER_SVR_ID","type":"string","length":"8","scale":"0","default":"SmartUNFront"},{"cd_key":"TRAN_CODE","type":"string","length":"4","scale":"0","source":"head","bums_key":"tranCode"},{"cd_key":"SERVICE_CODE","type":"string","length":"30","scale":"0","default":"ECIF1200003000"},{"cd_key":"MESSAGE_CODE","type":"string","length":"6","scale":"0","default":"0008"}],"app_head":[{"cd_key":"USER_ID","type":"string","length":"30","scale":"0","default":"BOBQZ2"},{"cd_key":"AGENT_BRANCH_ID","type":"string","length":"9","scale":"0","source":"head","bums_key":"branchId"},{"cd_key":"BRANCH_ID","type":"string","length":"9","scale":"0","source":"head","bums_key":"branchId"}],"local_head":[],"body":{}}`,
				header: header,
				value:  `{"body":{},"head":{"areaCode":"0000","branchId":"90012","consumerId":"QDT001","consumerSeqNo":"AQDT202108110670585","tranCode":"ACC582300","tranDate":"20210811","tranTimestamp":"123456","versionId":"0001"}}`,
			},
			want:    `<?xml version="1.0" encoding="UTF-8"?><service><sys-header><data name="SYS_HEAD"><struct><data name="CONSUMER_SVR_ID"><field length="8" scale="0" type="string">SmartUNFront</field></data><data name="SERVICE_SCENE"><field length="2" scale="0" type="string">03</field></data><data name="SERVICE_CODE"><field length="30" scale="0" type="string">ECIF1200003000</field></data><data name="MESSAGE_CODE"><field length="6" scale="0" type="string">0008</field></data><data name="CONSUMER_SEQ_NO"><field length="52" scale="0" type="string">AQDT202108110670585</field></data><data name="TRAN_TIMESTAMP"><field length="9" scale="0" type="string">123456</field></data><data name="MODULE_ID"><field length="4" scale="0" type="string">RB</field></data><data name="SOURCE_BRANCH_NO"><field length="10" scale="0" type="string">EsbBJFront</field></data><data name="TRAN_DATE"><field length="8" scale="0" type="string">20210811</field></data><data name="CONSUMER_ID"><field length="6" scale="0" type="string">QDT001</field></data><data name="TRAN_CODE"><field length="4" scale="0" type="string">ACC582300</field></data></struct></data></sys-header><app-header><data name="APP_HEAD"><struct><data name="AGENT_BRANCH_ID"><field length="9" scale="0" type="string">90012</field></data><data name="USER_ID"><field length="30" scale="0" type="string">BOBQZ2</field></data><data name="BRANCH_ID"><field length="9" scale="0" type="string">90012</field></data></struct></data></app-header><local-header><data name="LOCAL_HEAD"><struct/></data></local-header></service>`,
			wantErr: nil,
		},
		{
			name: "head-list",
			fields: fields{
				config: `{"sys_head":[{"cd_key":"TRAN_DATE","type":"string","length":"8","scale":"0","source":"head","bums_key":"tranDate"},{"head_iterms":[{"cd_key":"SUBH","type":"string","length":"6","scale":"0","source":"head","bums_key":"get#sub"}],"cd_key":"LISTH","type":"list","source":"head","bums_key":"listh"}]}`,
				header: header,
				value:  `{"body":{},"head":{"consumerSeqNo":"AQDT202108110670585","listh":[{"sub":"hello"}],"tranDate":"20210811"}}`,
			},
			want:    `<?xml version="1.0" encoding="UTF-8"?><service><sys-header><data name="SYS_HEAD"><struct><data name="LISTH"><array><struct><data name="SUBH"><field length="6" scale="0" type="string">hello</field></data></struct></array></data><data name="TRAN_DATE"><field length="8" scale="0" type="string">20210811</field></data></struct></data></sys-header><app-header><data name="APP_HEAD"><struct/></data></app-header><local-header><data name="LOCAL_HEAD"><struct/></data></local-header></service>`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var relation Relation
			json.Unmarshal([]byte(tt.fields.config), &relation)
			br2cd, err := NewBums2Cd(tt.fields.header, tt.fields.value, &relation)
			assert.NoError(t, err)

			cd := etree.NewDocument()
			cd.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
			element := cd.CreateElement("service")
			err = br2cd.BodyHead(element)
			assert.Equal(t, err, tt.wantErr)

			str, _ := cd.WriteToString()
			assert.Len(t, str, len(tt.want))
		})
	}
}
