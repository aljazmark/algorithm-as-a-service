package algorithms

/*
#include "basic.h"
#include "sort.h"
#include "stdlib.h"
*/
import "C"
import (
	"algo/models"
)

//AlgHandler prepares response objecet and calls the package with requested algorithm
func AlgHandler(algorithm string, request models.AlgoRequest) models.AlgoResponse {
	var result models.AlgoResponse
	var alg Algorithm
	result.Algorithm = algorithm
	result.Input = request.Input
	result.Parameters = request.Parameters
	input := Parameters{algorithm, request.Input, request.Parameters}

	switch algorithm {
	//Sorting algorithms
	case "InsertionSort", "BubbleSort", "MergeSort", "SelectionSort", "QuickSort":
		alg = Sorting{}
		output, time, err := RunAlgorithm(alg, input)
		if err != nil {
			result.Output = "Problem running algorithm"
			result.ExecutionTime = ""
		} else {
			result.ExecutionTime = time
			result.Output = output
		}

	//K-center algorithms
	case "bf",
		"bfbrec",
		"reduce",
		"greedy",
		"greedyplus",
		"greedyrnd",
		"gonzalez1c",
		"gonzalezplus",
		"gonzalezrnd",
		"plesnikrnd",
		"plesnikdeg",
		"hochbaumshmoys",
		"hochbaumshmoysbin",
		"bottleneck",
		"score",
		"CDSP",
		"CDSh",
		"CDSPh",
		"CDS":

		alg = Kcenter{}
		output, time, err := RunAlgorithm(alg, input)
		if err != nil {
			result.Output = "Problem running algorithm"
			result.ExecutionTime = ""
		} else {
			result.ExecutionTime = time
			result.Output = output
		}

	}
	return result
}
