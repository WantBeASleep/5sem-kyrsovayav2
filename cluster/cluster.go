//clearCommitHistory
package main

import (
	"os"
	"net/http"
	"fmt"

	is "lib/infostructs"
	ts "lib/tasks"

	"cluster/cinit"
	"cluster/taskhandler"
	"cluster/workers"
)

var clusterInfo *is.ClusterInfo
var workersPool chan *is.WorkerInfo
var deferClusterWorkerTaskPool chan ts.ClusterWorkerTask


func caddworker(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Добавление воркера в кластер")

	workers.Handler_caddworker(r.Body, clusterInfo, workersPool)
}

func cfreeworker(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Воркер выполнил свою работу и сообщил кластеру")

	workers.Handler_cfreeworker(r.Body, clusterInfo, workersPool)
}

func cbusyworker(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Воркер взял свою работу и сообщил кластеру")

	workers.Handler_cbusyworker(r.Body, clusterInfo, workersPool)
}

func csolveproblem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запрос на решение от менеджера")
	result := taskhandler.Handler_csolveproblem(workersPool, deferClusterWorkerTaskPool, r.Body)
	fmt.Println("Воркер решил задачу, отправка ответа менеджеру")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
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