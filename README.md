# Backend of MusicShop Website

Build CRUD website with golang, postgresql and deploy on docker


## Website
- [LPCLUB](https://lpclub.vn/)
- [Link crawl](https://github.com/MusicShopVersion1/crawl_data)

## Database
- [postgres](https://hub.docker.com/_/postgres) Database
- [gorm](https://github.com/go-gorm/gorm) ORM
- [pgx](https://github.com/jackc/pgx) Driver & toolkit

![Database Diagram](https://github.com/MusicShopVersion1/server/blob/master/images/MusicShop%20Database.png)

## Framework & Library
- [gin](https://github.com/gin-gonic/gin)
- [validator](https://github.com/go-playground/validator)
- [jwt](https://github.com/golang-jwt/jwt)

[//]: # (- [migrate]&#40;https://github.com/golang-migrate/migrate&#41;)

### Create image

**Golang**
```
docker build --tag golang .
```

### Run

**Run database**
```
docker compose up
```

**Run web**

docker network list #(server_default)
```
docker run --net server_default --name backend_golang -p 6868:6868 -d golang:latest
```

## APIs

**User**

| Method | Link           | Request         | Middleware | Response | Description                                |
|--------|----------------|-----------------|------------|----------|--------------------------------------------|
| POST   | /user/register | Form-data, JSON |            | 201      | Create a new "user" with rolename = user   |
| POST   | /user/login    | Form-data, JSON |            | 201      | Login for the website (Create a new token) |

**Brand**

| Method | Link       | Request         | Middleware     | Response | Description               |
|--------|------------|-----------------|----------------|----------|---------------------------|
| POST   | /brand     | Form-data, JSON | Token, isAdmin | 201      | Create a new "user" brand |
| GET    | /brand     |                 | Token          | 200      | Get brands                |
| GET    | /brand/:id |                 | Token          | 200      | Get a brand               |
| PUT    | /brand/:id | Form-data, JSON | Token, isAdmin | 200      | Update a brand            |
| DELETE | /brand/:id |                 | Token, isAdmin | 200      | Delete a brand            |

