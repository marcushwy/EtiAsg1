package main

//Import
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Passenger Struct
type Passenger struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	MobileNumber string `json:"MobileNumber"`
	Email        string `json:"Email"`
}

// Driver Struct
type Driver struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	MobileNumber string `json:"MobileNumber"`
	Email        string `json:"Email"`
	CarLicense   string `json:"CarLicense"`
	IcNo         string `json:"IcNo"`
	DriverStatus string `json:"DriverStatus"`
}

// Trip Struct
type Trip struct {
	PassengerId string `json:"PassengerId"`
	PickupCode  string `json:"PickupCode"`
	DropoffCode string `json:"DropoffCode"`
	DriverId    string `json:"DriverId"`
	TripStatus  string `json:"TripStatus"`
	TripDate    string `json:"TripDate"`
}

// Maps
type Passengers struct {
	Passengers map[string]Passenger `json:"Passengers"`
}

type Drivers struct {
	Drivers map[string]Driver `json:"Drivers"`
}

type Trips struct {
	Trips map[string]Trip `json:"Trips"`
}

// Variable to store random driver with status set as available (to auto assign him when a trip is created.
var randomdriver Driver
var randomdriverid string

func main() {
outer:
	for {
		fmt.Println("\n")

		fmt.Println(" ==== Welcome to Ride Share by Marcus ==== \n",
			"\n",
			"<----Passenger Menu---->\n",
			"1.  Create New Passenger \n",       //create new passsenger
			"2.  Update Existing Passenger \n",  //update existing passenger
			"3.  View All Passenger - Admin \n", //view all passenger
			"4.  Request for New Trip \n",       //create new trip
			"5.  View Trip History \n",          //view passenger trip history
			"\n",
			"<----Driver Menu---->\n",
			"6.  Create New Driver\n",             //create new driver, driver status default set to busy
			"7.  Update Existing Driver\n",        //update all existing driver information, except for IcNo.
			"8.  View All Driver - Admin \n",      //view all driver
			"9.  Go Online  (Available Mode)\n",   //update status to available (eligible for trip assignment)
			"10. Go Offline (Busy Mode)\n",        //update status to buys (not eligible for trip assignment)
			"11. View Driver's Assigned Trips \n", //view trips assigned to driver, then able to start end or exit
			"0.  Exit")

		fmt.Print("Please select an option: ")

		var option int
		fmt.Scanf("%d\n", &option)

		switch option {

		case 0:
			break outer
		//Passenger
		case 1:
			createPassenger()
		case 2:
			updatePassenger()
		case 3:
			viewPassengers()
		case 4:
			createTrip()
		case 5:
			viewPassengerTrips()

		//Driver
		case 6:
			createDriver()
		case 7:
			updateDriver()
		case 8:
			viewDrivers()
		case 9:
			updateAvailable()
		case 10:
			updateBusy()
		case 11:
			viewDriverTrips()

		default:
			fmt.Println("Sorry, we didn't catch that.")

		}
	}
}

//PASSENGER MENU FUNCTIONS

// Function - Create New Passenger
func createPassenger() {

	reader := bufio.NewReader(os.Stdin)

	var newpassenger Passenger
	fmt.Println("\n")
	fmt.Println("=== New Passenger Creation ===")

	var randpid string
	var randit int
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randit = r.Intn(5000)
	randpid = "P" + strconv.Itoa(randit)

	fmt.Print("Enter First Name: ")
	inputf, _ := reader.ReadString('\n')
	newpassenger.FirstName = strings.TrimSpace(inputf)

	fmt.Print("Enter Last Name: ")
	inputl, _ := reader.ReadString('\n')
	newpassenger.LastName = strings.TrimSpace(inputl)

	fmt.Print("Enter Mobile Number: ")
	inputm, _ := reader.ReadString('\n')
	newpassenger.MobileNumber = strings.TrimSpace(inputm)

	fmt.Print("Enter Email: ")
	inpute, _ := reader.ReadString('\n')
	newpassenger.Email = strings.TrimSpace(inpute)

	postBody, _ := json.Marshal(newpassenger)
	resBody := bytes.NewBuffer(postBody)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/passenger/"+randpid, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("* New Passenger with ID : ", randpid, " created! *")
			} else if res.StatusCode == 409 {
				fmt.Println("* Error - Passenger", randpid, "already exists! *")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

// Function - Update Existing Passenger
func updatePassenger() {

	reader := bufio.NewReader(os.Stdin)

	var updatePassenger Passenger

	fmt.Println("\n")
	fmt.Println("=== Please Enter New Passenger Information ===")
	fmt.Print("Enter Existing Passenger ID :")
	var passengerid string
	fmt.Scanf("%v\n", &passengerid)

	fmt.Print("Enter First Name: ")
	inputf, _ := reader.ReadString('\n')
	updatePassenger.FirstName = strings.TrimSpace(inputf)

	fmt.Print("Enter Last Name: ")
	inputl, _ := reader.ReadString('\n')
	updatePassenger.LastName = strings.TrimSpace(inputl)

	fmt.Print("Enter Mobile Number: ")
	inputm, _ := reader.ReadString('\n')
	updatePassenger.MobileNumber = strings.TrimSpace(inputm)

	fmt.Print("Enter Email: ")
	inpute, _ := reader.ReadString('\n')
	updatePassenger.Email = strings.TrimSpace(inpute)

	postBody, _ := json.Marshal(updatePassenger)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/passenger/"+passengerid, bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("* Passenger ", passengerid, " updated successfully! *")
			} else if res.StatusCode == 404 {
				fmt.Println("* Error - Passenger", passengerid, "does not exist! *")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

// Function - View All Passenger (Admin)
func viewPassengers() {

	fmt.Println("\n")
	fmt.Println("=== View All Passenger ===")

	//connect to client
	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5000/api/v1/passenger/view/", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {

				var res Passengers

				json.Unmarshal(body, &res)
				fmt.Println("=== Passenger Info ===")
				for k, v := range res.Passengers {
					fmt.Println("\n")
					fmt.Println("Passenger ID : ", k, " ")
					fmt.Println("First Name : ", v.FirstName)
					fmt.Println("Last Name : ", v.LastName)
					fmt.Println("MobileNumber : ", v.MobileNumber)
					fmt.Println("Email : ", v.Email)
					fmt.Println("\n")

				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

}

// Function - Create Trip Request, get single random Available driver, create trip record, update driver status to Busy
func createTrip() {

	resp, err := http.Get("http://localhost:3000/api/v1/drivers/")
	if err != nil {
		fmt.Println("All Drivers are busy")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var assigndriver map[string]map[string]Driver
	json.Unmarshal([]byte(body), &assigndriver)

	for key, element := range assigndriver["Selected Drivers"] {
		randomdriverid = key
		randomdriver = element
	}

	reader := bufio.NewReader(os.Stdin)

	if len(randomdriverid) != 0 {
		var newtrip Trip

		var randtripid string
		var randit int
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		randit = r.Intn(10000)
		randtripid = "T" + strconv.Itoa(randit)

		fmt.Print("Enter Passenger ID : ")
		var passengerid string
		fmt.Scanf("%v\n", &passengerid)

		fmt.Print("Enter Pickup Code: ")
		inputf, _ := reader.ReadString('\n')
		newtrip.PickupCode = strings.TrimSpace(inputf)

		fmt.Print("Enter Dropoff Code: ")
		inputl, _ := reader.ReadString('\n')
		newtrip.DropoffCode = strings.TrimSpace(inputl)

		newtrip.TripStatus = "Accepted"

		newtrip.DriverId = randomdriverid

		postBody, _ := json.Marshal(newtrip)
		resBody := bytes.NewBuffer(postBody)

		client := &http.Client{}
		if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/passenger/trip/"+randtripid+"/"+passengerid+"/"+randomdriverid, resBody); err == nil {
			if res, err := client.Do(req); err == nil {
				if res.StatusCode == 202 {
					fmt.Println("* New Trip:", randtripid, "created successfully! *")

					var updateBusy Driver

					updateBusy.DriverStatus = "Busy"

					driverpostBody, _ := json.Marshal(updateBusy)

					if req, err := http.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/driver/status/offline/"+randomdriverid, bytes.NewBuffer(driverpostBody)); err == nil {

						if res, err := client.Do(req); err == nil {
							if res.StatusCode == 202 {
								fmt.Println("* A new driver is assigned to your trip! * ")
							} else if res.StatusCode == 404 {
								fmt.Println("* Driver does not exist! * ")
							}
						} else {
							fmt.Println(2, err)
						}
					} else {
						fmt.Println(3, err)
					}
				} else if res.StatusCode == 409 {
					fmt.Println("*** Error - Passenger", passengerid, "already in Ongoing Trip! *** ")
				} else if res.StatusCode == 404 {
					fmt.Println("*** Error - Passenger", passengerid, "does not exist! *** ")
				}
			} else {
				fmt.Println(2, err)
			}
		} else {
			fmt.Println(3, err)
		}

	} else {

		fmt.Print("* No Drivers Available at the moment! *")
	}
}

// Function - View Passenger Trip History in reverse chronological order
func viewPassengerTrips() {

	var passengerid string
	fmt.Print("Please enter Passenger ID to view trips: ")
	fmt.Scanf("%v\n", &passengerid)

	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5000/api/v1/trip/"+passengerid, nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {

				var res map[string]map[string]Trip

				json.Unmarshal(body, &res)

				if len(res["Trip"]) == 0 {
					fmt.Println("* No Trips Available! *")
				} else {

					fmt.Println("=== Trip Info ===")
					for k, v := range res["Trip"] {

						fmt.Println("Trip ID : ", k, " ")
						fmt.Println("Pickup Code : ", v.PickupCode)
						fmt.Println("Dropoff Code : ", v.DropoffCode)
						fmt.Println("Trip Status : ", v.TripStatus)
						fmt.Println("Trip Date : ", v.TripDate, "\n ")

					}
				}
			}
		}
	}
}

//DRIVER MENU FUNCTIONS

// Function - Create New Driver
func createDriver() {

	reader := bufio.NewReader(os.Stdin)

	var newdriver Driver

	fmt.Println("=== Please Enter New Driver Information ===")
	var randdid string
	var randit int
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randit = r.Intn(5000)
	randdid = "D" + strconv.Itoa(randit)

	fmt.Print("Enter First Name: ")
	inputf, _ := reader.ReadString('\n')
	newdriver.FirstName = strings.TrimSpace(inputf)

	fmt.Print("Enter Last Name: ")
	inputl, _ := reader.ReadString('\n')
	newdriver.LastName = strings.TrimSpace(inputl)

	fmt.Print("Enter Mobile Number: ")
	inputm, _ := reader.ReadString('\n')
	newdriver.MobileNumber = strings.TrimSpace(inputm)

	fmt.Print("Enter Email: ")
	inpute, _ := reader.ReadString('\n')
	newdriver.Email = strings.TrimSpace(inpute)

	fmt.Print("Enter Car License: ")
	inputc, _ := reader.ReadString('\n')
	newdriver.CarLicense = strings.TrimSpace(inputc)

	fmt.Print("Enter IC: ")
	inputic, _ := reader.ReadString('\n')
	newdriver.IcNo = strings.TrimSpace(inputic)

	newdriver.DriverStatus = "Busy"

	postBody, _ := json.Marshal(newdriver)
	resBody := bytes.NewBuffer(postBody)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/driver/"+randdid, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("* Driver with ID: ", randdid, "created! *")
			} else if res.StatusCode == 409 {
				fmt.Println("* Error - Driver", randdid, "already exists! *")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

// Function - Update Existing Driver Info
func updateDriver() {
	reader := bufio.NewReader(os.Stdin)

	var updateDriver Driver

	fmt.Println("\n")
	fmt.Println("=== Please Enter New Driver Information ===")
	fmt.Print("Enter Existing Driver ID :")
	var driverid string
	fmt.Scanf("%v\n", &driverid)

	fmt.Print("Enter First Name: ")
	inputf, _ := reader.ReadString('\n')
	updateDriver.FirstName = strings.TrimSpace(inputf)

	fmt.Print("Enter Last Name: ")
	inputl, _ := reader.ReadString('\n')
	updateDriver.LastName = strings.TrimSpace(inputl)

	fmt.Print("Enter Mobile Number: ")
	inputm, _ := reader.ReadString('\n')
	updateDriver.MobileNumber = strings.TrimSpace(inputm)

	fmt.Print("Enter Email: ")
	inpute, _ := reader.ReadString('\n')
	updateDriver.Email = strings.TrimSpace(inpute)

	fmt.Print("Enter Car License: ")
	inputc, _ := reader.ReadString('\n')
	updateDriver.CarLicense = strings.TrimSpace(inputc)

	postBody, _ := json.Marshal(updateDriver)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/driver/"+driverid, bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("* Driver", driverid, "updated successfully! *")
			} else if res.StatusCode == 404 {
				fmt.Println("* Error - Driver", driverid, "does not exist! *")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

// Function - View All Drivers (Admin)
func viewDrivers() {

	fmt.Println("\n")
	fmt.Println("=== View All Drivers ===")
	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/driver/view/", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {

				var res Drivers

				json.Unmarshal(body, &res)
				fmt.Println("=== Driver Info ===")
				for k, v := range res.Drivers {
					fmt.Println("\n")

					fmt.Println("Driver ID : ", k, " ")
					fmt.Println("First Name : ", v.FirstName)
					fmt.Println("Last Name : ", v.LastName)
					fmt.Println("MobileNumber : ", v.MobileNumber)
					fmt.Println("Email : ", v.Email)
					fmt.Println("Car License : ", v.CarLicense)
					fmt.Println("IcNo : ", v.IcNo)
					fmt.Println("Status : ", v.DriverStatus)
					fmt.Println("\n")
				}
			}
		}
	}
}

// Function - Update Status to Busy, prompt for driverid
func updateBusy() {

	var updateBusy Driver

	var driverid string
	fmt.Print("Enter Driver ID : ")
	fmt.Scanf("%v\n", &driverid)

	fmt.Print("Changing Status...")
	updateBusy.DriverStatus = "Busy"

	postBody, _ := json.Marshal(updateBusy)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/driver/status/offline/"+driverid, bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("* Driver", driverid, " went offline! *")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Driver", driverid, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

// Function - Update Status to Available, prompt for driverid, only can update when driver does not have ongoing ride
func updateAvailable() {

	var updateAvailable Driver

	var driverid string

	fmt.Print("Enter Driver ID : ")
	fmt.Scanf("%v\n", &driverid)

	fmt.Print("Changing Status...")
	updateAvailable.DriverStatus = "Available"

	postBody, _ := json.Marshal(updateAvailable)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/driver/status/online/"+driverid, bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("* Driver", driverid, "went online! *")
			} else if res.StatusCode == 409 {
				fmt.Println("Error - Driver", driverid, "has an ongoing ride!")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Driver", driverid, "does not exist!")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

// Function - View Driver's Trips Assigned + Trip Management, get trips assigned to the input driverid
func viewDriverTrips() {

	var driverid string
	fmt.Print("Please enter Driver ID to view trips: ")
	fmt.Scanf("%v\n", &driverid)

	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/api/v1/driver/trips/"+driverid, nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {

				var res map[string]map[string]Trip

				json.Unmarshal(body, &res)

				fmt.Println("=== Trip Info ===")
				for k, v := range res["Driver's Trips"] {

					fmt.Println("Trip ID : ", k, " ")
					fmt.Println("Pickup Code : ", v.PickupCode)
					fmt.Println("Dropoff Code : ", v.DropoffCode)
					fmt.Println("Trip Status : ", v.TripStatus)
					fmt.Println("Trip Date : ", v.TripDate, "\n ")

				}

				if len(res["Driver's Trips"]) != 0 {

					var startend int
					fmt.Println("Do you want to :\n",
						"1. Start a Trip\n",
						"2. End a Trip\n",
						"3. Back to Menu")
					fmt.Print("Please select an option :")
					fmt.Scanf("%v\n", &startend)
					switch startend {

					case 1:
						initiateTrip()
					case 2:
						endTrip()
					case 3:
						break
					}
				} else {
					fmt.Println("* No Trips Assigned! *")
				}

			}
		}
	}

}

// Function - Start a Trip, prompt for driver and trip id, change trip status to ongoing
func initiateTrip() {

	var updateTrip Trip
	var driverid string
	var tripid string

	fmt.Print("Enter Driver ID : ")
	fmt.Scanf("%v\n", &driverid)
	fmt.Print("Enter Trip ID : ")
	fmt.Scanf("%v\n", &tripid)

	fmt.Print("Changing Status...")

	updateTrip.TripStatus = "Ongoing"

	postBody, _ := json.Marshal(updateTrip)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/driver/trips/start/"+tripid, bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("* Trip", tripid, "has started! *")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Driver", driverid, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}

// Function - End a Trip, prompt for driver and trip id, change trip status to completed
func endTrip() {
	var updateTrip Trip

	var driverid string
	var tripid string

	fmt.Print("Enter Driver ID : ")
	fmt.Scanf("%v\n", &driverid)
	fmt.Print("Enter Trip ID : ")
	fmt.Scanf("%v\n", &tripid)

	fmt.Print("Changing Status...")

	updateTrip.TripStatus = "Completed"

	postBody, _ := json.Marshal(updateTrip)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:3000/api/v1/driver/trips/end/"+tripid+"/"+driverid, bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("* Trip", tripid, "has ended! *")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Driver", driverid, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

}
