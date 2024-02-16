package minit

import (
	is "lib/infostructs"
	mc "manager/cluster"
)

func ManagerInit(args []string) (string, map[int]*is.ClusterInfo, chan chan *is.ClusterInfo, chan *is.ClusterInfo) {
	managerPort := args[1]
	clusters := map[int]*is.ClusterInfo{}
	freeClusterReq := make(chan chan *is.ClusterInfo, 100)
	updateClusterReq := make(chan *is.ClusterInfo, 100)
	go mc.ClusterDataManager(&clusters, freeClusterReq, updateClusterReq)
	return managerPort, clusters, freeClusterReq, updateClusterReq
}