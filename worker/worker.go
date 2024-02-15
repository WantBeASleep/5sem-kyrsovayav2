package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lib"
	"net/http"
	"os"
	"strconv"
)

var workerInfo lib.Worker
var klasterPort string

func main() {
	args := os.Args
	if len(args) < 5 {
		fmt.Println("Используй gor worker.go <id> <кол-во ядер> <порт> <порт кластера>")
		return
	}

	id, _ := strconv.Atoi(args[1])
	cores, _ := strconv.Atoi(args[2])
	workerPort := args[3]
	klasterPort = args[4]

	workerInfo = lib.Worker{
		Port: workerPort,
		Id: id,
		Cores: cores,
	}

	workerInfoJson, _ := json.Marshal(workerInfo)
	http.Post(fmt.Sprintf("http://localhost:%s/addworker", klasterPort), "application/json", bytes.NewReader(workerInfoJson))

	
	http.ListenAndServe(fmt.Sprintf(":%s", workerInfo.Port), nil)
}
