# Product API

## Architecture

This API uses a **clean layered architecture**:

- **Handler Layer:** Receives and validates HTTP requests, parses queries, and sends structured responses.
- **Service Layer:** Contains business logic, validation checks, and rules for create and list operations.
- **Repository Layer:** Handles direct database operations with PostgreSQL/MySQL using parameterized queries.
- **Redis Layer:** Caches list product request parameters and results for 1 minute to reduce DB load for repeated requests.
- **Logger Layer:** Centralized logger for debugging, error tracking, and structured logging.

### Why This Architecture?

- **Maintainability:** Clear separation of concerns making the codebase easy to maintain and extend.
- **Scalability:** Easy to replace the repository layer to switch databases or caching strategies.
- **Performance:** Redis caching minimizes database hits for repeated list queries.
- **Testability:** Each layer can rawuh-service unit tested in isolation.
- **Observability:** Logger integration enables easier debugging in development and tracing in production environments.

## Overview

This project creates **two APIs**:

- `POST /product/create` for creating products
- `GET /product/list` for listing products

It uses **Redis** to store request parameters or keys with the same request for **1 minute** for caching purposes.
It also uses a logger to display errors, debug information, and log request/response flows for easier debugging and monitoring.

## How to Run

1. Create the table under db public.products:

```sql
CREATE TABLE IF NOT EXISTS public.products (
    product_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(500),
    price DOUBLE PRECISION,
    description VARCHAR(500),
    quantity INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

2. Run your server on `localhost:8080`.

---

## API Details

### Create Product

`POST localhost:8080/product/create`

**Validations:**

- `Name` cannot rawuh-service empty.
- Length of `Name` and `Description` is checked against env/DB length.
- Duplicate `Name` is not allowed.
- `Price` and `Quantity` must rawuh-service > 0.
- Regex validation:
  - `Name`: `^[a-zA-Z0-9 _.,'-]+$`
  - `Description`: `^[\\s\\w\\d_.,-;()/]*$`

**Body Example:**

```json
{
  "Name": "Macbook Pro M4 6 TB",
  "Price": 40000000,
  "Description": "MAX PRO MAX BROO",
  "Quantity": 2
}
```

**Response Success:**

```json
{
  "Error": false,
  "Code": 200,
  "Message": "Success Create Product"
}
```

**Response Error:**

```json
{
  "Error": true,
  "Code": 403,
  "Message": "rpc error: code = Aborted desc = product name is empty"
}
```

---

### Get Product List

`GET localhost:8080/product/list?sort=created_at&dir=asc&page=2&limit=3&query=<base64>`

**Sorting & Direction:**

- Allowed columns:
  - `created_at`, `price`, `name`, `product_id`, `quantity`
- Allowed directions:
  - `asc`, `desc`
- Supports sorting:
  - Product newest (created_at desc)
  - Product lowest price (price asc)
  - Product highest price (price desc)
  - Product name (A-Z, Z-A)

**Query Filtering:**

- Uses Base64 encoded string for `name`, `description` filtering.
- Example filter: `name,description:%!Macbook Pro` (encoded to Base64) for `ILIKE` filtering in DB.

---

**Response Success:**

```json
{
  "Error": false,
  "Code": 200,
  "Message": "Success",
  "Data": [
    {
      "ProductID": 4,
      "Name": "Macbook Pro M4 512 GB",
      "price": 23400000,
      "Description": "MacBook Pro 14 inci memiliki tiga port Thunderbolt 4, port pengisian daya MagSafe 3, slot kartu SDXC, port HDMI, dan jek headphone.",
      "quantity": 234,
      "CreatedAt": "2025-07-17T22:30:10.51959+07:00",
      "UpdatedAt": "2025-07-17T22:30:10.51959+07:00"
    },
    {
      "ProductID": 10,
      "Name": "Macbook Pro M4 5 TB",
      "price": 34000000,
      "Description": "MAX PRO MAX BROO",
      "quantity": 2,
      "CreatedAt": "2025-07-17T22:32:27.175526+07:00",
      "UpdatedAt": "2025-07-17T22:32:27.175526+07:00"
    },
    {
      "ProductID": 11,
      "Name": "Macbook Pro M4 6 TB",
      "price": 40000000,
      "Description": "MAX PRO MAX BROO",
      "quantity": 2,
      "CreatedAt": "2025-07-17T22:32:38.106871+07:00",
      "UpdatedAt": "2025-07-17T22:32:38.106871+07:00"
    }
  ],
  "Pagination": {
    "Limit": 3,
    "Page": 2,
    "TotalRows": 6,
    "TotalPages": 2
  }
}
```

**Response Error:**

```json
{
  "Error": true,
  "Code": 500,
  "Message": "rpc error: code = InvalidArgument desc = Invalid Argument",
  "Data": null,
  "Pagination": null
}
```

---

This README keeps your test assignment clear, aligned with requirements, and structured for direct implementation and submission.
