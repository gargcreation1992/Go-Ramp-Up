package main

import (
	"fmt"
	"goFullStack/wordCount"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: go run <Program_Name> <Directory>")
		return
	}
	t1 := time.Now()
	wordCount.ProcessDir(args[0])
	elapsed1 := time.Since(t1)
	fmt.Println("Elapsed time is: ", elapsed1)
}
