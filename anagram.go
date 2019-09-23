package main

import (
	"fmt"
	"os"
	"reflect"
)

func getCharMap(str string) map[string]int {
	found := make(map[string]int)

	for i := 0; i < len(str); i++ {
		if _, ok := found[string(str[i])]; ok {
			found[string(str[i])] = found[string(str[i])] + 1
		} else {
			found[string(str[i])] = 1
		}
	}

	return found
}

func isAnagram(str1, str2 string) bool {
	if str1 == str2 {
		return true
	}
	if len(str1) != len(str2) {
		return false
	}
	found1 := getCharMap(str1)
	found2 := getCharMap(str2)
	return reflect.DeepEqual(found1, found2)
}

func main() {
	str1 := os.Args[1]
	str2 := os.Args[2]

	if isAnagram(str1, str2) {
		fmt.Println("anagram")
	} else {
		fmt.Println("not anagram")
	}
}
