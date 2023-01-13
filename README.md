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
| POST   | /user/register      | Form-data, JSON | JSON     | Create a new "user" with rolename = user   |
| POST   | /user/login         |                 |          | Login for the website (Create a new token) |
| GET    | /user/ValidateToken |                 |          | Validate user with token                   |

