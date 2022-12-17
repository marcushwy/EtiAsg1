
# ETI ASSIGNMENT 1 - Ride Sharing Platform

Name: Marcus Hon Wei Yang 
ID: S10207469

# Introduction

Hello and Welcome to the Readme for my console application named Ride-Share, developed for my first ETI assignment in Ngee Ann Polytechnic.
In this Readme, my thought processes and thought considerations to implement this assignment will be documented.


# Microservices and Domain-Driven Design

Before working on this application, I have conducted extensive research through referencing online material to write a report on the concepts of Domain-Driven-Design while using the Shopee application as an example. The research process played a big part to determine the structure of the current console application as I was exposed to new concepts when it comes to programming with DDD. For example:  Model-Driven-Design, Strategic Design and Tactical Design. 


# Strategic Design (Domains, Sub-Domains and Contexts)

The Strategic Design phase of DDD helped me to plan out the design of the microservice architecture from a high-level perspective and understand the crucial features that each microservice would need to contain without going too in-depth on “how” the application would be implemented (such as technology). This phase of design is important because it prevents me from straying away from the important business capabilities of the software. 
Ride-Share is a console application that has somewhat similar functionalities when compared to the famous ride-grabbing app Uber and Grab. Therefore, a defined domain is important in the planning phase of this assignment because it helps to identify the actual business problems that I am targeting to solve / assist with Ride-Share. 

Business Domain: Ride-Sharing Domain
Sub Domains: Passenger, Driver, Trip

For this case, the Passenger, Driver, and Trip sub domains would be considered as “Core Domains” as they are domains that hold the upmost importance when it comes to the main purpose of a typical ride-sharing application, which is to connect passengers and drivers together to form a ride-sharing ecosystem. Therefore, these subdomains should be our key focus on time and resource because they come together to enable the ride-sharing functionality.

Then, the application's business functionalities is recorded down, in terms of different contexts. 

In the Passenger and Driver Context:
-	Account Management (Create Account, Update Account, View Account) 
// GET, POST, PUT 

In the Trip Booking Context (which links Passenger and Driver together) 
-	Trip Management (Request Trip / Start, End Trip, Review Trip History) 
-	Retrieval of Trip History
//GET, POST, PUT 

To put simply, for the Passenger and Driver subdomains would have bounded contexts that can aid me in determining the appropriate features needed for each potential microservice. 



# Tactical Design

This design phases allows me to think about the application from a more technical perspective, which are the technologies / ways that I could potentially utilize to programme this console application. (With consideration of assignment requirements) 

Programming language used: Go Lang 
Database: MySQL Database 

The research process on DDD and the Software Engineering module helped me to understand that microservices should be developed in a way such that it is loosely coupled. Which means that each of the microservices should have lower dependencies on each other. Not only that, I understood that each bounded context (e.g., Passenger) should contain their own microservice to tackle a specific problem, separate databases and APIs. 

The microservices identified for the Ride-Share console application would be: 
-	M.Svc 1: Passenger 
-	M.Svc 2: Driver 
-	M.Svc 3: Menu / Console Interface 

The Passenger and Driver microservices would contain features that allow the user to create Passenger and Driver Accounts and manage their Trip bookings. Because of this, the word “Trip” in the Passenger and Driver contexts would mean different things because as a Passenger, they would be able to request trips and view their trip history. However, for Driver they would be assigned to trips automatically, view ongoing trips and have the option to initiate and end them. 

With those considerations in mind, I can now move on to create the separate data tables that would hold the information of my microservices.

The tables created through MySQL are: 
-	Passenger Table 
-	Driver Table 
-	Trip Table

Editor’s Note: This Assignment would be implemented using one database. However, to simulate and implement the concepts of DDD (e.g. to have individual databases for each svc.), there would be no use of the INNER JOIN function / Foreign Keys to combine tables for CRUD data purposes. Basically, each table in the database would not have convenient connections / dependencies on each other. 

Entities: 
-	Passenger Table -> Passenger ID 
-	Driver Table -> Driver ID 
-	Trip Table -> Trip ID

Aggregates/Tables broken down:
-	Passenger Table 
-> ID, FirstName, LastName, MobileNumber, Email

-	Driver Table 
-> ID, FirstName, LastName, MobileNumber, Email, CarLicense, IcNo, DriverStatus

-	Trip Table 
-> TripId, PassengerId, PickupCode, DropoffCode, DriverId, TripStatus, TripDate



# List of Features in each microsvc.

Passenger Svc. 
- View Passenger Account (GET)
- Create Passenger Account (POST) 
- Update Passenger Account (PUT)
- Request New Trip (POST) 
- Retrieve Trips in Reverse Chronological Order (GET) 

Driver Svc. 
- View Driver Account (GET)
- Create Driver Account (POST) 
- Update Passenger Account (PUT) 
- Update Driver Status (Busy / Available) (PUT) 
- Start Trip / End Trip (PUT) 
- View Assigned Trips (GET) 

Trip Svc.
- Main Menu of the console application 
- Calls the different APIs
- Requests for user input to handle different actions

# Passenger Service Documentation

 ## List of API and functions called:
 
1. ("/api/v1/passenger/view/", viewpassenger)
 - The viewpassenger function is called to retrieve the list of passengers.
2. ("/api/v1/passenger/{passengerid}", passenger)
 - The passenger function is called, if the request method is POST, the insertPassenger function is called to create a new passenger to the database provided that the passenger does not exist previously. If the request method is PUT, the updatePassenger function is called, allowing the existing passenger to update their information.
3. ("/api/v1/passenger/trip/{tripid}/{passengerid}/{driverid}", newtrip)
 - The newtrip function is called, if the request method is POST, the createTrip function is called to create a new trip taking in trip id, passengerid and driverid as parameters. Provided that the passenger exist in the database and is not already in an ongoing trip. 
4. ("/api/v1/trip/{passengerid}", viewpassengertrip)
 - The viewpassengertrip function is called, which will print out the trip history of a input passengerid in reverse chronological order. 


# Driver Service Documentation 

 ## List of API and functions called:
 
1. router.HandleFunc("/api/v1/driver/view/", viewdriver)
 - The viewdriver function is called, to retrieve the list of drivers.
2. router.HandleFunc("/api/v1/driver/{driverid}", driver)
 - The driver function is called, if the request method is POST, the insertDriver function is called to create a new passenger to the database provided that the driver does not exist previously. If the request method is PUT, the updateDriver function is called, allowing the existing driver to update their information.
3. router.HandleFunc("/api/v1/driver/status/online/{driverid}", driveronline)
 - The driveronline function is called, allowing the driver to update their status to Available provided that they are not in an ongoing trip.
4.	router.HandleFunc("/api/v1/driver/status/offline/{driverid}", driveroffline)
 - The driveroffline function is called, which updates the existing driver's status to Busy, which means they are not able to be assigned to trip requests.
5.	router.HandleFunc("/api/v1/drivers/", autoassigndriver)
 - the autoassigndriver function is called, which scans through the database for a single random driver with status set to Available so that they can be assigned to a new trip request. 
6.	router.HandleFunc("/api/v1/driver/trips/{driverid}", getdrivertrip)
 - the getdrivertrip function is called, and retrieves the assigned trips for the existing driver. 
7.	router.HandleFunc("/api/v1/driver/trips/start/{tripid}", begintrip)
 - the begintrip function is called, which allows the driver to begin an accepted trip. 
8.	router.HandleFunc("/api/v1/driver/trips/end/{tripid}/{driverid}", finishtrip)
 -  - the endtrip function is called, which allows the driver to end an ongoing trip. 


# Menu Service
 
## List of Options in the main menu and how they work: 

### Passenger-Side (Port 5000)

1. Create Passenger
 - The create passenger option allows the user to create a new passenger and will prompt for multiple input fields. A random passengerid is then assigned to the user. After connecting to the client the id is then appended to the "http://localhost:5000/api/v1/passenger/" as a POST request to create a new passenger. 
 - Status 202 is returned if passenger is created successfully. 
 - Status 409 is returned if passenger already exists. 

2. Update Passenger
 - The update passenger option allows the user to update an existing passenger's information. The user will also be prompted with multiple input fields. With a PUT request, "http://localhost:5000/api/v1/passenger/" with the passengerid and postbody to update the information of the passenger. 
 - Status 202 is returned if passenger is updated successfully.
 - Status 404 is returned if passenger does not exist. 

3. View Pasenger (Admin - feature) 
- The viewPassenger function is called to retrieve the list of passengers that are in the database. With a new GET request and API "http://localhost:5000/api/v1/passenger/view/", the list of passengers are then scanned and then displayed. 

4. Create Trip 
- This option calls the createTrip function in the Menu service and allows the passenger to request for a new trip. First, the Driver API with GET request "http://localhost:3000/api/v1/drivers/" is called to retrieved the drivers that is available to be assigned to a trip. Then the function prompts the user for input of postal code and dropoffcode if a driver is found. Then the user is prompted for pickupcode and dropoff code and a new random tripid is generated. Then with a new POST request and post body, the API "http://localhost:5000/api/v1/passenger/trip/" to create a new trip in the Trip table. 
- Status 202 is returned if a trip is created successfully. 
- Status 409 is returned if a passenger is already in an ongoing trip
- "http://localhost:3000/api/v1/driver/status/offline/" with a POST request is also called to update the assigned driver's status to busy. 

5. View Trip History
- viewPassengerTrips function will prompt the user for their passengerid, and with the GET request, "http://localhost:5000/api/v1/trip/", the trip history of a passenger will be scanned and listed out in reverse chronological order.

### Driver-Side (Port 3000)

6. Create Driver
 - The create driver function, will assign a random driverid to the user, prompt the user for input then with a new POST request and Driver API "http://localhost:3000/api/v1/driver/", create a new driver record in the database.
 - Status 202 is returned if a driver is created successfully.
 - Status 409 is returned if a driver already exists. 

7. Update Driver
- The updateDriver option, will prompt the user for the specific assigned driverid to make change to the existing driver. However the driver will not be able to update their Identification No. A POST request with "http://localhost:3000/api/v1/driver/" is used to update the existing driver's information. 
- Status 202 is returned if the driver is updated successfully
- Status 404 is returned if the driver does not exist.

8. View Driver (Admin - Feature)
 - The viewDriver function is called to retrieve the list of passengers that are in the database. With a new GET request and API "http://localhost:3000/api/v1/driver/view/", the list of drivers are then scanned and then displayed.
 
9. Update Busy
- Prompts the user for the unique driverid, then changes the status of the driver to Busy using the POST request and API "http://localhost:3000/api/v1/driver/status/offline/".
- Status 202 is returned if the status is updated successfully.
- Status 404 is returned if the input driverid does not exist.
 
10. Update Avaialable
 - Prompts the user for the unique driverid, then changes the status of the driver to Busy using the POST request and API "(http://localhost:3000/api/v1/driver/status/online/)".
- Status 202 is returned if the status is updated successfully.
- Status 409 is returned if the input driver has an ongoing ride.
- Status 404 is returned if the input driverid does not exist.

11. View Driver Trips 
- Retrieve assigned trips to the driver using the GET request and API "http://localhost:3000/api/v1/driver/trips/". The assigned trip is printed out. The response body is checked, if there is an assigned trip to the driver, the driver will be prompted with 3 other options which allow the driver to: 
1. Start Trip (prompt user for trip and driverid, with PUT "http://localhost:3000/api/v1/driver/trips/start/" which will change the trip status in the trip table to Ongoing. 
 - Status 202 is returned if the tripstatus is updated successfully. 
 - Status 409 is returned if the driverid / tripid does not exists 
2. End Trip (prompt user for trip and driverid, with PUT "(http://localhost:3000/api/v1/driver/trips/end/)" which will change the trip status in the trip table to Completed, and also set the driver status back to Available. 
 - Status 202 is returned if the tripstatus is updated successfully. 
 - Status 409 is returned if the driverid / tripid does not exists
3. Exit
 - This option will return the user to the main menu. 

12. Exit
 - This option allows the user to exit the application. 

# Possible Improvements in the Future: 
- Work on beautifying the console display, such as by reducing the amount of information in the main menu. 
- Work on error validation, not all of the aspects of the programme was validated due to time constraints, and I hope to have the opportunity to tackle those errors to provide better user experience
- Better structure and design and programming practices
- Continuous testing of the microservices in the future to further scout any hidden requirements 
- Implementation of front-end 
- Strive towards high coupling and low cohesion between microservices


# Instructions to set-up the microservices
1. 
