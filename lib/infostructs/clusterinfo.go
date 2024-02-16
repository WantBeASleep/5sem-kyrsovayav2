package infostructs

type ClusterInfo struct {
	Port        string
	ManagerPort string
	Id          int
	AllCores    int
	FreeCores   int
	WorkersList *map[int]*WorkerInfo
}
