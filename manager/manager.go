package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lib"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var managerPort string
var klasters map[int]*lib.Klaster
var CurrentKlaster *lib.Klaster

func addKlaster(w http.ResponseWriter, r *http.Request) {
	var newKlaster lib.Klaster
	json.NewDecoder(r.Body).Decode(&newKlaster)

	klasters[newKlaster.Id] = &newKlaster
	fmt.Println("New klaster add. Id:", newKlaster.Id)
}

func updateKlaster(w http.ResponseWriter, r *http.Request) {
	var updateKlaster lib.UpdateKlaster
	json.NewDecoder(r.Body).Decode(&updateKlaster)

	atomic.StoreUint32(&klasters[updateKlaster.Id].AllCores, updateKlaster.AllCores)
	atomic.StoreUint32(&klasters[updateKlaster.Id].FreeCores , updateKlaster.FreeCores)
}

func updateCurrentKlaster() {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		<-ticker.C

		CurrentKlaster = nil
		for _, v := range klasters {
			if CurrentKlaster == nil || atomic.LoadUint32(&v.FreeCores) > CurrentKlaster.FreeCores {
				CurrentKlaster = v
			}
		}
	}
}

func expr(w http.ResponseWriter, r *http.Request) {
	responseKlaster, err := http.Post(fmt.Sprintf("http://localhost:%s/expr", CurrentKlaster.Port), "application/json", r.Body)
	if err != nil {
		fmt.Println("Команда \"От винта\"!")
	}

	responseManagerJson, err := io.ReadAll(responseKlaster.Body)
	if err != nil {
		fmt.Println("Укуси меня пчела")
	}
	responseKlaster.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseManagerJson)
}

func getKlaster(w http.ResponseWriter, r *http.Request) {
	fmt.Println(klasters)
	fmt.Println()
	fmt.Println(CurrentKlaster)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Команда: go run manager.go <порт>")
		return
	}
	managerPort = args[1]

	klasters = map[int]*lib.Klaster{}
	go updateCurrentKlaster()
	
	
	http.HandleFunc("/expr", expr)
	
	http.HandleFunc("/addklaster", addKlaster)
	http.HandleFunc("/updateklaster", updateKlaster)
	http.HandleFunc("/diagnostic/klasterslist", getKlaster)
	
	
	fmt.Println("Manager started on port", managerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", managerPort), nil)
}