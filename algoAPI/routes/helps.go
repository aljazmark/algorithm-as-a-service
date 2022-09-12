package routes

import (
	"algoAPI/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

//Helps defines Helps controller struct
type Helps struct {
	log  *log.Logger
	help models.Help
}

//NewHelps initializes a new Helps struct with given logger and database client
func NewHelps(log *log.Logger, help models.Help) *Helps {
	return &Helps{log: log, help: help}
}

//PrepHelp reads help.json
func (hel *Helps) PrepHelp() bool {
	helpJSON, err := os.Open("help.json")
	if err != nil {
		log.Println("Error opening help.json file: " + err.Error())
		return false
	}
	defer helpJSON.Close()
	helpBytes, err := ioutil.ReadAll(helpJSON)
	if err != nil {
		log.Println("Error reading help.json file: " + err.Error())
		return false
	}
	err = json.Unmarshal(helpBytes, &hel.help)
	if err != nil {
		log.Println("Error parsing help.json file: " + err.Error())
		return false
	}
	return true
}

//GetAlgorithms returns all available algorithms
func (hel *Helps) GetAlgorithms(c echo.Context) error {
	algos, _ := json.Marshal(struct {
		Algorithms []string `json:"algorithms"`
	}{hel.help.Algorithms})
	return c.JSON(http.StatusOK, json.RawMessage(algos))
}

//GetAlgoHelp returns algorithm details for requested algorithm
func (hel *Helps) GetAlgoHelp(c echo.Context) error {
	algoName := c.Param("algorithm")
	for _, a := range hel.help.Help {
		if a.Algorithm == algoName {
			return c.JSON(http.StatusOK, a)
		}
	}
	return c.JSON(http.StatusBadRequest, json.RawMessage(`{"message": "Algorithm not found"}`))
}
