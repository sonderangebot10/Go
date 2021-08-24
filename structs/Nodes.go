package structs

import . "task/helpers"

type Nodes []Node

type Node struct {
	Name      string  `json:"name"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Zone      string  `json:"zone"`
	Cost      float64 `json:"cost"`
	PodsNames []string
}

func (n Nodes) Copy() Nodes {

	result := Nodes{}

	for _, x := range n {
		n := struct {
			Name      string  `json:"name"`
			CPU       float64 `json:"cpu"`
			Memory    float64 `json:"memory"`
			Zone      string  `json:"zone"`
			Cost      float64 `json:"cost"`
			PodsNames []string
		}{x.Name, x.CPU, x.Memory, x.Zone, x.Cost, x.PodsNames}

		result = append(result, n)
	}

	return result
}

type Req struct {
	CPU    float64
	Memory float64
}

func (n Nodes) TotalRemainingReq(zone string) Req {

	result := Req{}

	for _, x := range n {

		if x.Zone == zone || zone == "" {
			result.CPU = ToFixed(result.CPU+x.CPU, 2)
			result.Memory = ToFixed(result.Memory+x.Memory, 2)
		}
	}

	return result
}

func (n Nodes) HasFitNodes(cpu float64, memory float64, zone string) bool {

	for _, x := range n {

		if x.Zone == zone || zone == "" {

			if x.CPU >= cpu && x.Memory >= memory {
				return true
			}
		}
	}

	return false
}

func UnorderedEqual(first, second Nodes) bool {
	if len(first) != len(second) {
		return false
	}
	exists := make(map[string]bool)
	for _, value := range first {
		exists[value.Name] = true
	}
	for _, value := range second {
		if !exists[value.Name] {
			return false
		}
	}
	return true
}
