package lib

type Klaster struct {
	Port string
	ManagerPort string
	Id int
	AllCores uint32
	FreeCores uint32
	WorkersList map[int]*Worker
}

type UpdateKlaster struct {
	Id int
	FreeCores uint32
	AllCores uint32
}