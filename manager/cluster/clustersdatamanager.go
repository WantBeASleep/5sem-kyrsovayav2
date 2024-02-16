package cluster

import (
	is "lib/infostructs"
)

func ClusterDataManager(clusters *map[int]*is.ClusterInfo, freeClusterReq chan chan *is.ClusterInfo, updateClusterReq chan *is.ClusterInfo) {
	for {
		select {
		case x := <- freeClusterReq:
			x <- getFreeCluster(clusters)
		case x := <- updateClusterReq:
			updateCluster := x
			(*clusters)[updateCluster.Id] = updateCluster
		}
		
	}
}