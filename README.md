# Beli Mang Backend REST API

## Table of Contents

- [Description](#description)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Endpoints](#endpoints)
  - [Admin Authentication](#admin-authentication)
    - [Register](#register)
    - [Login](#login)
  - [User Authentication](#user-authentication)
    - [Register](#register-1)
    - [Login](#login-1)
  - [Image Upload](#image-upload)
  - [Managage Merchant](#managage-merchant)
    - [Add Merchant](#add-merchant)
    - [Get All Merchant](#get-all-merchant)
    - [Add Merchant Item](#add-merchant-item)
    - [Get All Merchant Items](#get-all-merchant-items)
  - [Purchase](#purchase)
    - [Get Merchant Nearby](#get-merchant-nearby)
    - [Estimate Delivery Time and Cost](#estimate-delivery-time-and-cost)
    - [Place Order](#place-order)
    - [Get Order Details of User](#get-order-details-of-user)

## Description

This API serves as the backbone for the BeliMang! food delivery application, allowing users to browse, order, and manage food and drink deliveries seamlessly. This documentation is designed to help developers integrate with our API quickly and efficiently. The API is built using Go and PostgreSQL, with the option to use AWS S3 for image uploads.

## Prerequisites

- Go 1.20 or later
- PostgreSQL 13 or later
- Docker (optional)
- AWS S3 bucket (optional; for image upload)

## Installation

Here are the steps to install the project:

1. Clone the repository:
   ```bash
   git clone https://github.com/hnfnfl/beli-mang-be.git
   ```
2. Navigate to the project directory:
   ```bash
   cd beli-mang-be
   ```
3. Install the required dependencies:
   ```bash
   go mod tidy
   ```
4. Build the binary using Makefile:
   ```bash
   make build
   ```
5. Set up the configuration file in `local_configuration/config.yaml`

6. Migrate the database schema:
   ```bash
   make migrate-up
   ```
7. Run the server:
   ```bash
    make run
   ```
8. The server should be running on `localhost:8080`

Also, you can run the server using Docker:

1. Build the Docker image:

   ```bash
   make docker-build
   ```

2. Update the environment variables in `docker-compose.yml`

3. Run the Docker container:
   ```bash
   make docker-run
   ```

## Configuration

The configuration file is located in `local_configuration/config.yaml`. This file contains the configuration for the database, logging, JWT, and AWS S3.

Here's an example of the configuration file:

```yaml
Environment: development # development, production
LogLevel: debug # debug, info, warn, error
AUTHEXPIRY: 1 # in hours

DB:
  Name: mydb # database name
  Port: 5432 # database port
  Host: localhost # database host
  Username: postgres # database username
  Password: admin # database password
  Params: sslmode=disable # database connection parameters
  MaxIdleConns: 20 # maximum idle connections
  MaxOpenConns: 20 # maximum open connections

Jwt:
  Secret: mysecret # JWT secret key
  BcryptSalt: 10 # Bcrypt salt

AWS:
  Access.Key.ID: myaccesskey # AWS access key ID
  Secret.Access.Key: mysecretaccesskey # AWS secret access key
  S3.Bucket.Name: mybucketname # AWS S3 bucket name
  Region: ap-southeast-1 # AWS region
```

## Endpoints

The following endpoints are available:

### Admin Authentication

#### Register

Endpoint: `POST /admin/register`

Request body:

```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

Response body:

```json
{
  "token": ""
}
```

#### Login

Endpoint: `POST /admin/login`

Request body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response body:

```json
{
  "token": ""
}
```

### User Authentication

#### Register

Endpoint: `POST /users/register`

Request body:

```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

Response body:

```json
{
  "token": ""
}
```

#### Login

Endpoint: `POST /users/login`

Request body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response body:

```json
{
  "token": ""
}
```

### Image Upload

Endpoint: `POST /image`

Request:

- Header: "Content-Type: multipart/form-data"
- Body: image file

Response body:

```json
{
  "message": "File uploaded sucessfully",
  "data": {
    "imageUrl": "https://awss3.d87801e9-fcfc-42a8-963b-fe86d895b51a.jpeg"
  }
}
```

### Managage Merchant

All requests here should include the JWT token in the header from `Admin`

#### Add Merchant

Endpoint: `POST /admin/merchants`

Request body:

```json
{
  "name": "",
  "merchantCategory": "",
  "imageUrl": "",
  "location": {
    "lat": 1,
    "long": 1
  }
}
```

Response body:

```json
{
  "merchantId": ""
}
```

#### Get All Merchant

Endpoint: `GET /admin/merchants`

Query parameters:

- `merchantId`: string
- `limit` & `offset`: integer; default `limit=5&offset=0`
- `name``: string
- `merchantCategory`: string
- `createdAt`: string; `asc` or `desc`

Response body:

```json
{
  "data": [
    {
      "merchantId": "",
      "name": "",
      "merchantCategory": "",
      "imageUrl": "",
      "location": {
        "lat": 1,
        "long": 1
      },
      "createdAt": ""
    }
  ],
  "meta": {
    "limit": 1,
    "offset": 0,
    "total": 10
  }
}
```

#### Add Merchant Item

Endpoint: `POST /admin/merchants/:merchantId/items`

Request body:

```json
{
  "name": "string",
  "productCategory": "",
  "price": 1,
  "imageUrl": ""
}
```

Response body:

```json
{
  "itemId": ""
}
```

#### Get All Merchant Items

Endpoint: `GET /admin/merchants/:merchantId/items`

Query parameters:

- `itemId`: string
- `limit` & `offset`: integer; default `limit=5&offset=0`
- `name`: string
- `productCategory`: string
- `createdAt`: string; `asc` or `desc`

Response body:

```json
{
  "data": [
    {
      "itemId": "",
      "name": "string",
      "productCategory": "",
      "price": 1,
      "imageUrl": "",
      "createdAt": ""
    }
  ],
  "meta": {
    "limit": 1,
    "offset": 0,
    "total": 10
  }
}
```

### Purchase

#### Get Merchant Nearby

Endpoint: `GET /merchants/nearby/${lat},${long}`

Query parameters:

- `merchantId`: string
- `limit` & `offset`: integer; default `limit=5&offset=0`
- `name`: string
- `merchantCategory`: string

Response body:

```json
{
  "data": [
    {
      "merchant": {
        "merchantId": "",
        "name": "",
        "merchantCategory": "",
        "imageUrl": "",
        "location": {
          "lat": 1,
          "long": 1
        },
        "createdAt": ""
      },
      "items": [
        {
          "itemId": "",
          "name": "string",
          "productCategory": "",
          "price": 1,
          "imageUrl": "",
          "createdAt": ""
        }
      ]
    }
  ]
}
```

#### Estimate Delivery Time and Cost

Endpoint: `POST /users/estimate`

Request body:

```json
{
  "userLocation": {
    "lat": 1,
    "long": 1
  },
  "orders": [
    {
      "merchantId": "string",
      "items": [
        {
          "itemId": "string",
          "quantity": 1
        }
      ]
    }
  ]
}
```

Response body:

```json
{
  "totalPrice": 1,
  "estimatedDeliveryTimeInMinutes": 1,
  "calculatedEstimateId": ""
}
```

#### Place Order

Endpoint: `POST /users/orders`

Request body:

```json
{
  "calculatedEstimateId": ""
}
```

Response body:

```json
{
  "orderId": ""
}
```

#### Get Order Details of User

Endpoint: `GET /users/orders`

Query parameters:

- `merchantId`: string
- `limit` & `offset`: integer; default `limit=5&offset=0`
- `name`: string
- `merchantCategory`: string

Response body:

```json
[
  {
    "orderId": "string",
    "orders": [
      {
        "merchant": {
          "merchantId": "",
          "name": "",
          "merchantCategory": "",
          "imageUrl": "",
          "location": {
            "lat": 1,
            "long": 1
          },
          "createdAt": ""
        },
        "items": [
          {
            "itemId": "",
            "name": "string",
            "productCategory": "",
            "price": 1,
            "imageUrl": "",
            "createdAt": ""
          }
        ]
      }
    ]
  }
]
```
