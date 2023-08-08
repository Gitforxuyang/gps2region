package main

import (
	"fmt"
	"github.com/Gitforxuyang/gps2region"
	"github.com/Gitforxuyang/gps2region/data"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var (
	lngBegin = 63.66
	lngEnd   = 135.05
	latBegin = 3.86
	latEnd   = 53.55
)

func main() {
	m, err := gps2region.InitGps2Region()
	if err != nil {
		panic(err)
	}
	m.Gps2GeoPosition(84.91479127007814, 44.64531241335134)

	country, err := gps2region.NewGeoMapFromBytes(data.China, "name")
	rand.Seed(time.Now().UnixNano())
	total := 0
	errCount := 0
	for i := 0; i < 10000; i++ {
		randLng := lngBegin + 0.001 + (lngEnd-lngBegin)*rand.Float64()
		randLat := latBegin + 0.001 + (latEnd-latBegin)*rand.Float64()
		//time.Sleep(time.Millisecond * 100)
		if !country.ContainLoc(randLng, randLat) {
			continue
		}
		fmt.Println(randLng, randLat)
		resp, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=&location=%f,%f&poitype=&radius=&extensions=all&roadlevel=0", randLng, randLat))
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json, err := fastjson.Parse(string(body))
		if err != nil {
			panic(err)
		}
		if string(json.GetStringBytes("status")) != "1" {
			fmt.Println(json)
			panic("111")
		}
		if !country.ContainLoc(randLng, randLat) {
			province := string(json.Get("regeocode").Get("addressComponent").GetStringBytes("province"))
			if province == "中华人民共和国" {
				continue
			}
			if string(json.Get("regeocode").GetStringBytes("formatted_address")) != "" {
				fmt.Println(json, randLng, randLat, "不在国内的地理位置判断不一致", province)
			}
			continue
		}
		p, err := m.Gps2GeoPosition(randLng, randLat)
		if err != nil {
			fmt.Println(randLng, randLat, err)
			continue
		}
		addr := json.Get("regeocode").Get("addressComponent")
		errFlag := false
		if string(addr.GetStringBytes("city")) != p.City && string(addr.GetStringBytes("district")) != p.District {
			if string(addr.GetStringBytes("district")) != p.City {
				fmt.Println(randLng, randLat, addr, p, "城市/区县结果不一致")
				errFlag = true
			}
		}
		if string(addr.GetStringBytes("province")) != p.Province {
			fmt.Println(randLng, randLat, addr, p, "省份结果不一致")
			errFlag = true
		}
		//if string(addr.GetStringBytes("adcode")) != p.AdCode {
		//	fmt.Println(addr, p, "adcode结果不一致")
		//	panic("4")
		//}
		if string(addr.GetStringBytes("citycode")) != p.CityCode {
			fmt.Println(addr, p, "citycode结果不一致")
			errFlag = true
		}
		if errFlag {
			errCount++
		}
		total++
		fmt.Println(p, randLng, randLat)
	}
	fmt.Println(total, errCount)
}
