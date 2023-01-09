# Music Shop backend website with Golang & Postgres using docker
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

**Run database**
```
docker run --name my_postgres -p 41234:5432 -e POSTGRES_USER=thanhliem -e POSTGRES_PASSWORD=liem1234 -d postgres:latest
docker exec -it my_postgres createdb --username=thanhliem music_shop
```

**Run web**
```
go run main.go
```

