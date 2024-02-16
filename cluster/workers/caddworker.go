package workers

import (
	"io"
	is "lib/infostructs"
	rq "lib/requests"
	"encoding/json"
)

func Handler_caddworker(reqData io.ReadCloser, clusterInfo *is.ClusterInfo, workersPool chan *is.WorkerInfo) {
	var newWorker is.WorkerInfo
	err := json.NewDecoder(reqData).Decode(&newWorker)
	if err != nil {
		panic("Ошибка парса воркера на кластере caddworker")
	}

	clusterInfo.AllCores += newWorker.Cores
	clusterInfo.FreeCores += newWorker.Cores
	(*clusterInfo.WorkersList)[newWorker.Id] = &newWorker
	workersPool <- &newWorker
	rq.SendRequest(clusterInfo.ManagerPort, "mupdatecluster", clusterInfo)
}