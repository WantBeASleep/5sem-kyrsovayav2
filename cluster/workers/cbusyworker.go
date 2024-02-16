package workers

import (
	"io"
	is "lib/infostructs"
	rq "lib/requests"
	"encoding/json"
)

func Handler_cbusyworker(reqData io.ReadCloser, clusterInfo *is.ClusterInfo, workersPool chan *is.WorkerInfo) {
	var busyWorkerId int
	err := json.NewDecoder(reqData).Decode(&busyWorkerId)
	if err != nil {
		panic("Ошибка парса воркера на кластере, cbusyworker")
	}

	busyWorker := (*clusterInfo.WorkersList)[busyWorkerId]
	clusterInfo.FreeCores -= busyWorker.Cores
	rq.SendRequest(clusterInfo.ManagerPort, "mupdatecluster", clusterInfo)
}