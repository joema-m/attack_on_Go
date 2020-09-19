package main

import (
	"fmt"
)

// 暴力解法
// 看来用Go的人还是少啊，这样的算法都能超过95的人

func longestCommonPrefixOne(strs []string) string{
	length := len(strs)
	if length == 0 {
		return ""
	}
	for i:=0;i<len(strs[0]);i++ {
		c := strs[0][i]
		for j:=1; j<len(strs);j++ {
			if i == len(strs[j]) || strs[j][i] != c {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}


func longestCommonPrefixTwo(strs []string) string{
	length := len(strs)
	if length == 0 {
		return ""
	}
	minstring := min(strs)
	maxstring := max(strs)
	//min := findMinString(strs)
	//max := findMaxString(strs)
	lesslength := lessLength(minstring, maxstring)
	//fmt.Println(min, max)
	for i:=0; i<lesslength; i++ {
		if minstring[i] != maxstring[i] {
			return maxstring[0:i]
		}
	}
	return minstring
}

func lessLength(strone, strtwo string) int {
	if len(strone) < len(strtwo) {
		return len(strone)
	}
	return len(strtwo)
}

func min(strs []string) string {
	if len(strs) == 0{
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	min := strs[0]
	for _,v := range strs[1:] {
		if stringCompare(min, v) == 1{
			min = v
		}
	}
	return min
}

func max(strs []string) string {
	if len(strs) == 0{
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	max := strs[0]
	for _,v := range strs[1:] {
		if stringCompare(max, v) == -1{
			max = v
		}
	}
	return max
}

// 重复造轮子  熟悉一下
// 可以直接用 strings.Compare()
func stringCompare(str1, str2 string) int {
	length := lessLength(str1,str2)
	if length == 0 {
		if len(str1) == 0 {
			return -1
		} else {
			return 1
		}
	}
	flag := false
	for i:=0; i< length; i++ {
		if str1[i] < str2[i] {
			return -1
		} else{
			if str1[i] == str2[i] {
				flag = true
				continue
			} else {
				return 1
			}
		}
	}
	if flag {
		if len(str1) < len(str2) {
			return -1
		} else {
			if len(str1) > len(str2){
				return 1
			}

		}
	}
	return 0
}


func main() {
	strs := []string{"abab","aba",""}
	fmt.Println(min(strs))
	fmt.Println(longestCommonPrefixTwo(strs))
	fmt.Println(stringCompare("a",""))
}