package main

// Import packages
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

// Structs and Json
type Passenger struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	MobileNumber string `json:"MobileNumber"`
	Email        string `json:"Email"`
}

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

type AllPassenger struct {
	Passengers map[string]Passenger `json:"Passenger"`
}

type AllDriver struct {
	Driver map[string]Driver `json:"Driver"`
}

type AllTrip struct {
	Trip map[string]Trip `json:"Trip"`
}

// Empty Map
var passengerlist = map[string]Passenger{}
var triplist = map[string]Trip{}

// Main
func main() {

	// Connect to Router
	router := mux.NewRouter()

	// Passenger Account Related
	router.HandleFunc("/api/v1/passenger/view/", viewpassenger).Methods("GET")             // view specific passenger
	router.HandleFunc("/api/v1/passenger/{passengerid}", passenger).Methods("POST", "PUT") // create and update passenger

	// Passenger Trip Related
	router.HandleFunc("/api/v1/passenger/trip/{tripid}/{passengerid}/{driverid}", newtrip).Methods("GET", "POST") // create new trip
	router.HandleFunc("/api/v1/trip/{passengerid}", viewpassengertrip).Methods("GET")                             // get passenger trip history

	//Port 5000
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}

// View Passenger
func viewpassenger(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	results, err := db.Query("select * from Passenger")
	if err != nil {
		fmt.Println(err)
	}

	for results.Next() {

		var p Passenger
		var passengerid string

		err := results.Scan(&passengerid, &p.FirstName, &p.LastName, &p.MobileNumber, &p.Email)
		if err != nil {
			fmt.Println("failed to scan")
		}
		fmt.Println("ID no.: ", passengerid, "Passenger Name: ", p.FirstName, p.LastName, "Mobile No.:", p.MobileNumber, "Email:", p.Email)

		passengerlist[passengerid] = p
	}
	data, _ := json.Marshal(map[string]map[string]Passenger{"Passengers": passengerlist})
	fmt.Fprintf(w, "%s\n", data)
}

// Create New Passenger, Update Existing Passenger
func passenger(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	//POST
	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Passenger //input from the user

			if err := json.Unmarshal(body, &data); err == nil {

				//If passenger does not exist, insert into database.
				_, pExist := passengerExist(params["passengerid"])

				if !pExist {

					fmt.Println("Passenger ID does not exist, Creating new Passenger...")
					insertPassenger(params["passengerid"], data) //insert to db //passing values into function
					w.WriteHeader(http.StatusAccepted)           //202

				} else {

					fmt.Println("Passenger ID exists, Please enter a new ID...")
					w.WriteHeader(http.StatusConflict) //409
					fmt.Fprintf(w, "Passenger ID exists")
				}

			} else {
				fmt.Println(err)
			}
		}

		//PUT
	} else if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Passenger //input from the user

			if err := json.Unmarshal(body, &data); err == nil {

				_, pExist := passengerExist(params["passengerid"])

				if pExist {

					fmt.Println("Passenger ID Exists, Updating Passenger Details...")
					updatePassenger(params["passengerid"], data)
					//passengerlist[params["passengerid"]] = data
					w.WriteHeader(http.StatusAccepted) //202

				} else {

					w.WriteHeader(http.StatusNotFound) //404
					fmt.Println("Passenger ID does not exist")

				}
			} else {
				w.WriteHeader(http.StatusConflict) //409
				fmt.Println("Passenger ID does not exist")
			}

		} else {
			fmt.Println(err)
		}

	}

}

// Create New Trip
func newtrip(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Trip
			if err := json.Unmarshal(body, &data); err == nil {

				_, pExist := passengerExist(params["passengerid"])
				_, pExistOngoing := passengerTripOngoing(params["passengerid"])

				if pExist {

					if !pExistOngoing {

						createTrip(params["tripid"], params["passengerid"], params["driverid"], data)

					} else {
						fmt.Println("Passenger is already in an Ongoing Trip!")
						w.WriteHeader(http.StatusConflict) //409
					}

				} else {
					fmt.Println("Please create a Passenger First!")
					w.WriteHeader(http.StatusNotFound) //404

				}
			}
		}
	}
}

// Get Trip History
func viewpassengertrip(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("failed to connect to db")
	}
	defer db.Close()

	results, err := db.Query("select * from Trip where PassengerId = ? ORDER BY TripDate Desc", params["passengerid"])
	if err != nil {
		fmt.Println("failed to connect ")
	}
	triplist = map[string]Trip{}
	//adding this map keeps it empty

	for results.Next() {
		var t Trip
		var passengerid string
		var tripid string
		var driverid string

		err := results.Scan(&tripid, &passengerid, &t.PickupCode, &t.DropoffCode, &driverid, &t.TripStatus, &t.TripDate)
		if err != nil {
			fmt.Println("failed to scan")
			fmt.Println(err)
		}

		triplist[tripid] = t
		//empty map created ontop ; tripid key ; Trip struct

	}
	data, _ := json.Marshal(map[string]map[string]Trip{"Trip": triplist})

	fmt.Fprintf(w, "%s\n", data)

}

//MySQL Database

// - check if passenger exist
func passengerExist(passengerid string) (Passenger, bool) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var p Passenger
	var id string

	results := db.QueryRow("Select * from Passenger where id =?", passengerid)
	err = results.Scan(&id, &p.FirstName, &p.LastName, &p.MobileNumber, &p.Email)
	fmt.Println(err)

	if err == sql.ErrNoRows {
		return p, false
	}
	return p, true
}

// - check if passenger exist in trip
func passengerExistinTrip(passengerid string) (Trip, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var t Trip
	var tripid string

	results := db.QueryRow("Select * from Trip where passengerid =?", passengerid)
	err = results.Scan(&tripid, &passengerid, &t.PickupCode, &t.DropoffCode, &t.DriverId, &t.TripStatus, &t.TripDate)
	fmt.Println(err)

	if err == sql.ErrNoRows {
		return t, false
	}
	return t, true
}

// - check if passenger is in an ongoing trip
func passengerTripOngoing(passengerid string) (Trip, bool) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var t Trip
	var tripid string

	results := db.QueryRow("Select * from Trip where passengerid =? AND TripStatus != 'Completed' ", passengerid)
	err = results.Scan(&tripid, &passengerid, &t.PickupCode, &t.DropoffCode, &t.DriverId, &t.TripStatus, &t.TripDate)
	fmt.Println(err)

	if err == sql.ErrNoRows {
		return t, false
	}
	return t, true
}

// - create new passenger
func insertPassenger(passengerid string, p Passenger) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("Error")
	}
	defer db.Close()

	_, err = db.Exec("INSERT into Passenger values(CAST(? as double),?,?,?,?)", passengerid, p.FirstName, p.LastName, p.MobileNumber, p.Email)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
	} else {
		fmt.Println("New Passenger Inserted.")
	}
}

// - update existing passenger
func updatePassenger(passengerid string, p Passenger) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Passenger SET FirstName=?, LastName=?, MobileNumber=?, Email=? WHERE ID = ?", p.FirstName, p.LastName, p.MobileNumber, p.Email, passengerid)
	//db exec will return rows and primary key

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Passenger Updated Successfully")
	}

}

// - create new trip request
func createTrip(tripid string, passengerid string, driverid string, t Trip) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		fmt.Println("failed to connect to db")
	}
	defer db.Close()

	var triptime = time.Now()

	_, err = db.Exec("INSERT into Trip (TripId, PassengerId, PickupCode, DropoffCode,DriverId, TripStatus, TripDate) VALUES (?,?,?,?,?,?,?)", tripid, passengerid, t.PickupCode, t.DropoffCode, driverid, t.TripStatus, triptime)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("New Trip Added")
	}
}
