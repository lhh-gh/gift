package util

import (
	"math/rand"
)

// 抽奖。给定每个奖品被抽中的概率（无需要做归一化，但概率必须大于0），返回被抽中的奖品下标
func Lottery(probs []float64) int {
	if len(probs) == 0 {
		return -1
	}
	sum := 0.0
	acc := make([]float64, 0, len(probs)) //累积概率
	for _, prob := range probs {
		sum += prob
		acc = append(acc, sum)
	}

	// 获取(0,sum] 随机数
	r := rand.Float64() * sum
	index := BinarySearch(acc, r)
	return index
}

// 二分法查找数组中>=target的最小的元素下标。arr是单调递增的(里面不能存在重复的元素)，如果target比arr的最后一个元素还大，则返回最后arr的长度；如果target比arr的第一个元素还小，则返回0
func BinarySearch(arr []float64, target float64) int {
	if len(arr) == 0 {
		return -1
	}
	begin, end := 0, len(arr)-1

	for {
		//不论len(arr)为多少，都适用如下2个if
		if target <= arr[begin] {
			return begin
		}
		if target > arr[end] {
			return end + 1
		}

		//len(arr)=2时，适用如下这个if
		if begin == end-1 {
			return end
		}

		//len(arr)>=3时，适用如下流程
		middle := (begin + end) / 2
		if arr[middle] > target {
			end = middle
		} else if arr[middle] < target {
			begin = middle
		} else {
			return middle
		}
	}
}
