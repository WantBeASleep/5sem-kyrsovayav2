package cluster

import (
	is "lib/infostructs"
	"io"
	"encoding/json"
)

func Hander_mupdatecluster(updateClusterReq chan *is.ClusterInfo, reqData io.ReadCloser) {
	var updateCluster is.ClusterInfo
	err := json.NewDecoder(reqData).Decode(&updateCluster)
	if err != nil {
		panic("Ошибка распарса кластера при обновлении")
	}
	updateClusterReq <- &updateCluster
}