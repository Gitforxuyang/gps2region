# gps2region

因为一直用的是高德的接口来实现GPS经纬度反解析地理位置，只为了这样一个功能开通高德的企业套餐不划算，所以需要找一个可以实现经纬度解析地理位置的库。

但是发现好像没有现成的，不知道是这个需求太简单还是啥其它原因，总之最后我找到了一个基本的库，可以通过地图信息反解析出位置信息，再通过阿里云的地图平台，完善了国内的数据。

目前只到省市一级（直辖市到区县）。


依赖的基础库地址： https://github.com/linvon/golang-geohelper

依赖的地图原始数据： http://datav.aliyun.com/portal/school/atlas/area_selector

依赖的citycode原始数据：https://www.showapi.com/book/view/3761/5


通过与高德API进行随机对比测试，发现在一些省市交界的地方有可能会有定位错误的情况，综合测试下来的错误率大概在不到千分之一左右。


###基本用法：

go get github.com/Gitforxuyang/gps2region

```
m, err := gps2region.InitGps2Region()
if err != nil {
    panic(err)
}
p, err := m.Gps2GeoPosition(79.98649941751329, 40.76153810752187)
if err != nil {
    panic(err)
}
fmt.Println(p)
```