package main

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func TestBeisResp2BumsRespTranscoder(t *testing.T) {
	type fields struct {
		root *etree.Element
	}
	tests := []struct {
		name string
		val  string
		want string
	}{
		{
			name: "normal",
			val:  `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead><AppHead><Traceid>CBT02100000000000011</Traceid><Spanid>CBT01100000000000011</Spanid><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead><A1>a1</A1><Cc><C3><D1>d1</D1><D2>D2</D2></C3><Ee><E1>E1</E1><E2>e2</E2></Ee><Ee><E3>E3</E3><E4>e4</E4></Ee><C1>c1</C1><C2>C2</C2></Cc><A2>A2</A2><details><dd><D1>d1</D1><d2>D2</d2></dd><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee></details></Document>`,
			want: `{"head":{"traceId":"CBT02100000000000011","branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","spanId":"CBT01100000000000011","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{"A1":"a1","Cc":{"C3":{"d1":"d1","d2":"D2"},"Ee":{"e3":"E3","e4":"e4"},"c1":"c1","c2":"C2"},"A2":"A2","details":{"dd":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}]}}}`,
		},
		{
			name: "body 走head 解析",
			val:  `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead><AppHead><Traceid>CBT02100000000000011</Traceid><Spanid>CBT01100000000000011</Spanid><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead><A1>a1</A1><Cc><C3><D1>d1</D1><D2>D2</D2></C3><EE><Er2><E3>E3</E3><E4>e4</E4></Er2><Er1><E1>E1</E1><E2>e2</E2></Er1></EE><C1>c1</C1><C2>C2</C2></Cc><A2>A2</A2><details><dd><D1>d1</D1><d2>D2</d2></dd><ee><Er2><e3>E3</e3><E4>e4</E4></Er2><er1><e1>E1</e1><E2>e2</E2></er1></ee></details></Document>`,
			want: `{"head":{"traceId":"CBT02100000000000011","branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","spanId":"CBT01100000000000011","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{"A1":"a1","Cc":{"C3":{"d1":"d1","d2":"D2"},"EE":{"e1":"E1","e2":"e2","e3":"E3","e4":"e4"},"c1":"c1","c2":"C2"},"A2":"A2","details":{"dd":[{"D1":"d1","d2":"D2"}],"ee":[{"e3":"E3","E4":"e4"},{"e1":"E1","E2":"e2"},{"Er2":{"e3":"E3","E4":"e4"},"er1":{"e1":"E1","E2":"e2"}}]}}}`,
		}, {
			name: "details 中包含子对象",
			val:  `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead><AppHead><Traceid>CBT02100000000000011</Traceid><Spanid>CBT01100000000000011</Spanid><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead><A1>a1</A1><A2>A2</A2><details><dd><D1>d1</D1><d2>D2</d2></dd><ee><Er2><e3>E3</e3><E4>e4</E4></Er2><er1><e1>E1</e1><E2>e2</E2></er1></ee><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee></details></Document>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br2br := NewBeisResp2BumsResp(tt.val)
			got := br2br.Transcoder()
			assert.Equal(t, got, tt.want)
		})
	}
}
