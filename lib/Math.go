package lib

import "math/rand"

// 生成区间-随机数
func RandInt(min, max int) int {
	if min >= max {
		return max
	}
	return rand.Intn(max-min) + min
}
