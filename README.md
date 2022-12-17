
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

 ## List of API and functions called when type of request is made:
 
- ("/api/v1/passenger/view/", viewpassenger)
- ("/api/v1/passenger/{passengerid}", passenger)
- ("/api/v1/passenger/trip/{tripid}/{passengerid}/{driverid}", newtrip)
- ("/api/v1/trip/{passengerid}", viewpassengertrip)


# Driver Service Documentation 

 ## List of API and functions called when type of request is made:
 
- router.HandleFunc("/api/v1/driver/view/", viewdriver)
-	router.HandleFunc("/api/v1/driver/{driverid}", driver)
-	router.HandleFunc("/api/v1/driver/status/online/{driverid}", driveronline)
-	router.HandleFunc("/api/v1/driver/status/offline/{driverid}", driveroffline)
-	router.HandleFunc("/api/v1/drivers/", autoassigndriver)
-	router.HandleFunc("/api/v1/driver/trips/{driverid}", getdrivertrip)
-	router.HandleFunc("/api/v1/driver/trips/start/{tripid}", begintrip)
-	router.HandleFunc("/api/v1/driver/trips/end/{tripid}/{driverid}", finishtrip)

# Trip Service
 
## List of Options in the main menu and how they work: 

-


# Instructions to set-up the microservices
- 
