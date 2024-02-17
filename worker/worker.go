package main

import (
	"encoding/json"
	rq "lib/requests"

	"fmt"
	gn "lib/generatelib"
	is "lib/infostructs"
	"net/http"
	"os"
	"worker/winit"

	"time"
)

var workerInfo *is.WorkerInfo

func wsolveproblem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запрос на решение от кластера")
	rq.SendRequest(workerInfo.ClusterPort, "cbusyworker", workerInfo.Id)

	time.Sleep(time.Second)
	ans := gn.GenerateRandMatrix(3, 3, 25)
	ansjs, err := json.Marshal(ans)
	if err != nil {
		panic("ошибка парса матрицы в json на воркере")
	}

	rq.SendRequest(workerInfo.ClusterPort, "cfreeworker", workerInfo.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ansjs)
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
