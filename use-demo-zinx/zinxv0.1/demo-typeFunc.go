package main

import "fmt"

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

// 其中sub可以看作函数名，是OP的类型
type op func(a, b int) int

// Ope 将函数类型作为形式参数传递
func Ope(o op, a, b int) int {
	return o(a, b)
}

func main() {
	// 在go语言中函数名可以看作函数类型的常量，所以我们可以直接把函数名作为参数传入函数
	res := Ope(sub, 3, 1)
	fmt.Println(res)
}
