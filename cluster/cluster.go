package main

import (
	"os"
	"net/http"
	"fmt"
	"io"
	"encoding/json"

	is "lib/infostructs"
	ts "lib/tasks"
	gn "lib/generatelib"

	
	"cluster/cinit"
	"cluster/workers"
)

var clusterInfo *is.ClusterInfo
var workersPool chan *is.WorkerInfo
var deferClusterWorkerTaskPool chan ts.ClusterWorkerTask


// func freeWorkerHandle(w http.ResponseWriter, r *http.Request) {
// 	var fWorker lib.FreeWorkerSignal
// 	json.NewDecoder(r.Body).Decode(&fWorker)
// 	freeWorker(klasterStatus, fWorker, workersPool)
// }

// func busyWorkerHandle(w http.ResponseWriter, r *http.Request) {
// 	var fWorker lib.BusyWorkerSignal
// 	json.NewDecoder(r.Body).Decode(&fWorker)
// 	busyWorker(klasterStatus, fWorker, workersPool)
// }

// func exprHandle(w http.ResponseWriter, r *http.Request) {
// 	var clientReq lib.ClientTask
// 	json.NewDecoder(r.Body).Decode(&clientReq)
// 	resultOfExpr := exprDisturbToWorkes(workersPool, deferTasksPool, clientReq)
// 	resultOfExprJson, _ := json.Marshal(resultOfExpr)
	
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(resultOfExprJson)
// }

func caddworker(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Добавление воркера в кластер")

	workers.Handler_caddworker(r.Body, clusterInfo, workersPool)
}

func cfreeworker(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Воркер выполнил свою работу и сообщил кластеру")
}

func cbusyworker(w http.ResponseWriter, r *http.Request) {

}

func csolveproblem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запрос на решение от менеджера:")
	body, _ := io.ReadAll(r.Body)
	fmt.Println(string(body))

	mtrx := gn.GenerateRandMatrix(3, 3, 25)
	mtrxjson, _ := json.Marshal(mtrx)

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(mtrxjson)
}

func main() {
	args := os.Args
	if len(args) < 4 {
		panic("<порт><id><порт менеджера>")
	}

	clusterInfo, workersPool, deferClusterWorkerTaskPool = cinit.ClusterInit(args)
	
	
	http.HandleFunc("/caddworker", caddworker)
	http.HandleFunc("/cfreeworker", cfreeworker)
	http.HandleFunc("/cbusyworker", cbusyworker)

	http.HandleFunc("/csolveproblem", csolveproblem)

	http.ListenAndServe(fmt.Sprintf(":%s", clusterInfo.Port), nil)
}