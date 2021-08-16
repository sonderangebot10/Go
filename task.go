package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

type Pods []struct {
	Name          string  `json:"name"`
	CPURequest    float64 `json:"cpuRequest"`
	MemoryRequest float64 `json:"memoryRequest"`
	Zone          string  `json:"zone"`
}

type Nodes []struct {
	Name      string  `json:"name"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Zone      string  `json:"zone"`
	Cost      float64 `json:"cost"`
	PodsNames []string
}

func main() {

	// PODS READ
	podsJsonFile, err := os.Open("pods.json")

	var pods Pods

	podsByteValue, _ := ioutil.ReadAll(podsJsonFile)
	json.Unmarshal(podsByteValue, &pods)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer podsJsonFile.Close()

	// NODES READ
	nodesJsonFile, err := os.Open("nodes1.json")

	var nodes Nodes

	nodesByteValue, _ := ioutil.ReadAll(nodesJsonFile)
	json.Unmarshal(nodesByteValue, &nodes)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer nodesJsonFile.Close()

	explore(nodes, pods, Nodes{}, 0, 0, 0)
}

// THE LOGIC
var globalMinPrice = math.MaxFloat64

func explore(nodes Nodes, pods Pods, currentSolution Nodes, count, startPods int, minPrice float64) {

	for i := startPods; i < len(pods); i++ {

		if len(currentSolution) > 0 {
			for j := 0; j < len(currentSolution); j++ {

				if currentSolution[j].CPU >= pods[i].CPURequest && currentSolution[j].Memory >= pods[i].MemoryRequest && (currentSolution[j].Zone == pods[i].Zone || currentSolution[j].Zone == "" || pods[i].Zone == "") {
					currentSolution[j].CPU = currentSolution[j].CPU - pods[i].CPURequest
					currentSolution[j].Memory = toFixed(currentSolution[j].Memory-pods[i].MemoryRequest, 2)
					currentSolution[j].PodsNames = append(currentSolution[j].PodsNames, pods[i].Name)

					explore(nodes, pods, currentSolution, count+1, i+1, minPrice)

					currentSolution[j].CPU = currentSolution[j].CPU + pods[i].CPURequest
					currentSolution[j].Memory = toFixed(currentSolution[j].Memory+pods[i].MemoryRequest, 2)
					currentSolution[j].PodsNames = currentSolution[j].PodsNames[:len(currentSolution[j].PodsNames)-1]
				}
			}
		}

		for _, node := range nodes {

			if node.CPU >= pods[i].CPURequest && node.Memory >= pods[i].MemoryRequest {

				if node.Zone == pods[i].Zone || node.Zone == "" || pods[i].Zone == "" {

					newNode := node
					newNode.CPU = toFixed(newNode.CPU-pods[i].CPURequest, 2)
					newNode.Memory = toFixed(newNode.Memory-pods[i].MemoryRequest, 2)
					newNode.PodsNames = append(newNode.PodsNames, pods[i].Name)

					currentSolution = append(currentSolution, newNode)

					explore(nodes, pods, currentSolution, count+1, i+1, minPrice+newNode.Cost)

					currentSolution = currentSolution[:len(currentSolution)-1]
				}
			}
		}
	}

	if count == len(pods) && globalMinPrice > minPrice {
		globalMinPrice = minPrice

		fmt.Print(minPrice)
		fmt.Print(": ")
		for _, v := range currentSolution {
			fmt.Print(" " + v.Name + ":")

			for _, n := range v.PodsNames {
				fmt.Print(" " + n)
			}
		}
		fmt.Println()
	}
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
