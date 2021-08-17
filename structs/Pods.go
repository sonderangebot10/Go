package structs

type Pods []Pod

type Pod struct {
	Name          string  `json:"name"`
	CPURequest    float64 `json:"cpuRequest"`
	MemoryRequest float64 `json:"memoryRequest"`
	Zone          string  `json:"zone"`
}
