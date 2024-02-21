package main

import (
	"fmt"
	gn "lib/generatelib"
	mt "lib/matrixes"
	rq "lib/requests"
)

const (
	m int = 1500
	n int = 1200
)

func getFormedTask() rq.ClientReq {
	expr := "(a + b * c) * (d + e * f)"
	a := gn.GenerateRandMatrix(m, n, 100)
	b := gn.GenerateRandMatrix(m, n, 100)
	c := gn.GenerateRandMatrix(n, m, 100)
	d := gn.GenerateRandMatrix(m, n, 100)
	e := gn.GenerateRandMatrix(m, n, 100)
	f := gn.GenerateRandMatrix(n, m, 100)
	
	h := map[string]mt.Matrix{
		"a": a,
		"b": b,
		"c": c,
		"d": d,
		"e": e,
		"f": f,
	}

	req := rq.ClientReq{
		Expr: expr,
		Matrixes: h,
	}
	return req
}

func SendRequest(port string) {
	Req := getFormedTask()
	var Ans mt.Matrix
	rq.SendRequest(port, "msolveproblem", Req, &Ans)
	fmt.Println(Ans)
}