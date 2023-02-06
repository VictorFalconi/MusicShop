# Backend of MusicShop Website

Build CRUD website with golang, postgresql and deploy on docker


## Website
- [LPCLUB](https://lpclub.vn/)
- [Link crawl](https://github.com/MusicShopVersion1/crawl_data)

## Database
- [postgres](https://hub.docker.com/_/postgres) Database
- [gorm](https://github.com/go-gorm/gorm) ORM
- [pgx](https://github.com/jackc/pgx) Driver & toolkit

![Database Diagram](https://github.com/MusicShopVersion1/server/blob/master/images/Web%20Online%20MusicShop.png)

## Framework & Library
- [gin](https://github.com/gin-gonic/gin)
- [validator](https://github.com/go-playground/validator)
- [jwt](https://github.com/golang-jwt/jwt)

[//]: # (- [migrate]&#40;https://github.com/golang-migrate/migrate&#41;)

### Run project

```
docker compose up
```

## APIs

**User**

| Method | Link           | Request         | Middleware | Response | Description                                |
|--------|----------------|-----------------|------------|----------|--------------------------------------------|
| POST   | /user/register | Form-data, JSON |            | 201      | Create a new "user" with rolename = user   |
| POST   | /user/login    | Form-data, JSON |            | 201      | Login for the website (Create a new token) |

**Brand**

| Method | Link        | Request         | Middleware     | Response | Description                 |
|--------|-------------|-----------------|----------------|----------|-----------------------------|
| POST   | /brand      | Form-data, JSON | Token, isAdmin | 201      | Create a new brand          |
| POST   | /brand/file | Form-data       | Token, isAdmin | 201      | Create brands with CSV file |
| GET    | /brand      |                 | Token          | 200      | Get brands                  |
| GET    | /brand/:id  |                 | Token          | 200      | Get a brand                 |
| PUT    | /brand/:id  | Form-data, JSON | Token, isAdmin | 200      | Update a brand              |
| DELETE | /brand/:id  |                 | Token, isAdmin | 204      | Delete a brand              |

**Product**

| Method | Link          | Request         | Middleware     | Response | Description                             |
|--------|---------------|-----------------|----------------|----------|-----------------------------------------|
| POST   | /product      | Form-data, JSON | Token, isAdmin | 201      | Create a new product                    |
| POST   | /product/file | Form-data       | Token, isAdmin | 201, 207 | Create some/all product with Excel file |
| GET    | /product      |                 | Token          | 200      | Get products                            |
| GET    | /product/:id  |                 | Token          | 200      | Get a product                           |
| PUT    | /product/:id  | Form-data, JSON | Token, isAdmin | 200      | Update a product                        |
| DELETE | /product/:id  |                 | Token, isAdmin | 204      | Delete a product                        |
