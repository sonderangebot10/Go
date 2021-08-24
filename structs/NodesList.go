package structs

import "fmt"

type NodesList []NodesListItem

type NodesListItem struct {
	NodeList  Nodes
	TotalCost float64
}

func (n NodesListItem) Print() {

	fmt.Print(n.TotalCost)
	fmt.Print(": ")
	for _, v := range n.NodeList {
		if len(v.PodsNames) > 0 {
			fmt.Println()
			fmt.Print(v.Name + ":")

			for _, n := range v.PodsNames {
				fmt.Print(" " + n)
			}
		}
	}
	fmt.Println()
}

func (n NodesListItem) PrintPrice() {

	fmt.Print(n.TotalCost)
	fmt.Println()
}
