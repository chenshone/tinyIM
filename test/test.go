package main

import "fmt"

type Person struct {
	Age int
}

func foo(bar *Person) {
	bar.Age = 1
}

func main() {
	a := &Person{}
	foo(a)
	fmt.Printf("%v", a)
}
