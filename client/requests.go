package main

import (
	"lib"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"math/rand"
)

func getFormedTask() (string, map[string]lib.Matrix) {
	expr := "(a+b)*(c+d)"
	m, n := 1 + rand.Intn(10), 1 + rand.Intn(10)

	a := lib.GenerateRandMatrix(m, n, 100)
	b := lib.GenerateRandMatrix(m, n, 100)
	c := lib.GenerateRandMatrix(m, n, 100)
	d := lib.GenerateRandMatrix(m, n, 100)
	
	h := map[string]lib.Matrix{
		"a": a,
		"b": b,
		"c": c,
		"d": d,
	}

	return expr, h
}

func SendRequest(port string) {
	expr, data := getFormedTask()
	newTask := lib.ClientTask{
		Expr: expr,
		Data: data,
	}

	newTaskJson, _ := json.Marshal(newTask)
	fmt.Println("client send:", string(newTaskJson))

	response, err := http.Post(fmt.Sprintf("http://localhost:%s/expr", port) , "application/json", bytes.NewReader(newTaskJson))
	if err != nil {
		panic("zxc tilt")
	}
	defer response.Body.Close()

	var result lib.Matrix
	json.NewDecoder(response.Body).Decode(&result)
	fmt.Println("client get:", result)
}
