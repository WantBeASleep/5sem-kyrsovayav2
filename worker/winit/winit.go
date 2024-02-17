package winit

import (
	is "lib/infostructs"
	rq "lib/requests"
	"strconv"
	"worker/workpool"
)

func WorkerInit(args []string) *is.WorkerInfo {
	workerPort := args[1]
	clusterPort := args[2]
	id, err := strconv.Atoi(args[3])
	if err != nil {
		panic("Ошибка парса Id")
	}
	cores, err := strconv.Atoi(args[4])
	if err != nil {
		panic("Ошибка парса кол-во ядер")
	}

	workerInfo := is.WorkerInfo{
		Port: workerPort,
		ClusterPort: clusterPort,
		Id: id,
		Cores: cores,
	}

	workpool.StartPool()
	rq.SendRequest(workerInfo.ClusterPort, "caddworker", workerInfo)
	return &workerInfo
}