# Backend of MusicShop Website

Build CRUD website with golang, postgresql and deploy on docker


## Crawl data using Python
- [Selenium](https://pypi.org/project/selenium/)

- [LPCLUB](https://lpclub.vn/)
- [Link crawl]()

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

