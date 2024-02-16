package cinit

import (
	"cluster/taskhandler"
	is "lib/infostructs"
	rq "lib/requests"
	ts "lib/tasks"

	"strconv"
)

func ClusterInit(args []string) (*is.ClusterInfo, chan *is.WorkerInfo, chan ts.ClusterWorkerTask) {
	clusterPort := args[1]
	clusterId, err := strconv.Atoi(args[2])
	if err != nil {
		panic("Ошибка парса ID кластера")
	}
	managerPort := args[3]

	clusterInfo := &is.ClusterInfo{
		Port: clusterPort,
		ManagerPort: managerPort,
		Id: clusterId,
		AllCores: 0,
		FreeCores: 0,
		WorkersList: &map[int]*is.WorkerInfo{},
	}

	workersPool := make(chan *is.WorkerInfo, 50)
	deferClusterWorkerTaskPool := make(chan ts.ClusterWorkerTask, 50)
	rq.SendRequest(managerPort, "maddcluster", clusterInfo)
	go taskhandler.DeferTasksPoolHandler(deferClusterWorkerTaskPool, workersPool)
	return clusterInfo, workersPool, deferClusterWorkerTaskPool
}