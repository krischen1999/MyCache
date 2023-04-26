package main

import "fmt"

func main() {
	a := []byte{'1', '2'}
	b := string(a)
	b = "222"

	fmt.Println(string(a))
	fmt.Println(b)
}
