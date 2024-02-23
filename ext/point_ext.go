package ext

import (
	"math"
)

// Point 经纬度
type Point struct {
	Lng  float64 //经度
	Lat  float64 //纬度
	Sort int     //序号
}

const (
	EARTH_RADIUS = 6371393 // 地球半径，单位：米
)

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// GetDistance 计算两个经纬度距离，单位米
func GetDistance(lon1, lat1, lon2, lat2 float64) float64 {
	if lon1 == 0 || lat1 == 0 || lon2 == 0 || lat2 == 0 {
		return 0
	}

	lat1 = toRadians(lat1)
	lon1 = toRadians(lon1)
	lat2 = toRadians(lat2)
	lon2 = toRadians(lon2)

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

// GetMultiDistance 计算多个经纬度 与 给定坐标的距离，单位米
func GetMultiDistance(lon1 float64, lat1 float64, points []*Point) map[int]float64 {
	var result = make(map[int]float64, 0)
	if len(points) == 0 {
		return result
	}

	for _, point := range points {
		result[point.Sort] = GetDistance(lon1, lat1, point.Lng, point.Lat)
	}

	return result
}
