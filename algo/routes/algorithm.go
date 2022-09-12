package routes

import (
	"algo/algorithms"
	"algo/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

//List of available algorithms
var algoList = []string{
	"InsertionSort",
	"BubbleSort",
	"SelectionSort",
	"QuickSort",
	"MergeSort",
	"bf",
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
	"CDS"}

//GetAlgorithms returns available algorithms
func GetAlgorithms(c echo.Context) error {
	return c.JSON(http.StatusOK, algoList)
}

//RunAlgorithm checks if algorithm is available and passes request
func RunAlgorithm(c echo.Context) error {
	u := models.AlgoRequest{}
	r := models.AlgoResponse{}
	responseCode := http.StatusInternalServerError
	algorithm := c.Param("algo")
	if err := c.Bind(&u); err != nil {
		responseCode = http.StatusBadRequest
		return c.JSON(responseCode, err.Error)
	}
	if algoCheck(algorithm) {
		r = algorithms.AlgHandler(algorithm, u)
		responseCode = http.StatusOK
	} else {
		r.Output = "Algorithm not found"
		responseCode = http.StatusBadRequest
	}
	return c.JSON(responseCode, r)
}

//Checks if algorithm is available
func algoCheck(algo string) bool {
	for _, i := range algoList {
		if i == algo {
			return true
		}
	}
	return false
}
