package main

import(
	"lib"
	"sync/atomic"
	"net/http"
	"bytes"
	"fmt"
	"encoding/json"
)

func addWorker(klasterStatus *lib.Klaster, newWorker *lib.Worker) {
	klasterStatus.WorkersList[newWorker.Id] = newWorker
	atomic.AddUint32(&klasterStatus.AllCores, uint32(newWorker.Cores))
	atomic.AddUint32(&klasterStatus.FreeCores, uint32(newWorker.Cores))

	updateServerInfo := lib.UpdateKlaster{
		Id: klasterStatus.Id,
		FreeCores: atomic.LoadUint32(&klasterStatus.FreeCores),
		AllCores: atomic.LoadUint32(&klasterStatus.AllCores),
	}
	updateServerInfoJson, _ := json.Marshal(updateServerInfo)
	http.Post(fmt.Sprintf("http://localhost:%s/updateklaster", klasterStatus.ManagerPort), "application/json", bytes.NewReader(updateServerInfoJson))
}