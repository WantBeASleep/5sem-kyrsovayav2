package cluster

import (
	"encoding/json"
	"io"
	is "lib/infostructs"
)

func Hander_maddcluster(updateClusterReq chan *is.ClusterInfo, reqData io.ReadCloser) {
	var newCluster is.ClusterInfo
	err := json.NewDecoder(reqData).Decode(&newCluster)
	if err != nil {
		panic("Ошибка распарса кластера при добавлении")
	}
	updateClusterReq <- &newCluster
}