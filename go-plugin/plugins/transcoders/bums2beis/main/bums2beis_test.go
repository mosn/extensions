package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/pkg/protocol/http"
)

func TestBums2beisGetConfig(t *testing.T) {
	type fields struct {
		cfg  map[string]interface{}
		bums api.HeaderMap
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   string
		wantErr error
	}{
		{
			name: "GetConfig",
			fields: fields{
				cfg:  map[string]interface{}{"details": `[{"uniqueId":"","service_code":"01","service_scene":"GENL.1400.0382.00","gw":"xxx@gw","req_mapping":{"sys_head":["ServiceCode","ServiceScene","MessageType","MessageCode","ConsumerSvrId","ConsumerSeqNo","ConsumerId","TranTimestamp","TranDate","TranCode"],"app_head":["UniqueId","AdminUserIdA","Traceid","Spanid","BranchId","AgentBranchId","UserId","VerifyUserId"],"detail_switch":false,"body_switch":false}}]`},
				bums: http.RequestHeader{},
			},
			want:    `{"sys_head":["ServiceCode","ServiceScene","MessageType","MessageCode","ConsumerSvrId","ConsumerSeqNo","ConsumerId","TranTimestamp","TranDate","TranCode"],"app_head":["UniqueId","AdminUserIdA","Traceid","Spanid","BranchId","AgentBranchId","UserId","VerifyUserId"],"detail_switch":false,"body_switch":false}`,
			want1:   `{"namespace":"genl.1400.0382.00.01","gw":"xxx@gw"}`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bmbi := &bums2beis{
				cfg:  tt.fields.cfg,
				bums: tt.fields.bums,
			}
			got, got1, err := bmbi.GetConfig(context.TODO())
			assert.Equal(t, tt.wantErr, err)

			str, _ := json.Marshal(got)
			assert.Equal(t, string(str), tt.want)

			str, _ = json.Marshal(got1)
			assert.Equal(t, string(str), tt.want1)
		})
	}
}
