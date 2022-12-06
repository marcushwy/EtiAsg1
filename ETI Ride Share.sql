CREATE DATABASE my_db; 

USE my_db; 

SET SQL_SAFE_UPDATES = 0;

CREATE TABLE `passenger` (
  `ID` varchar(5) NOT NULL,
  `FirstName` varchar(30) DEFAULT NULL,
  `LastName` varchar(30) DEFAULT NULL,
  `MobileNumber` varchar(30) DEFAULT NULL,
  `Email` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `driver` (
  `ID` varchar(5) NOT NULL,
  `FirstName` varchar(30) DEFAULT NULL,
  `LastName` varchar(30) DEFAULT NULL,
  `MobileNumber` varchar(30) DEFAULT NULL,
  `Email` varchar(50) DEFAULT NULL,
  `CarLicense` varchar(10) DEFAULT NULL,
  `IcNo` varchar(10) DEFAULT NULL,
  `DriverStatus` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `trip` (
  `TripId` varchar(10) NOT NULL,
  `PassengerId` varchar(10) DEFAULT NULL,
  `PickupCode` varchar(7) DEFAULT NULL,
  `DropoffCode` varchar(7) DEFAULT NULL,
  `DriverId` varchar(10) DEFAULT NULL,
  `TripStatus` varchar(20) DEFAULT NULL,
  `TripDate` datetime DEFAULT NULL,
  PRIMARY KEY (`TripId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci