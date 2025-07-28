package main

import "fmt"

func main() {
	var a int = 10
	var b float32 = 32
	const pi float64 = 3.14
	fmt.Println("Hello, World! This is a Go program.{pi:", pi, "}")

	x, y := 10, 20

	fmt.Println("a:", a, "b:", b, "pi:", pi)
	fmt.Println("x:", x, "y:", y)

	isbool := true
	fmt.Println("isbool:", isbool)
	changeValue(&x)
	fmt.Println("x after changeValue:", x)
	fmt.Println(&a, &b, &x, &y, &isbool)

}

func changeValue(x *int) {
	*x = 20
}
