package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/protocol/beis"
)

func TestBums2beisGetConfig(t *testing.T) {
	type fields struct {
		cfg  map[string]interface{}
		beis api.HeaderMap
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "GetConfig",
			fields: fields{
				cfg:  map[string]interface{}{"details": `[{"uniqueId":"","path":"/","method":"GET","gw":"","resp_mapping":{"sys_head":["ServiceCode","ServiceScene","MessageType","MessageCode","ConsumerSvrId","ConsumerSeqNo","ConsumerId","TranTimestamp","TranDate","TranCode"],"app_head":["UniqueId","AdminUserIdA","Traceid","Spanid","BranchId","AgentBranchId","UserId","VerifyUserId"],"detail_switch":false,"body_switch":false}}]`},
				beis: &beis.Request{},
			},
			want:    `{"uniqueId":"","path":"/","method":"GET","gw":"","resp_mapping":{"sys_head":["ServiceCode","ServiceScene","MessageType","MessageCode","ConsumerSvrId","ConsumerSeqNo","ConsumerId","TranTimestamp","TranDate","TranCode"],"app_head":["UniqueId","AdminUserIdA","Traceid","Spanid","BranchId","AgentBranchId","UserId","VerifyUserId"],"detail_switch":false,"body_switch":false}}`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bmbi := &beis2bums{
				cfg:  tt.fields.cfg,
				beis: tt.fields.beis,
			}
			got, err := bmbi.GetConfig()
			assert.Equal(t, tt.wantErr, err)

			str, _ := json.Marshal(got)
			assert.Equal(t, string(str), tt.want)
		})
	}
}
