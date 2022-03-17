package bumsbeis

import (
	"fmt"
	"strings"

	"mosn.io/api"
	"mosn.io/extensions/go-plugin/pkg/protocol/beis"
)

func (br2br *Bums2Beis) HeadRequest() (api.HeaderMap, error) {
	// 大小写转换
	beisHeader := &beis.Request{}
	br2br.header.Range(func(k, v string) bool {
		k = strings.ToLower(k)
		beisHeader.Set(k, v)
		return true
	})

	origsender, ok := beisHeader.Get("origsender")
	if !ok || len(origsender) > 10 {
		return nil, fmt.Errorf("the %s of origsender is illage", origsender)
	}
	beisHeader.OrigSender = origsender
	beisHeader.Del("origsender")

	ctrlbits, ok := beisHeader.Get("ctrlbits")
	if !ok || len(ctrlbits) > 8 {
		return nil, fmt.Errorf("the %s of ctrlbits is illage", ctrlbits)
	}
	beisHeader.CtrlBits = ctrlbits
	beisHeader.Del("ctrlbits")

	areacode, ok := beisHeader.Get("areacode")
	if !ok || len(areacode) != 4 {
		return nil, fmt.Errorf("the %s of areacode is illage", areacode)
	}
	beisHeader.AreaCode = areacode
	beisHeader.Del("areacode")

	versionid, ok := beisHeader.Get("versionid")
	if !ok || len(versionid) != 4 {
		return nil, fmt.Errorf("the %s of versionid is illage", versionid)
	}
	beisHeader.VersionID = versionid
	beisHeader.Del("versionid")

	// 校验traceid ，spanid
	traceid, ok := beisHeader.Get("traceid")
	if !ok {
		return nil, fmt.Errorf("the %s of traceid is illage", traceid)
	}

	spanid, ok := beisHeader.Get("spanid")
	if !ok {
		return nil, fmt.Errorf("the %s of spanid is illage", spanid)
	}

	if len(br2br.vo.MesgId) != 0 {
		beisHeader.MessageID = br2br.vo.MesgId
	}

	if len(br2br.vo.MesgRefId) != 0 {
		beisHeader.MessageRefID = br2br.vo.MesgRefId
	}

	if len(br2br.vo.Reserve) != 0 {
		beisHeader.Reserve = br2br.vo.Reserve
	}
	beisHeader.Set("service", br2br.vo.GWName)
	return beisHeader, nil
}

func (br2br *Bums2Beis) HeadRespone() (api.HeaderMap, error) {
	// 大小写转换
	beisHeader := &beis.Response{}
	br2br.header.Range(func(k, v string) bool {
		k = strings.ToLower(k)
		beisHeader.Set(k, v)
		return true
	})

	origsender, ok := beisHeader.Get("origsender")
	if !ok || len(origsender) > 10 {
		return nil, fmt.Errorf("the %s of origsender is illage", origsender)
	}
	beisHeader.OrigSender = origsender
	beisHeader.Del("origsender")

	ctrlbits, ok := beisHeader.Get("ctrlbits")
	if !ok || len(ctrlbits) > 8 {
		return nil, fmt.Errorf("the %s of ctrlbits is illage", ctrlbits)
	}
	beisHeader.CtrlBits = ctrlbits
	beisHeader.Del("ctrlbits")

	areacode, ok := beisHeader.Get("areacode")
	if !ok || len(areacode) != 4 {
		return nil, fmt.Errorf("the %s of areacode is illage", areacode)
	}
	beisHeader.AreaCode = areacode
	beisHeader.Del("areacode")

	versionid, ok := beisHeader.Get("versionid")
	if !ok || len(versionid) != 4 {
		return nil, fmt.Errorf("the %s of versionid is illage", versionid)
	}
	beisHeader.VersionID = versionid
	beisHeader.Del("versionid")

	// 校验traceid ，spanid
	traceid, ok := beisHeader.Get("traceid")
	if !ok {
		return nil, fmt.Errorf("the %s of traceid is illage", traceid)
	}

	spanid, ok := beisHeader.Get("spanid")
	if !ok {
		return nil, fmt.Errorf("the %s of spanid is illage", spanid)
	}

	if len(br2br.vo.MesgId) != 0 {
		beisHeader.MessageID = br2br.vo.MesgId
	}

	if len(br2br.vo.MesgRefId) != 0 {
		beisHeader.MessageRefID = br2br.vo.MesgRefId
	}

	if len(br2br.vo.Reserve) != 0 {
		beisHeader.Reserve = br2br.vo.Reserve
	}

	beisHeader.Set("service", br2br.vo.GWName)
	return beisHeader, nil
}
