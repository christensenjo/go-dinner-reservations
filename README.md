# Dinner Reservations API - Written in Golang

Let's *Go* Eat! 😋🍽️

This repository was my first stab at implementing an API in Golang. It implements a CRUD API with mockup dinner reservation data that can be queried.

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
