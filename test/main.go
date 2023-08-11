package main

import (
	"encoding/json"
	"fmt"
	"github.com/Gitforxuyang/gps2region"
)

func main() {
	m, err := gps2region.InitGps2Region()
	if err != nil {
		panic(err)
	}
	p, err := m.Gps2GeoPosition(79.98649941751329, 40.76153810752187)
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(p)
	fmt.Println(string(b))
}
