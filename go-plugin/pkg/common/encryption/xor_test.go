package encryption

import (
	"reflect"
	"testing"
)

func TestXorEncoder(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
		msg  []byte
		want []int8
	}{
		{
			name: "XorEncoder",
			key:  []byte{7, 2, 34, 9, 10, 6, 78},
			msg:  []byte{1, 5, 8, 3, 5, 24, 89, 0, 35, 10, 39, 127, 46, 97, 11, 16, 18, 49, 29, 47, 14, 29},
			want: []int8{6, 7, 42, 10, 15, 30, 23, 7, 33, 40, 46, 117, 40, 47, 12, 18, 48, 56, 23, 41, 64, 26},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := XorEncrypt(tt.msg, tt.key)
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("XorEncoder() got = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestXorDecoder(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
		msg  []int8
		want []byte
	}{
		{
			name: "XorDecoder",
			key:  []byte{7, 2, 34, 9, 10, 6, 78},
			msg:  []int8{6, 7, 42, 10, 15, 30, 23, 7, 33, 40, 46, 117, 40, 47, 12, 18, 48, 56, 23, 41, 64, 26},
			want: []byte{1, 5, 8, 3, 5, 24, 89, 0, 35, 10, 39, 127, 46, 97, 11, 16, 18, 49, 29, 47, 14, 29},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := XorDecrypt(tt.msg, tt.key)
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("XorDecoder() got = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestDecoder(t *testing.T) {
	tests := []struct {
		msg  string
		want string
	}{
		{
			msg:  "DQ1LWVkWQV1DQVpbWwsVCR8CERRQWFRXVVtdUwgUYmx3HwsWCgg9BHVdUEFYU1lMEUpeWFtFChpWV11YGwcFCAEcAwQFBxkIARwDBRcIPQRiS0B8UFdTBjsOYVFBdEJLWHZSQFAIBQgDAAMFBAcLF2NXR3ZARV58UEZWCj8KZV1FYUpHcVdDXQ8AAwYHBgYJAA4cZlBCZEFCdlJAUAg9BGNXR2dMRWNRXFcNBQAGBgwHAQ8bZ1NDa0hBZ11YUwkyDWBWQGZTRnZeDHJ6fGMHCAACAwQFBgcJBwcFDQUKGGpURmBRRHhYBjsOYVFBCD0EY1dHd1pSUgZzBwMFCRllXUVxXFBQCD0EY1dHeUZRCQ8BAlBbWBhETV9LUkZRGFJAUldDQFxZWRZiR11xdntyQFJXQ0BcWVkCEdW5gtO2tt6Rtdubs9mLonQCAwcZ0Li32Y2D24ms0bGI1J+V0YmW3rCd14y/0ouY1JaC3IGTG920hdeJptOxvdScituJrN+XlNS6jdOalt2GgNeMv9KLmNScv9K9ptiEvdqcg9C9iNG2v9aQuNChiNSQrQgaZFJMfEFUCj8KGGpURg0+CRlkQUJ6VlVRCD0EcEJDfFBXUxcPOA8bcVlUTVxXXUALPA==",
			want: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Document xmlns=\"genl.1200.0001.00.01\">\n<SysHead>\n<RetBusiDate>20220111</RetBusiDate>\n<RetSysDate>20220111</RetSysDate>\n<RetSysTime>1501463</RetSysTime>\n<RetSeqNo>ANIU001000000165690</RetSeqNo>\n<Ret>\n<RetCode>B501</RetCode>\n<RetMsg>700com.sunyard.exception.SunECMException: 状态标识：E003,描述：批次信息上传失败,具体内容：该批次已上传完成，请勿重复新增</RetMsg>\n</Ret>\n</SysHead>\n<AppHead/>\n</Document>\n",
		},
		{
			msg:  "DQ1LWVkWQV1DQVpbWwsVCR8CERRQWFRXVVtdUwgUYmx3HwsWCgg9BHVdUEFYU1lMEUpeWFtFChpTQF1XGwcDCAEcAAQHDhkIARwDBRcIPQRwcXBganh4BgEDAw0FBQMJBQIDBQcGBggEAwoBAQUFBB5zcHdhaXl3DzgPd3pycmd+YHQbCzwLa0hBe1FUUgkyDXFcWkZDWl1De1cKe3NjCAEDDxt2WVlLRF9WRnxSCTINf1xQQFpScVUMS01WChh1XlZGWFB/UwY7DmNGWlFFWVx7VwphdXRgDR1jRlpRRVlce1cKPwp0V19BRllQRGRdQHxcCnQGBwgBAgMEBQYHCAELAQwBBAQEHnFcWkZDWl1DYVZFe1kJMg1mQVVbclZMVAwBBAcEBwkBBg8bYURWVnVTR1ELPAtsQ1NdYFxbUktFU15ECwYHCQkBCggaYkVZX2ZaWVBFQ1lcQg0+CWJFWV9xXFBQCE92c2YPG2FEVlZyXVdRCzwLF2JLQHxQV1MGOw5yREV+UllVDDkId0RWVlJaelALDwcIAQEPG3dEVlZSWnpQCzwLeVZXXUB3RFZWUlp6UAsPBwgBAQ8bdFFSVkVwQVVbVV9xVQw5CGBFUkp4Vg0FBwcLF2RBVkZ8UgkyDWRWRlxQTm1CV0F9UQgGCgAOHGJQRF5eSGdAUUd/UwY7Dhx1RUZ/XVBWDT4JGXNXUkdeUVtCCTI=",
			want: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Document xmlns=\"brnc.1400.3028.00.01\">\n<ACCT_NO>01090341400120105195432</ACCT_NO>\n<CODE_ORG/>\n<SysHead>\n<ConsumerId>NET001</ConsumerId>\n<ModuleId>xyc</ModuleId>\n<ProgramId>TCCX</ProgramId>\n<ConsumerSeqNo>A000000000000928423</ConsumerSeqNo>\n<TranDate>20220104</TranDate>\n<TranTimestamp>001839</TranTimestamp>\n<TranCode>xNBT</TranCode>\n</SysHead>\n<AppHead>\n<BranchId>90003</BranchId>\n<AgentBranchId>90003</AgentBranchId>\n<UserId>121</UserId>\n<VerifyUserId>121</VerifyUserId>\n</AppHead>\n</Document>\n",
		},
	}

	for _, tt := range tests {
		got := string(XorDecrypt(Base64Decoder([]byte(tt.msg)), []byte("12345678")))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Base64Encoder() got = %s, want %s", got, tt.want)
		}
	}
}

func TestEncoder(t *testing.T) {
	tests := []struct {
		want string
		msg  string
	}{
		{
			want: "DQ1LWVkWQV1DQVpbWwsVCR8CERRQWFRXVVtdUwgUYmx3HwsWCgg9BHVdUEFYU1lMEUpeWFtFChpWV11YGwcFCAEcAwQFBxkIARwDBRcIPQRiS0B8UFdTBjsOYVFBdEJLWHZSQFAIBQgDAAMFBAcLF2NXR3ZARV58UEZWCj8KZV1FYUpHcVdDXQ8AAwYHBgYJAA4cZlBCZEFCdlJAUAg9BGNXR2dMRWNRXFcNBQAGBgwHAQ8bZ1NDa0hBZ11YUwkyDWBWQGZTRnZeDHJ6fGMHCAACAwQFBgcJBwcFDQUKGGpURmBRRHhYBjsOYVFBCD0EY1dHd1pSUgZzBwMFCRllXUVxXFBQCD0EY1dHeUZRCQ8BAlBbWBhETV9LUkZRGFJAUldDQFxZWRZiR11xdntyQFJXQ0BcWVkCEdW5gtO2tt6Rtdubs9mLonQCAwcZ0Li32Y2D24ms0bGI1J+V0YmW3rCd14y/0ouY1JaC3IGTG920hdeJptOxvdScituJrN+XlNS6jdOalt2GgNeMv9KLmNScv9K9ptiEvdqcg9C9iNG2v9aQuNChiNSQrQgaZFJMfEFUCj8KGGpURg0+CRlkQUJ6VlVRCD0EcEJDfFBXUxcPOA8bcVlUTVxXXUALPA==",
			msg:  "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Document xmlns=\"genl.1200.0001.00.01\">\n<SysHead>\n<RetBusiDate>20220111</RetBusiDate>\n<RetSysDate>20220111</RetSysDate>\n<RetSysTime>1501463</RetSysTime>\n<RetSeqNo>ANIU001000000165690</RetSeqNo>\n<Ret>\n<RetCode>B501</RetCode>\n<RetMsg>700com.sunyard.exception.SunECMException: 状态标识：E003,描述：批次信息上传失败,具体内容：该批次已上传完成，请勿重复新增</RetMsg>\n</Ret>\n</SysHead>\n<AppHead/>\n</Document>\n",
		},
		{
			want: "DQ1LWVkWQV1DQVpbWwsVCR8CERRQWFRXVVtdUwgUYmx3HwsWCgg9BHVdUEFYU1lMEUpeWFtFChpTQF1XGwcDCAEcAAQHDhkIARwDBRcIPQRwcXBganh4BgEDAw0FBQMJBQIDBQcGBggEAwoBAQUFBB5zcHdhaXl3DzgPd3pycmd+YHQbCzwLa0hBe1FUUgkyDXFcWkZDWl1De1cKe3NjCAEDDxt2WVlLRF9WRnxSCTINf1xQQFpScVUMS01WChh1XlZGWFB/UwY7DmNGWlFFWVx7VwphdXRgDR1jRlpRRVlce1cKPwp0V19BRllQRGRdQHxcCnQGBwgBAgMEBQYHCAELAQwBBAQEHnFcWkZDWl1DYVZFe1kJMg1mQVVbclZMVAwBBAcEBwkBBg8bYURWVnVTR1ELPAtsQ1NdYFxbUktFU15ECwYHCQkBCggaYkVZX2ZaWVBFQ1lcQg0+CWJFWV9xXFBQCE92c2YPG2FEVlZyXVdRCzwLF2JLQHxQV1MGOw5yREV+UllVDDkId0RWVlJaelALDwcIAQEPG3dEVlZSWnpQCzwLeVZXXUB3RFZWUlp6UAsPBwgBAQ8bdFFSVkVwQVVbVV9xVQw5CGBFUkp4Vg0FBwcLF2RBVkZ8UgkyDWRWRlxQTm1CV0F9UQgGCgAOHGJQRF5eSGdAUUd/UwY7Dhx1RUZ/XVBWDT4JGXNXUkdeUVtCCTI=",
			msg:  "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Document xmlns=\"brnc.1400.3028.00.01\">\n<ACCT_NO>01090341400120105195432</ACCT_NO>\n<CODE_ORG/>\n<SysHead>\n<ConsumerId>NET001</ConsumerId>\n<ModuleId>xyc</ModuleId>\n<ProgramId>TCCX</ProgramId>\n<ConsumerSeqNo>A000000000000928423</ConsumerSeqNo>\n<TranDate>20220104</TranDate>\n<TranTimestamp>001839</TranTimestamp>\n<TranCode>xNBT</TranCode>\n</SysHead>\n<AppHead>\n<BranchId>90003</BranchId>\n<AgentBranchId>90003</AgentBranchId>\n<UserId>121</UserId>\n<VerifyUserId>121</VerifyUserId>\n</AppHead>\n</Document>\n",
		},
	}

	for _, tt := range tests {

		got := string(Base64Encoder(XorEncrypt([]byte(tt.msg), []byte("12345678"))))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Base64Encoder() got = %s, want %v", got, tt.want)
		}
	}
}

func TestName(t *testing.T) {
	tests := []struct {
		want string
		msg  string
	}{
		{
			want: "DQ1LWVkWQV1DQVpbWwsVCR8CERRQWFRXVVtdUwgUYmx3HwsWCgg9BHVdUEFYU1lMEUpeWFtFChpWV11YGwcFCAEcAwQFBxkIARwDBRcIPQRiS0B8UFdTBjsOYVFBdEJLWHZSQFAIBQgDAAMFBAcLF2NXR3ZARV58UEZWCj8KZV1FYUpHcVdDXQ8AAwYHBgYJAA4cZlBCZEFCdlJAUAg9BGNXR2dMRWNRXFcNBQAGBgwHAQ8bZ1NDa0hBZ11YUwkyDWBWQGZTRnZeDHJ6fGMHCAACAwQFBgcJBwcFDQUKGGpURmBRRHhYBjsOYVFBCD0EY1dHd1pSUgZzBwMFCRllXUVxXFBQCD0EY1dHeUZRCQ8BAlBbWBhETV9LUkZRGFJAUldDQFxZWRZiR11xdntyQFJXQ0BcWVkCEdW5gtO2tt6Rtdubs9mLonQCAwcZ0Li32Y2D24ms0bGI1J+V0YmW3rCd14y/0ouY1JaC3IGTG920hdeJptOxvdScituJrN+XlNS6jdOalt2GgNeMv9KLmNScv9K9ptiEvdqcg9C9iNG2v9aQuNChiNSQrQgaZFJMfEFUCj8KGGpURg0+CRlkQUJ6VlVRCD0EcEJDfFBXUxcPOA8bcVlUTVxXXUALPA==",
			msg:  "{\"body\":{\"acctNo\":\"01090341400120105195432\",\"orgCode\":\"\"},\"head\":{\"agentBranchId\":\"90003\",\"tranDate\":\"20220104\",\"verifyUserId\":\"121\",\"branchId\":\"90003\",\"tranCode\":\"NBT302800\",\"tranTimestamp\":\"001839\",\"userId\":\"121\",\"consumerSeqNo\":\"A000000000000928423\",\"consumerId\":\"NET001\"}}",
		},
		{
			want: "DQ1LWVkWQV1DQVpbWwsVCR8CERRQWFRXVVtdUwgUYmx3HwsWCgg9BHVdUEFYU1lMEUpeWFtFChpTQF1XGwcDCAEcAAQHDhkIARwDBRcIPQRwcXBganh4BgEDAw0FBQMJBQIDBQcGBggEAwoBAQUFBB5zcHdhaXl3DzgPd3pycmd+YHQbCzwLa0hBe1FUUgkyDXFcWkZDWl1De1cKe3NjCAEDDxt2WVlLRF9WRnxSCTINf1xQQFpScVUMS01WChh1XlZGWFB/UwY7DmNGWlFFWVx7VwphdXRgDR1jRlpRRVlce1cKPwp0V19BRllQRGRdQHxcCnQGBwgBAgMEBQYHCAELAQwBBAQEHnFcWkZDWl1DYVZFe1kJMg1mQVVbclZMVAwBBAcEBwkBBg8bYURWVnVTR1ELPAtsQ1NdYFxbUktFU15ECwYHCQkBCggaYkVZX2ZaWVBFQ1lcQg0+CWJFWV9xXFBQCE92c2YPG2FEVlZyXVdRCzwLF2JLQHxQV1MGOw5yREV+UllVDDkId0RWVlJaelALDwcIAQEPG3dEVlZSWnpQCzwLeVZXXUB3RFZWUlp6UAsPBwgBAQ8bdFFSVkVwQVVbVV9xVQw5CGBFUkp4Vg0FBwcLF2RBVkZ8UgkyDWRWRlxQTm1CV0F9UQgGCgAOHGJQRF5eSGdAUUd/UwY7Dhx1RUZ/XVBWDT4JGXNXUkdeUVtCCTI=",
			msg:  "{\"head\":{\"ctrlBits\":\"10000000\",\"retCode\":\"0000\",\"retMsg\":\"xxx\"},\"body\":{\"dataArray\":[{\"ENTRUSTRESULT\":\"2\",\"UUID\":\"2022010405884506\",\"FAILREASON\":\"xxxx\"},{\"ENTRUSTRESULT\":\"2\",\"UUID\":\"2022010405885199\",\"FAILREASON\":\"xxx\"},{\"ENTRUSTRESULT\":\"2\",\"UUID\":\"2022010405885201\",\"FAILREASON\":\"xxx\"}]}}",
		},
	}

	for _, tt := range tests {

		got := string(Base64Encoder(XorEncrypt([]byte(tt.msg), []byte("12345678"))))
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Base64Encoder() got = %s, want %v", got, tt.want)
		}
	}
}
