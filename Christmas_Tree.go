package main

import (
    "fmt"
)

func star(i int) {
	for j:=0; j<=i; j++{
			fmt.Print("*")
		}
		fmt.Println()
}

func main() {
	
    for i:=0; i<=10; i++{
		for j:=0; j<=i; j++{
			fmt.Print("*")
		}
		fmt.Println()
	}
	fmt.Println()
	
	for i:=0; i<=9; i++{
		star(i)
	}
}

    