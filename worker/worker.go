//clearCommitHistory
package main

import (
	"encoding/json"
	rq "lib/requests"
	tr "lib/trees"
	"fmt"
	is "lib/infostructs"
	"net/http"
	"os"
	"worker/winit"

	"worker/workpool"
	mt "lib/matrixes"
)

var workerInfo *is.WorkerInfo

func wsolveproblem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запрос на решение от кластера")
	rq.SendRequest(workerInfo.ClusterPort, "cbusyworker", workerInfo.Id)
	
	type uncodeHelper struct {
		Root json.RawMessage
		Matrixes map[string]mt.Matrix
	}
	var partReq uncodeHelper

	err := json.NewDecoder(r.Body).Decode(&partReq)
	if err != nil {
		panic("Ошибка распарса запроса на воркере!")
	}

	root := tr.UpparseJson(partReq.Root)

	ans := workpool.Solve(root, partReq.Matrixes)
	ansjs, err := json.Marshal(ans)
	if err != nil {
		panic("Ошибка парса на воркере при отправке ответа!")
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
