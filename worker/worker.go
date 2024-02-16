package main

import (
	"encoding/json"
	gn "lib/generatelib"
	rq "lib/requests"
	"time"
	
	
	"fmt"
	is "lib/infostructs"
	"net/http"
	"os"
	"worker/winit"
)

var workerInfo *is.WorkerInfo

func wsolveproblem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запрос на решение от кластера")
	rq.SendRequest(workerInfo.ClusterPort, "cbusyworker", workerInfo.Id)

	mtrx := gn.GenerateRandMatrix(3, 3, 25)
	mtrxjson, _ := json.Marshal(mtrx)

	time.Sleep(2 * time.Second)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(mtrxjson)

	rq.SendRequest(workerInfo.ClusterPort, "cfreeworker", workerInfo.Id)
}

func main() {
	args := os.Args
	if len(args) < 5 {
		panic("<порт><порт кластера><id><кол-во ядер>")
	}
	workerInfo = winit.WorkerInit(args)

	http.HandleFunc("/wsolveproblem", wsolveproblem)
	http.ListenAndServe(fmt.Sprintf(":%s", workerInfo.Port), nil)
}
