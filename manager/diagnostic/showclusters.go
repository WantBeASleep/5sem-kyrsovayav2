package diagnostic

import (
	"fmt"
	is "lib/infostructs"
)

func ShowClusters(Clusters map[int]*is.ClusterInfo) {
	fmt.Println(Clusters)
}