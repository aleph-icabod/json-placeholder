package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aleph-icabod/json-placeholder/apiClient"
	"github.com/aleph-icabod/json-placeholder/controller"
)

// Config helper struct for get the url for slack notification
type Config struct {
	SlackWebhook string
}

func main() {

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("error reading the configuration file " + err.Error())
	}
	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Fatal("error marshalling configuration file: " + err.Error())
	}

	photoApi := apiClient.NewClient("https://jsonplaceholder.typicode.com/photos")
	slack := apiClient.NewSlack(c.SlackWebhook)
	ctrl := controller.NewController(photoApi, slack)
	for {
		option := printMenu()
		var err error
		switch option {
		case 1:
			err = ctrl.CreatePhoto()
		case 2:
			err = ctrl.ListPhotos()
		case 3:
			err = ctrl.GetPhoto()
		case 4:
			err = ctrl.UpdatePhoto()
		case 5:
			err = ctrl.DeletePhoto()
		case 6:
			os.Exit(0)
		}
		if err != nil {
			fmt.Println("Error ", err)
		}
	}
}

// prinMenu helper to print the user menu
func printMenu() int {
	fmt.Println("Json placeholder API")
	fmt.Println("--Select an option--")
	fmt.Println("1 CREATE")
	fmt.Println("2 LIST")
	fmt.Println("3 GET")
	fmt.Println("4 UPDATE")
	fmt.Println("5 DELETE")
	fmt.Println("6 EXIT")
	var option int
	_, err := fmt.Scan(&option)
	if err != nil {
		return -1
	}
	if option < 1 || option > 6 {
		return -1
	}
	return option
}
