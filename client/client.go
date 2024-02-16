package main

import (
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		panic("<порт> <кол-во запросов>")
	}

	port := args[1]
	requests, err := strconv.Atoi(args[2])
	if err != nil {
		panic("ошибка конвертации кол-во запросов")
	}

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
