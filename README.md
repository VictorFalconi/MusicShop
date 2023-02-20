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
- [excelize](https://github.com/qax-os/excelize)

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
| GET    | /user/login    | Form-data, JSON |            | 201      | Login for the website (Create a new token) |
| GET    | /user          |                 | Auth       | 200      | Read information of a user                 |
| PUT    | /user          | Form-data, JSON | Auth       | 200      | Update a user                              |

**Brand**

| Method | Link        | Request         | Middleware    | Response | Description                                                   |
|--------|-------------|-----------------|---------------|----------|---------------------------------------------------------------|
| POST   | /brand      | Form-data, JSON | Auth, isAdmin | 201      | Create a new brand                                            |
| POST   | /brand/file | Form-data       | Auth, isAdmin | 201      | Create brands with CSV file (at crawl_data/product/brand.csv) |
| GET    | /brand      |                 |               | 200      | Get brands                                                    |
| GET    | /brand/:id  |                 |               | 200      | Get a brand                                                   |
| PUT    | /brand/:id  | Form-data, JSON | Auth, isAdmin | 200      | Update a brand                                                |
| DELETE | /brand/:id  |                 | Auth, isAdmin | 204      | Delete a brand                                                |

**Product**

| Method | Link          | Request         | Middleware    | Response | Description                                                                  |
|--------|---------------|-----------------|---------------|----------|------------------------------------------------------------------------------|
| POST   | /product      | Form-data, JSON | Auth, isAdmin | 201      | Create a new product                                                         |
| POST   | /product/file | Form-data       | Auth, isAdmin | 201, 207 | Create some/all product with Excel file (at crawl_data/product/product.xlsx) |
| GET    | /product      |                 |               | 200      | Get products                                                                 |
| GET    | /product/:id  |                 |               | 200      | Get a product                                                                |
| PUT    | /product/:id  | Form-data, JSON | Auth, isAdmin | 200      | Update a product                                                             |
| DELETE | /product/:id  |                 | Auth, isAdmin | 204      | Delete a product                                                             |

**Order**

| Method    | Link              | Request         | Middleware    | Response | Description                                    |
|-----------|-------------------|-----------------|---------------|----------|------------------------------------------------|
| **User**  |
| POST      | /order            | Form-data, JSON | Auth          | 201      | Create a new order                             |
| GET       | /order            |                 | Auth          | 200      | Get orders of user                             |
| GET       | /order/:id        |                 | Auth          | 200      | Get a order of user                            |
| PUT       | /order/:id        |                 | Auth          | 200      | Cancel a order of user                         |
| **Admin** |
| GET       | /admin_order      |                 | Auth, isAdmin | 200      | Get orders of users                            |
| GET       | /admin_order/:id  |                 | Auth, isAdmin | 200      | Get orders of user                             |
| PUT       | /accept_order/:id |                 | Auth, isAdmin | 200      | Accept a order of user                         |
| PUT       | /cancel_order/:id |                 | Auth, isAdmin | 200      | Cancel a order of user with all type of status |
