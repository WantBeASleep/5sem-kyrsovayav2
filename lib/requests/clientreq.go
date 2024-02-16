package requests

import (
	mt "lib/matrixes"
)

type ClientReq struct {
	Expr     string
	Matrixes map[string]mt.Matrix
}
