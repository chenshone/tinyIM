package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func foo() (a int, b int) {
	a = 1
	b = 2
	logrus.Fatalf("1111")
	fmt.Println("2222")
	return
}

func main() {
	a, b := foo()
	fmt.Println("1", a, b)
}
