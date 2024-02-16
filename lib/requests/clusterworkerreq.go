package requests

import (
	mt "lib/matrixes"
	tr "lib/trees"
)

type ClusterWorkerReq struct {
	Root     tr.ASTNode
	Matrixes map[string]mt.Matrix
}