# Backend of MusicShop Website

Build CRUD website with golang, postgresql and deploy on docker


## Website
- [LPCLUB](https://lpclub.vn/)
- [Link crawl](https://github.com/MusicShopVersion1/crawl_data)

## Database
- [postgres](https://hub.docker.com/_/postgres) Database
- [gorm](https://github.com/go-gorm/gorm) ORMS
- [pgx](https://github.com/jackc/pgx) Driver & toolkit

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

| Method | Link                | Request         | Response | Decription                                 |
|--------|---------------------|-----------------|----------|--------------------------------------------|
| POST   | /user/register      | Form-data, JSON | 201      | Create a new "user" with rolename = user   |
| POST   | /user/login         | Form-data, JSON | 201      | Login for the website (Create a new token) |
| GET    | /user/ValidateToken | Cookie          | 200      | Validate user with token                   |

**Brand**

| Method | Link       | Request         | Response | Decription                |
|--------|------------|-----------------|----------|---------------------------|
| POST   | /brand     | Form-data, JSON | 201      | Create a new "user" brand |
| GET    | /brand     |                 | 200      | Get brands                |
| GET    | /brand/:id |                 | 200      | Get a brand               |
| PUT    | /brand/:id | Form-data, JSON | 200      | Update a brand            |
| DELETE | /brand/:id |                 | 204      | Delete a brand            |

