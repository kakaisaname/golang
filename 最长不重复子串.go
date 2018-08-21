package main

import "fmt"

func lengthOfNonRepeatingSubStr(s string) int {
	lastOccurred := make(map[rune]int) //字典
	//初始化开始和最大长度
	start := 0
	maxLength := 0
	for i,ch := range []rune(s) {
		if lastI,ok := lastOccurred[ch];ok && lastI >= start {
			start = lastI + 1
		}

		if i -start+1 > maxLength { //i -start+1 子串从start开始到i结束 中间的长度 > maxLength 更新最大长度
			maxLength = i-start+1
		}
		//记录每个字母最后出现的位置 更新lastOccurred[ch]
		lastOccurred[ch] = i
	}
	return maxLength
}
func main() {
	fmt.Println(
		lengthOfNonRepeatingSubStr("abcabcbb"))
	fmt.Println(
		lengthOfNonRepeatingSubStr("abcdef"))
	fmt.Println(
		lengthOfNonRepeatingSubStr("我是中国人名的人"))
}
