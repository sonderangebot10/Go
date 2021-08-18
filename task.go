package main

import (
	"fmt"
	"math"
	"strconv"
	. "task/helpers"
	. "task/structs"
)

var allVariationsNodes = NodesList{}
var globalMinPrice = math.MaxFloat64

var globalLowestCostSolution = NodesListItem{}
var found = false

func main() {

	var pods Pods
	GenericReadFile("tests/pods.json", &pods)

	var nodes Nodes
	GenericReadFile("tests/nodes.json", &nodes)

	totals := TotalReq{}
	for _, v := range pods {

		if v.Zone == "A" {
			totals.A_totalCPU = ToFixed(v.CPURequest+totals.A_totalCPU, 2)
			totals.A_totalMemory = ToFixed(v.MemoryRequest+totals.A_totalMemory, 2)
		} else if v.Zone == "B" {
			totals.B_totalCPU = ToFixed(v.CPURequest+totals.B_totalCPU, 2)
			totals.B_totalMemory = ToFixed(v.MemoryRequest+totals.B_totalMemory, 2)
		} else if v.Zone == "" {
			totals.X_totalCPU = ToFixed(v.CPURequest+totals.X_totalCPU, 2)
			totals.X_totalMemory = ToFixed(v.MemoryRequest+totals.X_totalMemory, 2)
		}
	}

	var tree = Leaf{Node: Node{}, TotalPrice: 0}

	queue := make([]Leaf, 0)

	for _, v := range nodes {

		t := totals
		if v.Zone == "A" {
			t.A_totalCPU = ToFixed(t.A_totalCPU-v.CPU, 2)
			t.A_totalMemory = ToFixed(t.A_totalMemory-v.Memory, 2)
		}
		if v.Zone == "B" {
			t.B_totalCPU = ToFixed(t.B_totalCPU-v.CPU, 2)
			t.B_totalMemory = ToFixed(t.B_totalMemory-v.Memory, 2)
		}

		var l = Leaf{Node: v, TotalPrice: v.Cost, Totals: t}
		tree.Parent = nil

		queue = append(queue, l)
	}

	traverseTree(queue, nodes)

	fmt.Println("OK")
}

func traverseTree(_queue []Leaf, nodes Nodes) {

	queue := _queue
	level := 1

	for len(queue) > 0 {

		size := len(queue)

		fmt.Println("checking level " + strconv.Itoa(level))
		level++

		for i := 0; i < size; i++ {

			u := queue[0]
			queue = queue[1:]

			if u.Totals.AllSatisfied() {
				// ENDING LEAF POINTING TO PARENT TO FORM A POSSIBILITY PATH
			}

			if !u.Totals.AllSatisfied() {
				for _, node := range nodes {

					t := u.Totals
					if node.Zone == "A" && (t.A_totalCPU > 0 || t.A_totalMemory > 0) {
						t.A_totalCPU = ToFixed(t.A_totalCPU-node.CPU, 2)
						t.A_totalMemory = ToFixed(t.A_totalMemory-node.Memory, 2)

						var l = Leaf{Node: node, TotalPrice: ToFixed(u.TotalPrice+node.Cost, 15), Totals: t, Parent: &u}
						queue = append(queue, l)

					} else if node.Zone == "B" && (t.B_totalCPU > 0 || t.B_totalMemory > 0) {
						t.B_totalCPU = ToFixed(t.B_totalCPU-node.CPU, 2)
						t.B_totalMemory = ToFixed(t.B_totalMemory-node.Memory, 2)

						var l = Leaf{Node: node, TotalPrice: ToFixed(u.TotalPrice+node.Cost, 15), Totals: t, Parent: &u}
						queue = append(queue, l)

					} else if t.X_totalCPU > 0 || t.X_totalMemory > 0 {
						t.X_totalCPU = ToFixed(t.X_totalCPU-node.CPU, 2)
						t.X_totalMemory = ToFixed(t.X_totalMemory-node.Memory, 2)

						var l = Leaf{Node: node, TotalPrice: ToFixed(u.TotalPrice+node.Cost, 15), Totals: t, Parent: &u}
						queue = append(queue, l)
					}
				}
			}
		}
	}
}

func permute(pods Pods, currentSolution Nodes, count int, price float64, startPods int) {

	for i := startPods; i < len(pods); i++ {

		for j := 0; j < len(currentSolution); j++ {

			if currentSolution[j].CPU >= pods[i].CPURequest && currentSolution[j].Memory >= pods[i].MemoryRequest && (currentSolution[j].Zone == pods[i].Zone || currentSolution[j].Zone == "" || pods[i].Zone == "") {

				newNode := Node{
					Name:      currentSolution[j].Name,
					CPU:       ToFixed(currentSolution[j].CPU-pods[i].CPURequest, 2),
					Memory:    ToFixed(currentSolution[j].Memory-pods[i].MemoryRequest, 2),
					Zone:      currentSolution[j].Zone,
					Cost:      currentSolution[j].Cost,
					PodsNames: append(currentSolution[j].PodsNames, pods[i].Name)}
				oldNode := Node{
					Name:      currentSolution[j].Name,
					CPU:       currentSolution[j].CPU,
					Memory:    currentSolution[j].Memory,
					Zone:      currentSolution[j].Zone,
					Cost:      currentSolution[j].Cost,
					PodsNames: currentSolution[j].PodsNames}

				currentSolution[j] = newNode

				permute(pods, currentSolution, count+1, price, i+1)

				currentSolution[j] = oldNode
			}
		}
	}

	if count == len(pods) && !found {

		globalLowestCostSolution = NodesListItem{NodeList: currentSolution.Copy(), TotalCost: price}
		globalLowestCostSolution.Print()
	}
}
