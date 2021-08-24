package main

import (
	"fmt"
	"math"
	"sort"
	. "task/helpers"
	. "task/structs"
)

const MAX_RETRY_COUNT = 1000000

var globalMinPrice = math.MaxFloat64
var globalLowestCostSolution = NodesListItem{}

func main() {

	var pods Pods
	GenericReadFile("data/pods.json", &pods)

	var nodes Nodes
	GenericReadFile("data/nodes.json", &nodes)

	sort.Slice(nodes[:], func(i, j int) bool {
		return (nodes[i].CPU+nodes[i].Memory)/nodes[i].Cost > (nodes[j].CPU+nodes[j].Memory)/nodes[j].Cost
	})

	totalRequirements := TotalReq{}
	for _, v := range pods {
		if v.Zone == "A" {
			totalRequirements.A_totalCPU = ToFixed(v.CPURequest+totalRequirements.A_totalCPU, 2)
			totalRequirements.A_totalMemory = ToFixed(v.MemoryRequest+totalRequirements.A_totalMemory, 2)
		} else if v.Zone == "B" {
			totalRequirements.B_totalCPU = ToFixed(v.CPURequest+totalRequirements.B_totalCPU, 2)
			totalRequirements.B_totalMemory = ToFixed(v.MemoryRequest+totalRequirements.B_totalMemory, 2)
		} else if v.Zone == "" {
			totalRequirements.X_totalCPU = ToFixed(v.CPURequest+totalRequirements.X_totalCPU, 2)
			totalRequirements.X_totalMemory = ToFixed(v.MemoryRequest+totalRequirements.X_totalMemory, 2)
		}
	}

	StartTimePrinter()

	explore(nodes, pods, Nodes{}, 0, 0, totalRequirements, 0)

	fmt.Print("cheapest solution found: ")
	globalLowestCostSolution.Print()
}

func explore(nodes Nodes, pods Pods, currentSolution Nodes, count int, minPrice float64, t TotalReq, ni int) {

	if minPrice >= globalMinPrice {
		return
	}

	if t.AllSatisfied() {

		retries = 0
		tryPack(pods, currentSolution.Copy(), 0, minPrice, 0)
		return
	}

	for nodei, node := range nodes[ni:] {

		tmpT := t
		if node.Zone == "A" && ((t.A_totalCPU > 0 || t.A_totalMemory > 0) || (t.X_totalCPU > 0 || t.X_totalMemory > 0)) {
			tmpT.A_totalCPU = ToFixed(tmpT.A_totalCPU-node.CPU, 2)
			tmpT.A_totalMemory = ToFixed(tmpT.A_totalMemory-node.Memory, 2)

			if tmpT.A_totalCPU < 0 && tmpT.A_totalMemory < 0 {
				tmpT.X_totalCPU = ToFixed(tmpT.X_totalCPU+tmpT.A_totalCPU, 2)
				tmpT.X_totalMemory = ToFixed(tmpT.X_totalMemory+tmpT.A_totalMemory, 2)

				tmpT.A_totalCPU = 0
				tmpT.A_totalMemory = 0
			}
		} else if node.Zone == "B" && ((tmpT.B_totalCPU > 0 || t.B_totalMemory > 0) || (t.X_totalCPU > 0 || t.X_totalMemory > 0)) {
			tmpT.B_totalCPU = ToFixed(tmpT.B_totalCPU-node.CPU, 2)
			tmpT.B_totalMemory = ToFixed(tmpT.B_totalMemory-node.Memory, 2)

			if tmpT.B_totalCPU < 0 && tmpT.B_totalMemory < 0 {
				tmpT.X_totalCPU = ToFixed(tmpT.X_totalCPU+tmpT.B_totalCPU, 2)
				tmpT.X_totalMemory = ToFixed(tmpT.X_totalMemory+tmpT.B_totalMemory, 2)

				tmpT.B_totalCPU = 0
				tmpT.B_totalMemory = 0
			}
		} else {
			continue
		}

		currentSolution = append(currentSolution, node)
		explore(nodes, pods, currentSolution, count+1, minPrice+node.Cost, tmpT, nodei)
		currentSolution = currentSolution[:len(currentSolution)-1]
	}
}

var retries int

func tryPack(pods Pods, currentSolution Nodes, count int, minPrice float64, startPods int) bool {

	if retries > MAX_RETRY_COUNT {
		return false
	}

	for i := startPods; i < len(pods); i++ {

		if !currentSolution.HasFitNodes(pods[i].CPURequest, pods[i].MemoryRequest, pods[i].Zone) ||
			(pods[i].CPURequest > currentSolution.TotalRemainingReq(pods[i].Zone).CPU ||
				pods[i].MemoryRequest > currentSolution.TotalRemainingReq(pods[i].Zone).Memory) ||
			count != i {

			return false
		}

		for j := 0; j < len(currentSolution); j++ {

			if currentSolution[j].CPU >= pods[i].CPURequest &&
				currentSolution[j].Memory >= pods[i].MemoryRequest &&
				(currentSolution[j].Zone == pods[i].Zone || pods[i].Zone == "") {

				currentSolution[j].CPU = ToFixed(currentSolution[j].CPU-pods[i].CPURequest, 2)
				currentSolution[j].Memory = ToFixed(currentSolution[j].Memory-pods[i].MemoryRequest, 2)
				currentSolution[j].PodsNames = append(currentSolution[j].PodsNames, pods[i].Name)

				if tryPack(pods, currentSolution, count+1, minPrice, i+1) {
					return true
				}

				retries++

				currentSolution[j].CPU = ToFixed(currentSolution[j].CPU+pods[i].CPURequest, 2)
				currentSolution[j].Memory = ToFixed(currentSolution[j].Memory+pods[i].MemoryRequest, 2)
				currentSolution[j].PodsNames = currentSolution[j].PodsNames[:len(currentSolution[j].PodsNames)-1]
			}
		}

	}

	if count == len(pods) {

		if globalMinPrice > minPrice {

			globalMinPrice = minPrice

			globalLowestCostSolution = NodesListItem{NodeList: currentSolution.Copy(), TotalCost: minPrice}
			fmt.Print("new lowest price found: ")
			globalLowestCostSolution.PrintPrice()

			return true
		}
	}

	return false
}
