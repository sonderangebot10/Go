package structs

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

func RemoveUnusedNodesPrice(cost float64, n Nodes) float64 {

	for _, v := range n {
		if len(v.PodsNames) == 0 {
			cost -= v.Cost
		}
	}

	return cost
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
