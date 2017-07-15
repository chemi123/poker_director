package main

import "fmt"

func main() {
	tmp := make([]int, 0)
	fmt.Println("tmp:", tmp, " len(tmp):", len(tmp), " cap(tmp):", cap(tmp))
	tmp = append(tmp, 1)
	fmt.Println("append 1 to tmp")
	fmt.Println("tmp:", tmp, " len(tmp):", len(tmp), " cap(tmp):", cap(tmp))
	tmp = append(tmp, 1)
	fmt.Println("append 1 to tmp")
	fmt.Println("tmp:", tmp, " len(tmp):", len(tmp), " cap(tmp):", cap(tmp))
	tmp = append(tmp, 1)
	fmt.Println("append 1 to tmp")
	fmt.Println("tmp:", tmp, " len(tmp):", len(tmp), " cap(tmp):", cap(tmp))
	tmp = append(tmp, 1)
	fmt.Println("append 1 to tmp")
	fmt.Println("tmp:", tmp, " len(tmp):", len(tmp), " cap(tmp):", cap(tmp))
	tmp = append(tmp, 1)
	fmt.Println("append 1 to tmp")
	fmt.Println("tmp:", tmp, " len(tmp):", len(tmp), " cap(tmp):", cap(tmp))
}
