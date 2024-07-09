package main

import "fmt"

func foo() (int, int) {
	return 1, 2
}

func main() {
	s := "aaa"
	sprintf := fmt.Sprintf("%s", s)
	fmt.Println(sprintf)
}
