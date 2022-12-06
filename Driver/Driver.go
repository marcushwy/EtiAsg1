package main

// Import packages
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

// Structs and Json
type Driver struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	MobileNumber string `json:"MobileNumber"`
	Email        string `json:"Email"`
	CarLicense   string `json:"CarLicense"`
	IcNo         string `json:"IcNo"`
	DriverStatus string `json:"DriverStatus"`
}

type Trip struct {
	PassengerId string `json:"PassengerId"`
	PickupCode  string `json:"PickupCode"`
	DropoffCode string `json:"DropoffCode"`
	DriverId    string `json:"DriverId"`
	TripStatus  string `json:"TripStatus"`
	TripDate    string `json:"TripDate"`
}

type AllDriver struct {
	Driver map[string]Driver `json:"Driver"`
}

type AllTrip struct {
	Trip map[string]Trip `json:"Trip"`
}

// Empty Maps
var driverlist = map[string]Driver{}
var drivertriplist = map[string]Trip{}
var newdriver = map[string]Driver{}

// Main
func main() {

	// Connect to Router
	router := mux.NewRouter()

	// Driver Account Related
	router.HandleFunc("/api/v1/driver/view/", viewdriver).Methods("GET")          // view driver
	router.HandleFunc("/api/v1/driver/{driverid}", driver).Methods("POST", "PUT") // create and update driver

	// Driver Status and Trip related
	router.HandleFunc("/api/v1/driver/status/online/{driverid}", driveronline).Methods("GET", "PUT")
	router.HandleFunc("/api/v1/driver/status/offline/{driverid}", driveroffline).Methods("GET", "PUT")
	router.HandleFunc("/api/v1/drivers/", autoassigndriver) // auto assign driver

	router.HandleFunc("/api/v1/driver/trips/{driverid}", getdrivertrip).Methods("GET") // get trips assigned to driver

	router.HandleFunc("/api/v1/driver/trips/start/{tripid}", begintrip).Methods("PUT")           //start trip
	router.HandleFunc("/api/v1/driver/trips/end/{tripid}/{driverid}", finishtrip).Methods("PUT") //end trip

	//Port 3000
	fmt.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// Retrieve Driver
func viewdriver(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("failed to connect to db")
	}
	defer db.Close()

	results, err := db.Query("select * from Driver")
	if err != nil {
		fmt.Println("failed to connect ")
	}

	for results.Next() {
		var d Driver
		var driverid string
		err := results.Scan(&driverid, &d.FirstName, &d.LastName, &d.MobileNumber, &d.Email, &d.CarLicense, &d.IcNo, &d.DriverStatus)
		if err != nil {
			fmt.Println("failed to scan")
		}

		fmt.Println("Driver ID no.: ", driverid, "Driver Name: ", d.FirstName, d.LastName, "Mobile No.:", d.MobileNumber, "Email:", d.Email, "Car License:", d.CarLicense, "IcNo.:", d.IcNo)
		driverlist[driverid] = d
	}

	data, _ := json.Marshal(map[string]map[string]Driver{"Drivers": driverlist})
	fmt.Fprintf(w, "%s\n", data)
}

// Create, Update Driver
func driver(w http.ResponseWriter, r *http.Request) {

	//Retrieve the variable params variable
	params := mux.Vars(r)

	//POST
	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {

			var data Driver //input from the user

			if err := json.Unmarshal(body, &data); err == nil {
				_, pExist := driverExist(params["driverid"])
				if !pExist {

					fmt.Println(data)
					fmt.Println(params)
					fmt.Println("hi im at 202, driver does not exists ")
					insertDriver(params["driverid"], data) //insert to db //passing values into function
					//passengerlist[params["passengerid"]] = data  //insert to map

					w.WriteHeader(http.StatusAccepted) //202

				} else {
					fmt.Println("hi im at 409, driver id exists ")
					w.WriteHeader(http.StatusConflict) //409
					fmt.Fprintf(w, "Driver ID exists")
				}

			} else {
				fmt.Println(err)
			}
		}

		//PUT
	} else if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Driver //input from the user

			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				_, pExist := driverExist(params["driverid"])
				if pExist {

					fmt.Println("Driver Exists, updating")
					w.WriteHeader(http.StatusAccepted)
					updateDriver(params["driver"], data)
					//passengerlist[params["passengerid"]] = data

				} else {

					fmt.Println(data)
					fmt.Println(params)
					w.WriteHeader(http.StatusNotFound)
					fmt.Println("Driver ID does not exist")

				}
			} else {
				w.WriteHeader(http.StatusConflict)
				fmt.Println("Passenger ID does not exist")

			}
		}
	}

}

// Update Driver Status
func driveronline(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Driver

			fmt.Println(data.DriverStatus)

			if err := json.Unmarshal(body, &data); err == nil {

				_, dExist := driverExist(params["driverid"])
				_, dOngoing := driverongoing(params["driverid"])

				if dExist {

					if !dOngoing {

						fmt.Println("Driver Exists, Updating Status to Available")
						w.WriteHeader(http.StatusAccepted) //202
						updateAvailable(params["driverid"])

					} else {
						w.WriteHeader(http.StatusConflict) //409
						fmt.Println("Driver has ongoing ride! ")
					}

				} else {

					w.WriteHeader(http.StatusNotFound) //404
					fmt.Println("Driver ID does not exist")

				}

			}

		}
	}
}

func driveroffline(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Driver

			fmt.Println(data.DriverStatus)

			if err := json.Unmarshal(body, &data); err == nil {

				_, dExist := driverExist(params["driverid"])

				if dExist {

					fmt.Println("Driver Exists, Updating Status to Available")
					w.WriteHeader(http.StatusAccepted) //202
					updateBusy(params["driverid"], data)

				} else {

					w.WriteHeader(http.StatusNotFound) //404
					fmt.Println("Driver ID does not exist")

				}

			}

		}
	}
}

// Start Trip
func begintrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Trip

			if err := json.Unmarshal(body, &data); err == nil {

				starttrip(params["tripid"])

				w.WriteHeader(http.StatusAccepted) //202

			} else {
				w.WriteHeader(http.StatusNotFound) //404
			}

		}
	}
}

// End Trip
func finishtrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Trip

			if err := json.Unmarshal(body, &data); err == nil {

				endtrip(params["tripid"])
				updateAvailable(params["driverid"])

				w.WriteHeader(http.StatusAccepted) //202

			} else {
				w.WriteHeader(http.StatusNotFound) //404
			}

		}
	}
}

// Get a Random Driver with Status Available
func autoassigndriver(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("failed to connect to db")
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM DRIVER WHere DriverStatus = 'Available' ORDER BY RAND() LIMIT 1")
	if err != nil {
		fmt.Println("failed to connect to insert")
		fmt.Println(err)
	} else {
		fmt.Println("Connected succesfully")
	}
	for results.Next() {

		var d Driver
		var driverid string

		err := results.Scan(&driverid, &d.FirstName, &d.LastName, &d.MobileNumber, &d.Email, &d.CarLicense, &d.IcNo, &d.DriverStatus)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(&driverid)
		}
		newdriver[driverid] = d
	}

	data, _ := json.Marshal(map[string]map[string]Driver{"Selected Drivers": newdriver})
	fmt.Fprintf(w, "%s\n", data)

}

// Get Trips assigned to driver
func getdrivertrip(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(err)
	}
	defer db.Close()

	results, err := db.Query("Select * from Trip where DriverId = ? AND TripStatus !='Completed'", params["driverid"])
	//db exec will return rows and primary key

	if err != nil {
		fmt.Println(err)
	}
	drivertriplist = map[string]Trip{}

	for results.Next() {
		var t Trip
		var tripid string

		err := results.Scan(&tripid, &t.PassengerId, &t.PickupCode, &t.DropoffCode, &t.DriverId, &t.TripStatus, &t.TripDate)
		if err != nil {
			fmt.Println("failed to scan")
		}
		fmt.Println("ID no.: ", tripid, "Passenger ID: ", t.PassengerId, "PickupCode:", t.PickupCode, "DropoffCode:", t.DropoffCode, "Driver Id:", t.DriverId, "Trip Status:", t.TripStatus, "Trip Date:", t.TripDate)

		drivertriplist[tripid] = t
	}

	data, _ := json.Marshal(map[string]map[string]Trip{"Driver's Trips": drivertriplist})
	fmt.Fprintf(w, "%s\n", data)

}

// MySQL Database

// - create new driver
func insertDriver(driverid string, d Driver) {

	var status string = "Busy"

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("failed to connect to db")
	}
	defer db.Close()

	_, err = db.Exec("INSERT into Driver values(CAST(? as double),?,?,?,?,?,?,?)", driverid, d.FirstName, d.LastName, d.MobileNumber, d.Email, d.CarLicense, d.IcNo, status)
	if err != nil {
		fmt.Println("failed to connect to insert")
		fmt.Println(err)
	} else {
		fmt.Println("New Driver Added")
	}

}

// - check if driver exist
func driverExist(driverid string) (Driver, bool) {
	//Connect to Database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var d Driver
	var id string

	results := db.QueryRow("Select * from Driver where id =?", driverid)
	err = results.Scan(&id, &d.FirstName, &d.LastName, &d.MobileNumber, &d.Email, &d.CarLicense, &d.IcNo, &d.DriverStatus)
	fmt.Println(err)

	if err == sql.ErrNoRows {
		return d, false
	}

	return d, true
}

// - update driver information
func updateDriver(driverid string, d Driver) {
	//Connect to database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Driver SET FirstName=?, LastName=?, MobileNumber=?, Email=?, CarLicense=? WHERE ID = ?", d.FirstName, d.LastName, d.MobileNumber, d.Email, d.CarLicense, driverid)
	//db exec will return rows and primary key

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Driver Updated Successfully")
	}
}

// - update driver status to busy
func updateBusy(driverid string, d Driver) {
	//Connect to database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Driver SET DriverStatus=? WHERE ID=?", d.DriverStatus, driverid)
	//db exec will return rows and primary key

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Driver Status updated Successfully")
	}
}

// - update driver status to available
func updateAvailable(driverid string) {
	//Connect to database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(err)
	}
	defer db.Close()

	var driverstatus = "Available"

	_, err = db.Exec("UPDATE Driver SET DriverStatus=? WHERE ID=?", driverstatus, driverid)
	//db exec will return rows and primary key

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Driver Status updated Successfully")
	}
}

// - Start Trip
func starttrip(tripid string) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(err)
	}
	defer db.Close()

	var tripstat string = "Ongoing"

	_, err = db.Exec("UPDATE Trip SET TripStatus=? WHERE TripId=? AND TripStatus = 'Accepted'", tripstat, tripid)
	//db exec will return rows and primary key

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Trip is currently ongoing")
	}
}

// - End Trip
func endtrip(tripid string) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(err)
	}
	defer db.Close()

	var tripstat string = "Completed"

	_, err = db.Exec("UPDATE Trip SET TripStatus=? WHERE TripId=? AND TripStatus != 'Completed' ", tripstat, tripid)
	//db exec will return rows and primary key

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Trip is currently ongoing")
	}
}

// - Driver in ongoing trip, not allowed to change to available
func driverongoing(driverid string) (Trip, bool) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var t Trip

	results := db.QueryRow("Select * from trip where driverid =? AND Trip != 'Completed'", driverid)
	err = results.Scan(&driverid, &t.TripStatus)
	fmt.Println(err)

	if err == sql.ErrNoRows {
		return t, false
	}

	return t, true
}
