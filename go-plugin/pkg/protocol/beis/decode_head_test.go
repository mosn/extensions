package beis

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
)

func TestBeisHeaderMarshalIndent(t *testing.T) {

	v := &SystemHeader{}
	(*v)["MessageType"] = "1410"
	(*v)["MessageCode"] = "3028"
	(*v)["RetStatus"] = "1"

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)

}

func TestBeisHeaderUnmarshal(t *testing.T) {

	v := SystemHeader{}

	data := `
		<SysHead>
			<Ret>
				<RetMsg>[DFDK_BRNC140030280001#ESBOUT1_0570_E007]:[BOBITSAdapter]IO异常！</RetMsg>
				<RetCode>E007</RetCode>
			</Ret>
			<MessageType>1410</MessageType>
			<MessageCode>3028</MessageCode>
			<RetStatus>1</RetStatus>
			<ServiceScene>01</ServiceScene>
			<ServiceCode>BRNC1400302800</ServiceCode>
		</SysHead>
	`
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	delete(v, "Ret")

	if len(v) > 0 {
		for k, d := range v {
			t.Logf("xml key=%s, val=%s", k, d)
		}
	}
}
