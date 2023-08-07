package gps2region

import (
	"encoding/json"
	"errors"
	"github.com/Gitforxuyang/gps2region/data"
	"strings"
)

type Position struct {
	Province string //省或直辖市
	City     string //市 当属于直辖市时，此字段返回为空
	District string //区县
	CityCode string //城市编码
	AdCode   string //行政区编码
}
type Gps2Region interface {
	Gps2GeoPosition(lng, lat float64) (Position, error)
}

type gps2Region struct {
	country            *GeoMap
	province           map[string]*GeoMap
	adCode2CityMapping map[string]*data.City
}
type province struct {
	data []byte
}

var (
	provinceMap = map[string]province{
		"上海市":      {data: data.ShangHai},
		"北京市":      {data: data.BeiJing},
		"重庆市":      {data: data.ChongQing},
		"天津市":      {data: data.TianJin},
		"江苏省":      {data: data.JiangSu},
		"江西省":      {data: data.JiangXi},
		"河南省":      {data: data.HeNan},
		"河北省":      {data: data.HeiBei},
		"山东省":      {data: data.ShanDong},
		"山西省":      {data: data.ShanXi},
		"陕西省":      {data: data.ShaanXi},
		"浙江省":      {data: data.ZheJiang},
		"安徽省":      {data: data.AnHui},
		"福建省":      {data: data.FuJian},
		"湖北省":      {data: data.HuBei},
		"湖南省":      {data: data.HuNan},
		"广东省":      {data: data.GuangDong},
		"广西壮族自治区":  {data: data.GuangXi},
		"海南省":      {data: data.HaiNan},
		"四川省":      {data: data.SiChuan},
		"贵州省":      {data: data.GuiZhou},
		"云南省":      {data: data.YunNan},
		"西藏自治区":    {data: data.XiZang},
		"甘肃省":      {data: data.GanSu},
		"青海省":      {data: data.QingHai},
		"宁夏回族自治区":  {data: data.NingXia},
		"新疆维吾尔自治区": {data: data.XinJiang},
		"台湾省":      {data: data.TaiWan},
		"香港":       {data: data.XiangGang},
		"澳门":       {data: data.AoMen},
		"内蒙古自治区":   {data: data.NeiMengGu},
		"辽宁省":      {data: data.LiaoNing},
		"吉林省":      {data: data.JiLin},
		"黑龙江省":     {data: data.HeiLongJiang},
	}
)

func InitGps2Region() (*gps2Region, error) {
	m := gps2Region{
		province:           map[string]*GeoMap{},
		adCode2CityMapping: map[string]*data.City{},
	}
	ff := func(ks []string) string {
		return ks[0] + "-" + ks[1] + "-" + ks[2]
	}
	country, err := NewGeoMapFromBytes(data.China, "name")
	if err != nil {
		return nil, err
	}
	m.country = country
	for k, v := range provinceMap {
		province, err := NewGeoMapFormatFromBytes(v.data, []string{"name", "adcode", "level"}, ff)
		if err != nil {
			return nil, err
		}
		m.province[k] = province
	}
	allCityCode := []*data.City{}
	err = json.Unmarshal(data.CityData, &allCityCode)
	if err != nil {
		return nil, err
	}
	for k, v := range allCityCode {
		m.adCode2CityMapping[v.AdCode] = allCityCode[k]
	}
	return &m, nil
}

var (
	UnknownPosition = errors.New("未知位置")
)

func (m *gps2Region) Gps2GeoPosition(lng, lat float64) (Position, error) {
	province := m.country.FindLoc(lng, lat)
	if province == StrNotFound {
		return Position{}, UnknownPosition
	}
	position := Position{}
	position.Province = province
	cityResult := m.province[position.Province].FindLoc(lng, lat)
	if cityResult == StrNotFound {
		return position, UnknownPosition
	}
	cityResultArr := strings.Split(cityResult, "-")
	if cityResultArr[2] == "city" {
		position.City = cityResultArr[0]
		position.AdCode = cityResultArr[1]
	} else if cityResultArr[2] == "district" {
		position.District = cityResultArr[0]
		position.AdCode = cityResultArr[1]
	} else if cityResultArr[2] == "province" {
		position.AdCode = cityResultArr[1]
	} else {
		return position, errors.New("未知的level")
	}
	position.CityCode = m.adCode2CityMapping[position.AdCode].CityCode
	return position, nil
}
