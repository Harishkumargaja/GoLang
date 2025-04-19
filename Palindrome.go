package main

import (
    "fmt"
)

func palindrome(s string) bool {
	for i:=0; i < len(s)/2; i++{
		if s[i]!=s[len(s)-1-i]{
			return false
		}
	}
	
	return true
}

func main() {

    const s = "roar"
	var t = palindrome(s)
	
	
	
	if t{
		fmt.Println("It is a palindrome")
	}else{
		fmt.Println("It is not a palindrome")
	}
}
