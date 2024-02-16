package cluster

import (
	rq "lib/requests"
	mt "lib/matrixes"
	gn "lib/generatelib"
	"io"
	"encoding/json"
)

func Handler_wsolveproblem(reqData io.ReadCloser) mt.Matrix {
	var task rq.ClusterWorkerReq
	err := json.NewDecoder(reqData).Decode(&task)
	if err != nil {
		panic("Ошибка распарса задачи на воркере от клиента!")
	}

	newMatrix := gn.GenerateRandMatrix(3, 3, 25)
	return newMatrix
}