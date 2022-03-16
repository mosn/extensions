package cd

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
)

func TestCdHeaderMarshalIndent(t *testing.T) {

	v := &SystemHeader{Name: "SYS_HEAD",
		WrapData: []WrapData{
			{
				Name:  "SERVICE_CODE",
				Field: &Field{Length: 14, Scale: 0, Type: "string", Value: "BRNC1200952700"},
			}, {
				Name:  "SERVICE_SCENE",
				Field: &Field{Length: 2, Scale: 0, Type: "string", Value: "01"},
			}, {
				Name: "RET",
				ArrayField: &[]WrapData{
					{
						Name:  "RET_CODE",
						Field: &Field{Length: 6, Scale: 0, Type: "string", Value: "999999"},
					},
					{
						Name:  "RET_MSG",
						Field: &Field{Length: 9, Scale: 0, Type: "string", Value: "JDBC调用失败!"},
					},
				},
			},
		},
	}

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)

}

func TestCdHeaderUnmarshal(t *testing.T) {

	v := SystemHeader{}

	data := `
		  <data name="SYS_HEAD">
			  <struct>
				  <data name="SERVICE_CODE">
					  <field length="14" scale="0" type="string">BRNC1200952700</field>
				  </data>
				  <data name="SERVICE_SCENE">
					  <field length="2" scale="0" type="string">01</field>
				  </data>
				  <data name="RET">
					  <array>
						  <struct>
							  <data name="RET_CODE">
								  <field length="6" scale="0" type="string">999999</field>
							  </data>
							  <data name="RET_MSG">
								  <field length="9" scale="0" type="string">JDBC调用失败!</field>
							  </data>
						  </struct>
					  </array>
				  </data>
			  </struct>
		  </data>
	`
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Printf("XMLName: %#v\n", v.XMLName)
	fmt.Printf("XMLAttr: %q\n", v.Name)
	if len(v.WrapData) > 0 {
		for _, d := range v.WrapData {
			if d.Field != nil { // plain field
				fmt.Printf("field: %s %q\n", d.Name, d.Field.Value)
			} else if d.ArrayField != nil {
				for _, f := range *d.ArrayField {
					if f.Field != nil {
						// struct field
						fmt.Printf("struct field: %s %q\n", f.Name, f.Field.Value)
					}
				}
			}
		}
	}
}
