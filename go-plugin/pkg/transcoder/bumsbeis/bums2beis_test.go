package bumsbeis

import (
	"context"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"mosn.io/extensions/go-plugin/pkg/protocol/beis"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/protocol/http"
)

func TestBmr2BirTrancoder(t *testing.T) {
	cfg := &Bums2BeisConfig{
		SysHead:      []string{"ServiceCode", "ServiceScene", "MessageType", "MessageCode", "ConsumerSvrId", "ConsumerSeqNo", "ConsumerId", "TranTimestamp", "TranDate", "TranCode"},
		AppHead:      []string{"UniqueId", "AdminUserIdA", "Traceid", "Spanid", "BranchId", "AgentBranchId", "UserId", "VerifyUserId"},
		DetailSwitch: true,
		BodySwitch:   false,
	}
	vo := &Bums2BeisVo{Namespace: "genl.1400.0382.00.01"}

	header := &fasthttp.RequestHeader{}
	header.Set("TraceId", "CBT02100000000000011")
	header.Set("SpanId", "CBT01100000000000011")
	header.Set("origsender", "1234")
	header.Set("ctrlbits", "4321")
	header.Set("areacode", "0001")
	header.Set("versionid", "0002")

	tests := []struct {
		name   string
		fields string
		want   string
	}{
		{
			name:   "normal",
			fields: `{"head":{"branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","versionId":"1.0.2","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{"A1":"a1","cc":{"c3":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"C1":"c1","c2":"C2"},"dd":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"a2":"A2"}}`,
			want:   `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead><AppHead><Traceid>CBT02100000000000011</Traceid><Spanid>CBT01100000000000011</Spanid><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead><A1>a1</A1><cc><c3><D1>d1</D1><d2>D2</d2></c3><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee><C1>c1</C1><c2>C2</c2></cc><a2>A2</a2><details><Dd><D1>d1</D1><D2>D2</D2></Dd><Ee><E1>E1</E1><E2>e2</E2></Ee><Ee><E3>E3</E3><E4>e4</E4></Ee></details></Document>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buffer.NewIoBufferString(tt.fields)
			br2br, _ := NewBums2Beis(context.Background(), http.RequestHeader{header}, buffer, cfg, vo)
			_, body, err := br2br.Transcoder(true)
			assert.NoError(t, err)
			assert.Len(t, body.String(), len(tt.want))
			t.Log(body.String())
			// xml 校验平台https://extendsclass.com/xml-diff.html#result
			// assert.Equal(t, got, tt.want)
		})
	}
}

func TestBmr2BirDetailFalse(t *testing.T) {
	vo := &Bums2BeisVo{Namespace: "genl.1400.0382.00.01"}
	cfg := &Bums2BeisConfig{
		SysHead:      []string{"ServiceCode", "ServiceScene", "MessageType", "MessageCode", "ConsumerSvrId", "ConsumerSeqNo", "ConsumerId", "TranTimestamp", "TranDate", "TranCode"},
		AppHead:      []string{"UniqueId", "AdminUserIdA", "Traceid", "Spanid", "BranchId", "AgentBranchId", "UserId", "VerifyUserId"},
		DetailSwitch: false,
		BodySwitch:   true,
	}
	header := &fasthttp.RequestHeader{}
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
			name:   "normal",
			fields: `{"head":{"branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","versionId":"1.0.2","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{"A1":"a1","cc":{"c3":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"C1":"c1","c2":"C2"},"dd":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"a2":"A2"}}`,
			want:   `<?xml version="1.0" encoding="UTF-8"?><Document xmlns="genl.1400.0382.00.01"><SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead><AppHead><Traceid>CBT02100000000000011</Traceid><Spanid>CBT01100000000000011</Spanid><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead><A1>a1</A1><cc><c3><D1>d1</D1><d2>D2</d2></c3><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee><C1>c1</C1><c2>C2</c2></cc><a2>A2</a2><details><Dd><D1>d1</D1><D2>D2</D2></Dd><Ee><E1>E1</E1><E2>e2</E2></Ee><Ee><E3>E3</E3><E4>e4</E4></Ee></details></Document>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buffer.NewIoBufferString(tt.fields)
			br2br, _ := NewBums2Beis(context.Background(), http.RequestHeader{header}, buffer, cfg, vo)
			_, body, err := br2br.Transcoder(true)
			assert.NoError(t, err)
			assert.Len(t, body.String(), len(tt.want))
			// xml 校验平台https://extendsclass.com/xml-diff.html#result
			//			assert.Equal(t, got, tt.want)
		})
	}
}

func TestBm2BiSysHead(t *testing.T) {
	vo := &Bums2BeisVo{Namespace: "genl.1400.0382.00.01"}
	cfg := &Bums2BeisConfig{
		SysHead:      []string{"ServiceCode", "ServiceScene", "MessageType", "MessageCode", "ConsumerSvrId", "ConsumerSeqNo", "ConsumerId", "TranTimestamp", "TranDate", "TranCode"},
		AppHead:      []string{},
		DetailSwitch: false,
		BodySwitch:   false,
	}
	header := &fasthttp.RequestHeader{}
	tests := []struct {
		name   string
		fields string
		cfg    Bums2BeisConfig
		want   string
	}{
		{
			name:   "normal",
			fields: `{"head":{"branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","versionId":"1.0.2","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{}}`,
			want:   `<SysHead><ServiceCode>SDKCC140000500</ServiceCode><ServiceScene>01</ServiceScene><MessageType>1400</MessageType><MessageCode>0382</MessageCode><ConsumerSvrId>SmartESB</ConsumerSvrId><ConsumerSeqNo>ACBT001100000000000011</ConsumerSeqNo><ConsumerId>CBT001</ConsumerId><TranTimestamp>144256</TranTimestamp><TranDate>20210824</TranDate><TranCode>BDP071</TranCode></SysHead>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buffer.NewIoBufferString(tt.fields)
			br2br, _ := NewBums2Beis(context.Background(), http.RequestHeader{header}, buffer, cfg, vo)

			doc := etree.NewDocument()
			sysHead := doc.CreateElement("SysHead")
			err := br2br.SysHead(sysHead, &beis.Request{})
			assert.NoError(t, err)
			str, _ := doc.WriteToString()
			assert.Equal(t, str, tt.want)
		})
	}
}

func TestBm2BiAppHead(t *testing.T) {
	vo := &Bums2BeisVo{Namespace: "genl.1400.0382.00.01"}
	cfg := &Bums2BeisConfig{
		SysHead:      []string{},
		AppHead:      []string{"UniqueId", "AdminUserIdA", "Traceid", "Spanid", "BranchId", "AgentBranchId", "UserId", "VerifyUserId"},
		DetailSwitch: false,
		BodySwitch:   false,
	}
	header := &fasthttp.RequestHeader{}
	header.Add("TraceId", "CBT02100000000000011")
	header.Add("SpanId", "CBT01100000000000011")

	tests := []struct {
		name   string
		fields string
		cfg    Bums2BeisConfig
		want   string
	}{
		{
			name:   "normal",
			fields: `{"head":{"branchId":"90153","agentBranchId":"90153","verifyUserId":"B0BQZ6","serviceCode":"SDKCC140000500","consumerId":"CBT001","tranTimestamp":"144256","userId":"CBT","consumerSeqNo":"ACBT001100000000000011","versionId":"1.0.2","serviceScene":"01","messageType":"1400","consumerSvrId":"SmartESB","messageCode":"0382","tranDate":"20210824","tranCode":"BDP071","adminUserIdA":"B0BCBT","uniqueId":"10000000000011"},"body":{}}`,
			want:   `<AppHead><Traceid/><Spanid/><Uniqueid>10000000000011</Uniqueid><AdminUserIdA>B0BCBT</AdminUserIdA><BranchId>90153</BranchId><AgentBranchId>90153</AgentBranchId><UserId>CBT</UserId><VerifyUserId>B0BQZ6</VerifyUserId></AppHead>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buffer.NewIoBufferString(tt.fields)
			br2br, _ := NewBums2Beis(context.Background(), http.RequestHeader{header}, buffer, cfg, vo)
			doc := etree.NewDocument()
			sysHead := doc.CreateElement("AppHead")
			err := br2br.AppHead(sysHead, &beis.Request{})
			assert.NoError(t, err)
			str, _ := doc.WriteToString()
			assert.Equal(t, str, tt.want)
		})
	}
}

func TestBm2BiBodyTrue(t *testing.T) {
	vo := &Bums2BeisVo{
		Namespace: "genl.1400.0382.00.01",
	}
	cfg := &Bums2BeisConfig{
		SysHead:      []string{},
		AppHead:      []string{},
		DetailSwitch: false,
		BodySwitch:   false,
	}
	header := &fasthttp.RequestHeader{}

	tests := []struct {
		name   string
		fields string
		cfg    Bums2BeisConfig
		want   string
	}{
		{
			name:   "normal",
			fields: `{"head":{},"body":{"A1":"a1","cc":{"c3":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"C1":"c1","c2":"C2"},"a2":"A2"}}`,
			want:   `<Document><A1>a1</A1><cc><c3><D1>d1</D1><d2>D2</d2></c3><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee><C1>c1</C1><c2>C2</c2></cc><a2>A2</a2></Document>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buffer.NewIoBufferString(tt.fields)
			br2br, _ := NewBums2Beis(context.Background(), http.RequestHeader{header}, buffer, cfg, vo)

			doc := etree.NewDocument()
			element := doc.CreateElement("Document")
			err := br2br.Body(element)
			assert.NoError(t, err)

			str, _ := doc.WriteToString()
			assert.Equal(t, str, tt.want)
			// xml 校验平台https://extendsclass.com/xml-diff.html#result
		})
	}
}

func TestBm2BiBodyFalse(t *testing.T) {
	vo := &Bums2BeisVo{Namespace: "genl.1400.0382.00.01"}
	cfg := &Bums2BeisConfig{
		SysHead:      []string{},
		AppHead:      []string{},
		DetailSwitch: false,
		BodySwitch:   true,
	}
	header := &fasthttp.RequestHeader{}

	tests := []struct {
		name   string
		fields string
		cfg    Bums2BeisConfig
		want   string
	}{
		{
			name:   "normal",
			fields: `{"head":{},"body":{"A1":"a1","cc":{"c3":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}],"C1":"c1","c2":"C2"},"a2":"A2"}}`,
			want:   `<Document><A1>a1</A1><Cc><C3><D1>d1</D1><D2>D2</D2></C3><Ee><E1>E1</E1><E2>e2</E2></Ee><Ee><E3>E3</E3><E4>e4</E4></Ee><C1>c1</C1><C2>C2</C2></Cc><A2>A2</A2></Document>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buffer.NewIoBufferString(tt.fields)
			br2br, _ := NewBums2Beis(context.Background(), http.RequestHeader{header}, buffer, cfg, vo)

			doc := etree.NewDocument()
			element := doc.CreateElement("Document")
			err := br2br.Body(element)
			assert.NoError(t, err)

			str, _ := doc.WriteToString()
			assert.Equal(t, str, tt.want)
		})
	}
}

func TestBm2BiDetailTrue(t *testing.T) {
	vo := &Bums2BeisVo{Namespace: "genl.1400.0382.00.01"}
	cfg := &Bums2BeisConfig{
		SysHead:      []string{},
		AppHead:      []string{},
		DetailSwitch: true,
		BodySwitch:   true,
	}
	header := &fasthttp.RequestHeader{}

	tests := []struct {
		name   string
		fields string
		cfg    Bums2BeisConfig
		want   string
	}{
		{
			name:   "normal",
			fields: `{"head":{},"body":{"dd":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}]}}`,
			want:   `<Document><details><Dd><D1>d1</D1><D2>D2</D2></Dd><Ee><E1>E1</E1><E2>e2</E2></Ee><Ee><E3>E3</E3><E4>e4</E4></Ee></details></Document>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buffer.NewIoBufferString(tt.fields)
			br2br, _ := NewBums2Beis(context.Background(), http.RequestHeader{header}, buffer, cfg, vo)

			doc := etree.NewDocument()
			element := doc.CreateElement("Document")
			err := br2br.Body(element)
			assert.NoError(t, err)

			str, _ := doc.WriteToString()
			assert.Equal(t, str, tt.want)
			// xml 校验平台https://extendsclass.com/xml-diff.html#result
		})
	}
}

func TestBm2BiDetailFalse(t *testing.T) {
	vo := &Bums2BeisVo{Namespace: "genl.1400.0382.00.01"}
	cfg := &Bums2BeisConfig{
		SysHead:      []string{},
		AppHead:      []string{},
		DetailSwitch: false,
		BodySwitch:   true,
	}
	header := &fasthttp.RequestHeader{}

	tests := []struct {
		name   string
		fields string
		cfg    Bums2BeisConfig
		want   string
	}{
		{
			name:   "normal",
			fields: `{"head":{},"body":{"dd":[{"D1":"d1","d2":"D2"}],"ee":[{"e1":"E1","E2":"e2"},{"e3":"E3","E4":"e4"}]}}`,
			want:   `<Document><details><dd><D1>d1</D1><d2>D2</d2></dd><ee><e1>E1</e1><E2>e2</E2></ee><ee><e3>E3</e3><E4>e4</E4></ee></details></Document>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buffer.NewIoBufferString(tt.fields)
			br2br, _ := NewBums2Beis(context.Background(), http.RequestHeader{header}, buffer, cfg, vo)

			doc := etree.NewDocument()
			element := doc.CreateElement("Document")
			err := br2br.Body(element)
			assert.NoError(t, err)

			str, _ := doc.WriteToString()
			assert.Equal(t, str, tt.want)
			// xml 校验平台https://extendsclass.com/xml-diff.html#result
		})
	}
}
