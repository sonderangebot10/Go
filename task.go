package main

import (
	"fmt"
	"math"
	"strconv"
	. "task/helpers"
	. "task/structs"
	"time"
)

var globalMinPrice = math.MaxFloat64
var globalLowestCostSolution = NodesListItem{}

var foundSolutions = NodesList{}
var allCheckedSolutions = NodesList{}

func main() {

	var pods Pods
	GenericReadFile("tests/pods1.json", &pods)

	var nodes Nodes
	GenericReadFile("tests/nodes1.json", &nodes)

	explore(nodes, pods, Nodes{}, 0, 0, 0)

	fmt.Println("found " + strconv.Itoa((len(foundSolutions))) + " different possible node combinations")

	time.Sleep(3 * time.Second)

	for i, v := range foundSolutions {
		fmt.Println("Permuting " + strconv.Itoa(i) + " element")
		permute(pods, v.NodeList, 0, v.TotalCost, 0)
	}

	globalLowestCostSolution.Print()
}

func explore(nodes Nodes, pods Pods, currentSolution Nodes, count, startPods int, minPrice float64) {

	/// /// /// /// /// ///
	for _, v := range append(foundSolutions, allCheckedSolutions...) {
		if UnorderedEqual(v.NodeList, currentSolution) {
			return
		}
	}
	/// /// /// /// /// ///

	for i := startPods; i < len(pods); i++ {

		for _, node := range nodes {

			if node.CPU >= pods[i].CPURequest && node.Memory >= pods[i].MemoryRequest {

				if node.Zone == pods[i].Zone || node.Zone == "" || pods[i].Zone == "" {

					node.CPU = ToFixed(node.CPU-pods[i].CPURequest, 2)
					node.Memory = ToFixed(node.Memory-pods[i].MemoryRequest, 2)
					// node.PodsNames = append(newNode.PodsNames, pods[i].Name)

					currentSolution = append(currentSolution, node)
					explore(nodes, pods, currentSolution, count+1, i+1, minPrice+node.Cost)
					currentSolution = currentSolution[:len(currentSolution)-1]
				}
			}
		}

	}

	/// /// /// /// /// ///
	contains := false

	if count == len(pods) {

		for _, v := range foundSolutions {
			if UnorderedEqual(v.NodeList, currentSolution) {
				contains = true
			}
		}

		if !contains {

			fmt.Println("Found " + strconv.Itoa(len(foundSolutions)) + " solutions without permutation")

			foundSolutions = append(foundSolutions, NodesListItem{NodeList: currentSolution.Copy(), TotalCost: minPrice})
			return
		}
	}

	if !contains {

		allCheckedSolutions = append(allCheckedSolutions, NodesListItem{NodeList: currentSolution.Copy(), TotalCost: minPrice})
		return
	}
	/// /// /// /// /// ///
}

func permute(pods Pods, currentSolution Nodes, count int, minPrice float64, startPods int) {

	for i := startPods; i < len(pods); i++ {

		for j := 0; j < len(currentSolution); j++ {

			if currentSolution[j].CPU >= pods[i].CPURequest && currentSolution[j].Memory >= pods[i].MemoryRequest && (currentSolution[j].Zone == pods[i].Zone || currentSolution[j].Zone == "" || pods[i].Zone == "") {
				currentSolution[j].CPU = currentSolution[j].CPU - pods[i].CPURequest
				currentSolution[j].Memory = ToFixed(currentSolution[j].Memory-pods[i].MemoryRequest, 2)
				currentSolution[j].PodsNames = append(currentSolution[j].PodsNames, pods[i].Name)

				permute(pods, currentSolution, count+1, minPrice, i+1)

				currentSolution[j].CPU = currentSolution[j].CPU + pods[i].CPURequest
				currentSolution[j].Memory = ToFixed(currentSolution[j].Memory+pods[i].MemoryRequest, 2)
				currentSolution[j].PodsNames = currentSolution[j].PodsNames[:len(currentSolution[j].PodsNames)-1]
			}
		}
	}

	if count == len(pods) {

		minPrice = RemoveUnusedNodesPrice(minPrice, currentSolution)

		if globalMinPrice > minPrice {
			globalMinPrice = minPrice

			globalLowestCostSolution = NodesListItem{NodeList: currentSolution.Copy(), TotalCost: minPrice}
			globalLowestCostSolution.Print()
		}
	}
}
