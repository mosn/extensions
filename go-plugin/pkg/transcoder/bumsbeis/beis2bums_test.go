package bumsbeis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/common"
	"mosn.io/extensions/go-plugin/pkg/protocol/beis"
	"mosn.io/pkg/buffer"
)

func TestBeis2BumsHeader(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		header api.HeaderMap
		want   string
	}{
		{
			name:   "head",
			val:    `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead><AppHead><Traceid>CBT02100000000000011</Traceid><Spanid>CBT01100000000000011</Spanid><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead></Document>`,
			want:   `{"head":{"traceId":"CBT02100000000000011","branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","spanId":"CBT01100000000000011","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{}}`,
			header: &beis.Request{Header: common.Header{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br2br, err := NewBeis2Bums(context.Background(), tt.header, buffer.NewIoBufferString(tt.val), &Beis2BumsConfig{})
			assert.NoError(t, err)
			body, err := br2br.BodyJson(tt.header)
			assert.NoError(t, err)
			assert.Len(t, body, len(tt.want))
			// json :https://www.sojson.com/jsondiff.html
			// t.Log(tt.want)
			// t.Log(body)
		})
	}
}

func TestBeis2BumsHead(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		header api.HeaderMap
		want   string
	}{
		{
			name:   "body",
			val:    `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><Cc><EE><Er2><E3>E3</E3><E4>e4</E4></Er2><Er1><E1>E1</E1><E2>e2</E2></Er1></EE><C3><D1>d1</D1><D2>D2</D2></C3><C1>c1</C1><C2>C2</C2></Cc></Document>`,
			want:   `{"head":{},"body":{"Cc":{"C1":"c1","C2":"C2","C3":{"D1":"d1","D2":"D2"},"EE":{"Er1":{"E1":"E1","E2":"e2"},"Er2":{"E3":"E3","E4":"e4"}}}}}`,
			header: &beis.Request{Header: common.Header{}},
		},
		{
			name:   "head_slice",
			val:    `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><Cc><EE><Er2><E3>E3</E3><E4>e4</E4></Er2><Er1><E1>E1</E1><E2>e2</E2></Er1></EE><EE><E1>E1</E1><E2>e2</E2></EE><C3><D1>d1</D1><D2>D2</D2></C3><C1>c1</C1><C2>C2</C2></Cc></Document>`,
			want:   `{"body":{"Cc":{"C1":"c1","C2":"C2","C3":{"D1":"d1","D2":"D2"},"EE":[{"Er1":{"E1":"E1","E2":"e2"},"Er2":{"E3":"E3","E4":"e4"}},{"E1":"E1","E2":"e2"}]}},"head":{}}`,
			header: &beis.Request{Header: common.Header{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br2br, err := NewBeis2Bums(context.Background(), tt.header, buffer.NewIoBufferString(tt.val), &Beis2BumsConfig{})
			assert.NoError(t, err)
			body, err := br2br.BodyJson(tt.header)
			assert.NoError(t, err)
			assert.Len(t, body, len(tt.want))
			// json :https://www.sojson.com/jsondiff.html
			// t.Log(tt.want)
			// t.Log(body)
		})
	}
}

func TestBeis2BumsDetail(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		header api.HeaderMap
		want   string
	}{
		{
			name:   "details_bug",
			val:    `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><details><cc><EE><Er2><e3>E3</e3><E4>e4</E4></Er2><er1><e1>E1</e1><E2>e2</E2></er1></EE><EE><e1>E1</e1><E2>e2</E2></EE><c3><D1>d1</D1><d2>D2</d2></c3><C1>c1</C1><c2>C2</c2></cc></details></Document>`,
			want:   `{"head":{},"body":{"details":{"cc":{"EE":[{"Er2":{"e3":"E3","E4":"e4"},"er1":{"e1":"E1","E2":"e2"}},{"e1":"E1","E2":"e2"}],"c3":{"D1":"d1","d2":"D2"},"C1":"c1","c2":"C2"}}}}`,
			header: &beis.Request{},
		},
		{
			name:   "details_objectbug",
			val:    `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><details><dd><D1>d1</D1><d2>D2</d2></dd><ee><Er2><e3>E3</e3><E4>e4</E4></Er2><er1><e1>E1</e1><E2>e2</E2></er1></ee><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee></details></Document>`,
			want:   `{"head":{},"body":{"details":{"ee":[{"Er2":{"e3":"E3","E4":"e4"},"er1":{"e1":"E1","E2":"e2"}},{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"dd":{"D1":"d1","d2":"D2"}}}}`,
			header: &beis.Request{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br2br, err := NewBeis2Bums(context.Background(), tt.header, buffer.NewIoBufferString(tt.val), &Beis2BumsConfig{})
			assert.NoError(t, err)
			body, err := br2br.BodyJson(tt.header)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(body))
			t.Log(tt.want)
			t.Log(body)
		})
	}
}
