package main

import (
	"encoding/json"
	rq "lib/requests"
	
	
	"fmt"
	is "lib/infostructs"
	"net/http"
	"os"
	"worker/winit"
	"worker/workpool"
)

var workerInfo *is.WorkerInfo

func wsolveproblem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запрос на решение от кластера")
	rq.SendRequest(workerInfo.ClusterPort, "cbusyworker", workerInfo.Id)

	var newReq rq.ClusterWorkerReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		panic("Ошибка парса на воркере, wsolveproblem")
	}
	resultjson, err := json.Marshal(workpool.DoExpr(newReq))
	if err != nil {
		panic("Ошибка парса на ответ в воркере wsolveproblem")
	}


	rq.SendRequest(workerInfo.ClusterPort, "cfreeworker", workerInfo.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultjson)
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
