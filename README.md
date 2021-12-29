# ScaleFlix

Media APIs for managing Movies and Series. 

# How to install

git clone https://github.com/havvaozdemir/scaleflixapi.git

# Requirements

Install the requirements listed below.

* Golang 
* postgreSQL
    Create Database with the information also listed in .env file;
    databese : postgresql
    user: *****
    password: ****

## How to use
* For Movie Library, you need api key, you can get it from http://www.omdbapi.com. Set API_KEY config and .env file.

* go build, run , test options are in Makefile
    >Make build
    >Make run
    >Make test
    >Make swagger
    >Make run-compose
    >Make run-docker

* Start application with : 
    Make run or
    Make run-compose or
    Make run-docker

* Go to http://localhost:8080

* Insert database user and admin roles to your database. if not, you can not do some functinalities. Use scripts file from utils/scripts

* Authentication: 

    Get token as a user or admin and add Authorization header with Barear and token.
    
By using the endpoints listed below; you can search movies and series, add or remove movies and series to a favorite list as a user role. You can search movies and series from [http://omdbapi.com/] library, add or remove them to the system as an admin role.

## Endpoint Table

| Endpoint        | Method | Description                       |
| ----------------|--------|-----------------------------------|
| /token          | GET    | Returns token for authorization   |
| /movies         | GET    | Get movies list                   |
| /movies         | POST   | Add movie to the system           |
| /movies/{id}    | GET    | Get movie by ID                   |
| /series         | GET    | Get series list                   |
| /series         | POST   | Add series and details to system  |
| /series/{id}    | GET    | Get series by ID                  |
| /movies/{id}    | DELETE | Remove movie from system by ID    |
| /series/{id}    | DELETE | Remove series from system by ID   |
| /suggestions    | GET    | Get movies and series from library|
| /favorites      | GET    | Get movies and series from favorite list|
| /favorites      | POST   | Add movies and series to favorite list|
| /favorites/{id} | DELETE | Remove movie or series from favorite list|
