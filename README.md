# greenlight

This repo contains my code from reading the book [Let's Go Further](https://lets-go-further.alexedwards.net/)

Finished Nov 29 2024. This book was also a gold mine, highly recommend.

## Local database setup

```zsh
sudo -u postgres psql

CREATE DATABASE greenlight;
\c greenlight
CREATE ROLE greenlight WITH LOGIN PASSWORD '<password>';
CREATE EXTENSION IF NOT EXISTS citext;
ALTER DATABASE greenlight OWNER TO greenlight;
```

Run database migrations

```zsh
make db/migrations/up
```

## Running the server

Create a `.env` file and populate the variables from `.env.example`

```zsh
make run/api
```
