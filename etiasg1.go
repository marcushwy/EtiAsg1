package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/passenger", allPassenger).Methods("GET")
	router.HandleFunc("/api/v1/create", createPassenger).Methods("POST")
	fmt.Printf("Listening at Port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

//Defining the structs

type Passenger struct {
	PFirstName    string `json:"PFirstName"`
	PLastName     string `json:"PLastName"`
	PMobileNumber int    `json:"PMobileNumber"`
	PEmailAddr    string `json:"PEmailAddr"`
}

var passengers map[string]Passenger = map[string]Passenger{

	"1": Passenger{"Lee", "Wei Yang", 82207177, "lwy@hotmail.com"},
	"2": Passenger{"See", "Dong Hua", 82207097, "sdh@hotmail.com"},
	"3": Passenger{"Lim", "Wen Song", 19807177, "lws@hotmail.com"},
}

type Driver struct {
	DFirstName    string `json:"FirstName"`
	DLastName     string `json:"LastName"`
	DMobileNumber int    `json:"MobileNumber"`
	DEmailAddr    string `json:"EmailAddr"`
	DIcNumber     string `json:"DIcNumber"`
	DCarPlate     string `json:"DCarPlate`
}

type Trip struct {
	PickUp           int `json:"PickUp"`
	DropOff          int `json:"DropOff"`

}

func allPassenger(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Hello from API")
	json.NewEncoder(w).Encode(passengers)

}

func createPassenger(w http.ResponseWriter, r *http.Request) {

	//get the body of the POST request
	//return te string response containing the request body

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))

	//unmarshal into new Passenger struct
	var newpassenger Passenger
	json.Unmarshal(reqBody, &newpassenger)

	fmt.Print(&newpassenger)
	//update global map

}

func adminConsole() {

	fmt.Println("Please select option:")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	for _, c := range input {

		digit := c - 48

		switch digit {
		case 1:
			fmt.Print("Create Passenger Account")
		case 2:
			fmt.Print("Create Driver Account")
		}
	}
}
