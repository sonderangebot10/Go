package structs

type Pods []Pod

type Pod struct {
	Name          string  `json:"name"`
	CPURequest    float64 `json:"cpuRequest"`
	MemoryRequest float64 `json:"memoryRequest"`
	Zone          string  `json:"zone"`
}

type TotalReq struct {
	A_totalCPU    float64
	A_totalMemory float64

	B_totalCPU    float64
	B_totalMemory float64

	X_totalCPU    float64
	X_totalMemory float64
}

func (t TotalReq) AllSatisfied() bool {
	return t.A_totalCPU <= 0 && t.A_totalMemory <= 0 && t.B_totalCPU <= 0 && t.B_totalMemory <= 0 && t.X_totalCPU <= 0 && t.X_totalMemory <= 0
}
