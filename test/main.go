package main

import (
	"fmt"
	"github.com/Gitforxuyang/gps2region"
)

func main() {
	m, err := gps2region.InitGps2Region()
	if err != nil {
		panic(err)
	}
	p, err := m.Gps2GeoPosition(121.509062,25.044332)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
}
