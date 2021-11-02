package main

import (
	"fmt"
	"os"
)

const n = 3

func main() {
	var arr [n]int
	for i := 0; i < n; i++ {
		cnt, err := fmt.Scan(&arr[i])
		if cnt != 1 || err != nil {
			os.Exit(1)
		}
	}
}
