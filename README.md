# Pets test case

A simple Golang HTTP REST API application that does the simple CRUD cycle within Postgresql database

## Table of Contents

- [Endpoints](#endpoints)
    - [GetPets](#getpets)
    - [CreatePet](#createpet)
    - [UpdatePet](#updatepet)
    - [DeletePet](#deletepet)
- [Error Handling](#error-handling)
- [Usage](#usage)

## Endpoints

### GetPets

- **HTTP Method:** GET
- **Route:** /pet
- **Description:** Retrieves a list of pets.
- **Parameters:**
    - `limit` (optional): Limits the number of pets returned.
    - `offset` (optional): Sets the offset for paginating through the list of pets.
    - `order` (optional): Specifies the order in which pets should be returned. Can receive "asc" and "desc" strings
- **Response:**
    - 200 OK: Returns a JSON response containing a list of pets ant total int value if pets are found.
    - 404 Not Found: Returns a "pets not found" message if no pets are found.
    - 500 Internal Server Error: Returns an error message if a database error or encoding error occurs.

### CreatePet

- **HTTP Method:** POST
- **Route:** /pet
- **Description:** Creates a new pet record.
- **Request Body:**
    - JSON object with a "name" field (string) specifying the name of the pet.
- **Response:**
    - 201 Created: Returns a JSON response containing the ID of the newly created pet.
    - 400 Bad Request: Returns an error message if the request body is missing or if the "name" field is blank.
    - 500 Internal Server Error: Returns an error message if a database error or encoding error occurs.

### UpdatePet

- **HTTP Method:** PUT
- **Route:** /pet
- **Description:** Updates an existing pet record.
- **Request Body:**
    - JSON object with "id" (number) and "name" (string) fields specifying the ID and new name of the pet.
- **Response:**
    - 200 OK: Returns a success message if the update is successful.
    - 400 Bad Request: Returns an error message if the request body is missing, if the "name" field is blank, if the "id"
  is less than or equal to 0, or if the specified pet ID does not exist.
    - 500 Internal Server Error: Returns an error message if a database error occurs.

### DeletePet

- **HTTP Method:** DELETE
- **Route:** /pet
- **Description:** Deletes an existing pet record.
- **Request Body:**
    - JSON object with an "id" field (number) specifying the ID of the pet to be deleted.
- **Response:**
    - 200 OK: Returns a success message if the deletion is successful.
    - 400 Bad Request: Returns an error message if the request body is missing, if the "id" is less than or equal to 0, 
  or if the specified pet ID does not exist.
    - 500 Internal Server Error: Returns an error message if a database error occurs.

## Error Handling

The handlers handle errors gracefully and return appropriate HTTP status codes along with error messages in case of errors. 
The error codes and messages are as follows:

- 400 Bad Request: Indicates a client-side error, such as missing or invalid request parameters.
- 404 Not Found: Indicates that the requested resource (pets) was not found.
- 500 Internal Server Error: Indicates a server-side error, such as a database error or encoding error.

## Usage

Project contains Dockerfile for the project and docker-compose file to build and run. Use command:

```shell
docker-compose up
```