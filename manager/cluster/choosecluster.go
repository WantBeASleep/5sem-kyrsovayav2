package cluster

import (
	is "lib/infostructs"
)

func getFreeCluster(Clusters *map[int]*is.ClusterInfo) *is.ClusterInfo {
	var ans *is.ClusterInfo
	for _, v := range *Clusters {
		if ans == nil || ans.FreeCores < v.FreeCores {
			ans = v
		}
	}
	return ans
}