# Chirpy GO
## _Backend server connected to DB - chirpy go_

Project of server backend similar to twitter functions connected to PostgreSQL and have authentication with email and password to login, and JWT to generate chirps.

Go lang and PostgreSQL.

## Features

- Handlers to work with APIs give request and respond with JSON
- File Server path to access documents
- Create Users and Chirps by User saved in Database
- Authentication by JWT to create chirps

## Tech

Gator uses modules to work properly

- [pq](github.com/lib/pq) - PostgreSQL Library to connect to Database
- [uuid](github.com/google/uuid) - Google UUID Library that add uuid generation functionality
- [dotenv](https://github.com/joho/godotenv) - Library that can use environment variables from .env files
- [argon2id](github.com/alexedwards/argon2id) - Library to hash password and compare hashed password for validation
- [jwt](github.com/golang-jwt/jwt/v5) - Library to use JWT for Authentication

## Installation

Chirpy go require postgreSQL, Go and Goose
To develop more functionalities with the database it will also be necessary SQLC command line tool

[PostgreSQL](https://www.postgresql.org/download/) - Download and install

[Golang download doc](https://go.dev/doc/install) - Download files and follow instructions

Download or clone repository

Install goose by command line:
```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Install SQLC by command line:
```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Generate env file with variables:
- DB_URL = url to connect to Database
- PLATFORM = Encironment that is being used, to develop and test `dev`
- SECRET = String used for JWT
- POLKA_KEY = String used for Webhook

## Endpoints

- **/app/** - File Server
- **GET api/healthz** - Check Server Status
- **GET /admin/metrics** - Check how many visits given to File Server
- **POST /admin/reset** - Reset visits to File Server
- **POST /api/users** - Create New User, hash password and save in DB
- **PUT /api/users** - Update User email and Password, needs valid JWT Auth header
- **POST /api/polka/webhooks** - Change User is chirpy red to True, used by service
- **POST /api/chirps** - Create Chirps and save in DB, needs valid JWT Auth header
- **GET /api/chirps** - Get JSON with all Chirps, can use query parameters author_id and sort to just show Chirps from specific user or sort in descending order
- **GET /api/chirps/{chirpID}** - Get JSON with Specific Chirp by Chrip ID
- **DELETE /api/chirps/{chirpID}** - Delete Chirp by ID
- **POST /api/login** - Login by Valid email and password, return User, JWT and Refresh token. JWT has a duration of 1 hour and Refresh Token of 60 days
- **POST /api/refresh** - Refresh JWT, require valid Refresh Token find in DB, not revoked or pass Expiration Time
- **POST /api/revoke** - Revoke Refresh Token