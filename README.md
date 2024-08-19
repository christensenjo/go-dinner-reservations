# Dinner Reservations API - Written in Golang

Let's *Go* Eat! 😋🍽️

This repository was my first stab at implementing an API in Golang. It creates a CRUD API on a local PostgreSQL database instance which can be interacted with.

## Setup

First, make sure that you have PostgreSQL installed. You can do so [here](https://www.postgresql.org/download/). You'll need to run the correct executable on Windows, use Homebrew on Mac, or sudo on Linux.

Second, verify the installation. You can do so by running `psql` in your terminal or command prompt. (eg Powershell on Windows) 

> [!NOTE]
> *Windows Users*
> 
> You may need to add psql to your PATH variable on windows if the command is not recognized after installing. Navigate to environment variables, find PATH, edit it, and add:
> ```
> C:/Program Files/PostgreSQL/{version_num}/bin
> ```
> Or the correct alternative bin location on your machine. Restart your terminal.

Once you've verified `psql` works, create a database.
`CREATE DATABASE dinner_reservations`

Create a user. Replace myuser and mypassword with the login credentials you want.
```
CREATE USER myuser with PASSWORD 'mypassword';
GRANT ALL PRIVILEGES ON DATABASE dinner_reservations TO myuser;
```
*You may need to play with the permissions on myuser. If so, login as the root user and GRANT privileges to the user both on the correct schema and on the table itself if necessary.*

Next, create the reservations table.
```
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    date DATE,
    time TIME,
    guests INTEGER,
    phone VARCHAR(15)
)
```

Optionally, insert and query some data manually.
```
INSERT INTO reservations (name, date, time, guests, phone) VALUES
('John Doe', '2024-08-03', '19:00', 4, '123-456-7890');
```
`SELECT * FROM reservations;`

Lastly, I've setup the db authentication in Go to use a gitignored .env file, so you'll need to create that:
eg
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_real_username
DB_PASSWORD=your_real_password
DB_NAME=dinner_reservations

```
Make sure to include the correct credentials from when you setup your dinner_reservations database.


At this point, you should be ready to use the API! From the repository execute `go run .`, then send requests from Postman or a browser, etc.


## Endpoints

### GET /reservations
Returns a list of all dinner reservations.

Sample Request:
```
GET /reservations
```

### GET /reservations/{id}
Returns a specific dinner reservation based on the provided ID.

Sample Request:
```
GET /reservations/1
```

### POST /reservations
Creates a new dinner reservation.

Sample Request:
```
POST /reservations
Content-Type: application/json

{
    "name": "John Doe",
    "date": "2022-10-15",
    "time": "19:00",
    "partySize": 4
}
```

### PUT /reservations/{id}
Updates an existing dinner reservation based on the provided ID.

Sample Request:
```
PUT /reservations/1
Content-Type: application/json

{
    "name": "Jane Smith",
    "date": "2022-10-15",
    "time": "19:30",
    "partySize": 6
}
```

### DELETE /reservations/{id}
Deletes a specific dinner reservation based on the provided ID.

Sample Request:
```
DELETE /reservations/1
```
