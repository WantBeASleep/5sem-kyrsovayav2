package main

import (
	"fmt"
	"net/http"
	"os"

	"manager/client"
	"manager/cluster"
	"manager/diagnostic"

	is "lib/infostructs"
	mi "manager/minit"
)

var ManagerPort string
var Clusters map[int]*is.ClusterInfo
var freeClusterReq chan chan *is.ClusterInfo
var updateClusterReq chan *is.ClusterInfo

func maddcluster(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("Добавление кластера")

	cluster.Hander_maddcluster(updateClusterReq, r.Body)
}

func mupdatecluster(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Обновление кластера")

	cluster.Hander_mupdatecluster(updateClusterReq, r.Body)
}

func msolveproblem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Запрос на решение от клиента")

	exprResult := client.Handler_msolveproblem(r.Body, freeClusterReq)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(exprResult)
}

func showclusters(w http.ResponseWriter, r *http.Request) {
	diagnostic.ShowClusters(Clusters)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		panic("<порт>")
	}

	ManagerPort, Clusters, freeClusterReq, updateClusterReq = mi.ManagerInit(args)

	http.HandleFunc("/msolveproblem", msolveproblem)

	http.HandleFunc("/maddcluster", maddcluster)
	http.HandleFunc("/mupdatecluster", mupdatecluster)
	http.HandleFunc("/diagnostic/showclusters", showclusters)

	http.ListenAndServe(fmt.Sprintf(":%s", ManagerPort), nil)
}
