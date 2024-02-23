package ext

import (
	"fmt"
	"testing"
)

func Test_haversineDistance(t *testing.T) {
	//value, ok := isTable[Device]()
	//if ok == false || value == nil {
	//	t.Errorf("IsTable faild")
	//}

	lat1, lon1 := 31.2304, 121.4737 // 上海
	lat2, lon2 := 39.9042, 116.4074 // 北京

	distance := GetDistance(lat1, lon1, lat2, lon2)
	fmt.Printf("上海和北京之间的距离为：%.2f米	", distance)
	fmt.Printf("上海和北京之间的距离为：%v米	", distance)
}
