package bumsbeis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/extensions/go-plugin/pkg/protocol/beis"
)

var (
	// jsonValue = `{"head":{"branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","versionId":"1.0.2","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{"Slice":["AAA","AAA3"],"BuBU":"BBB","details":{"A":[{"AAA":"AAA","BBB":"BBB"},{"AAA":"AAA2","BBB":"BBB2"}]},"hehe":"BBB"}}`
	jsonValue = `{"head":{"branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","versionId":"1.0.2","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{"hehe":"BBB"}}`
)

func TestBmr2BirDetailTrue(t *testing.T) {
	cfg := Bums2BeisConfig{
		SysHead:      []string{"ServiceCode", "ServiceScene", "MessageType", "MessageCode", "ConsumerSvrId", "ConsumerSeqNo", "ConsumerId", "TranTimestamp", "TranDate", "TranCode"},
		AppHead:      []string{"UniqueId", "AdminUserIdA", "Traceid", "Spanid", "BranchId", "AgentBranchId", "UserId", "VerifyUserId"},
		DetailSwitch: true,
		BodySwitch:   false,
		Namespace:    "genl.1400.0382.00.01",
	}

	header := &beis.Request{}
	header.Set("TraceId", "CBT02100000000000011")
	header.Set("SpanId", "CBT01100000000000011")
	header.Set("origsender", "1234")
	header.Set("ctrlbits", "4321")
	header.Set("areacode", "0001")
	header.Set("versionid", "0002")

	tests := []struct {
		name   string
		fields string
		cfg    Bums2BeisConfig
		want   string
	}{
		{
			name:   "1",
			cfg:    cfg,
			fields: `{"head":{"branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","versionId":"1.0.2","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{"A1":"a1","cc":{"c3":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"C1":"c1","c2":"C2"},"dd":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"a2":"A2"}}`,
			want:   `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead><AppHead><Traceid>CBT02100000000000011</Traceid><Spanid>CBT01100000000000011</Spanid><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead><A1>a1</A1><cc><c3><D1>d1</D1><d2>D2</d2></c3><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee><C1>c1</C1><c2>C2</c2></cc><a2>A2</a2><details><Dd><D1>d1</D1><D2>D2</D2></Dd><Ee><E1>E1</E1><E2>e2</E2></Ee><Ee><E3>E3</E3><E4>e4</E4></Ee></details></Document>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br2br, _ := NewBums2Beis(header, tt.fields, tt.cfg)
			_, body, err := br2br.Transcoder()
			assert.NoError(t, err)
			assert.Len(t, body.String(), len(tt.want))
			t.Log(body.String())
			// xml 校验平台https://extendsclass.com/xml-diff.html#result
			// assert.Equal(t, got, tt.want)
		})
	}
}

func TestBmr2BirDetailFalse(t *testing.T) {
	cfg := Bums2BeisConfig{
		SysHead:      []string{"ServiceCode", "ServiceScene", "MessageType", "MessageCode", "ConsumerSvrId", "ConsumerSeqNo", "ConsumerId", "TranTimestamp", "TranDate", "TranCode"},
		AppHead:      []string{"UniqueId", "AdminUserIdA", "Traceid", "Spanid", "BranchId", "AgentBranchId", "UserId", "VerifyUserId"},
		DetailSwitch: false,
		BodySwitch:   true,
		Namespace:    "genl.1400.0382.00.01",
	}
	header := &beis.Request{}
	header.Set("TraceId", "CBT02100000000000011")
	header.Set("SpanId", "CBT01100000000000011")
	header.Set("origsender", "1234")
	header.Set("ctrlbits", "4321")
	header.Set("areacode", "0001")
	header.Set("versionid", "0002")

	tests := []struct {
		name   string
		fields string
		cfg    Bums2BeisConfig
		want   string
	}{
		{
			name:   "1",
			cfg:    cfg,
			fields: `{"head":{"branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","versionId":"1.0.2","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{"A1":"a1","cc":{"c3":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"C1":"c1","c2":"C2"},"dd":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"a2":"A2"}}`,
			want:   `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead><AppHead><Traceid>CBT02100000000000011</Traceid><Spanid>CBT01100000000000011</Spanid><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead><A1>a1</A1><cc><c3><D1>d1</D1><d2>D2</d2></c3><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee><C1>c1</C1><c2>C2</c2></cc><a2>A2</a2><details><Dd><D1>d1</D1><D2>D2</D2></Dd><Ee><E1>E1</E1><E2>e2</E2></Ee><Ee><E3>E3</E3><E4>e4</E4></Ee></details></Document>`,
		}, /*
			{
				name: " no details slice",
			},
			{
				name: "details slice and  detail object",
			},
			{
				name: " no details slice and no detail object",
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br2br, _ := NewBums2Beis(header, tt.fields, tt.cfg)
			_, body, err := br2br.Transcoder()
			assert.NoError(t, err)
			assert.Len(t, body.String(), len(tt.want))
			// xml 校验平台https://extendsclass.com/xml-diff.html#result
			//			assert.Equal(t, got, tt.want)
		})
	}
}
