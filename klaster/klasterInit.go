package main

import (
	"lib"
	"strconv"
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
)

func klasterInit(startArgs []string) (*lib.Klaster, chan *lib.Worker, chan lib.DeferKlasterToWorkerTask)  {
	if len(startArgs) < 4 {
		panic("Команда: go run klaster.go <id> <порт> <порт менеджера>")
	}

	klasterId, _ := strconv.Atoi(startArgs[1])
	klasterPort := startArgs[2]
	managerPort := startArgs[3]

	klasterStatus = &lib.Klaster{
		Port: klasterPort,
		ManagerPort: managerPort,
		Id: klasterId,
		AllCores: 0,
		FreeCores: 0,
		WorkersList: map[int]*lib.Worker{},
	}

	klasterInfoJson, _ := json.Marshal(klasterStatus)
	http.Post(fmt.Sprintf("http://localhost:%s/addklaster", klasterStatus.ManagerPort), "application/json", bytes.NewReader(klasterInfoJson))
	
	workersPool := make(chan *lib.Worker, 10)
	taskQ := make(chan lib.DeferKlasterToWorkerTask, 100)

	return klasterStatus, workersPool, taskQ
}