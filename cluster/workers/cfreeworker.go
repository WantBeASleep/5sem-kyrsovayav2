package workers

import (
	"io"
	is "lib/infostructs"
	rq "lib/requests"
	"encoding/json"
)

func Handler_cfreeworker(reqData io.ReadCloser, clusterInfo *is.ClusterInfo, workersPool chan *is.WorkerInfo) {
	var freeWorkerId int
	err := json.NewDecoder(reqData).Decode(&freeWorkerId)
	if err != nil {
		panic("Ошибка парса воркера на кластере, cfreeworker")
	}

	freeWorker := (*clusterInfo.WorkersList)[freeWorkerId]
	clusterInfo.FreeCores += freeWorker.Cores
	workersPool <- freeWorker
	rq.SendRequest(clusterInfo.ManagerPort, "mupdatecluster", clusterInfo)
}