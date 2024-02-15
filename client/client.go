package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Команда: go run client.go <порт> <кол-во запросов>")
		return
	}

	port := args[1]
	requests, _ := strconv.Atoi(args[2])

	var wg sync.WaitGroup
	for range requests {
		wg.Add(1)
		go func() {
			SendRequest(port)
			wg.Done()
		}()
	}

	wg.Wait()
}
