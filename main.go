package main

import "fmt"

func func1(a int, b int) int {
	return a + b;
}

func main() {
	result := func1(1, 3);
	fmt.Println(result)
}