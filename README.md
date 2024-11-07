# greenlight

This repo contains my code from reading the book [Let's Go Further](https://lets-go-further.alexedwards.net/)

## Local database setup

```zsh
sudo -u postgres psql

CREATE DATABASE greenlight;
\c greenlight
CREATE ROLE greenlight WITH LOGIN PASSWORD '<password>';
CREATE EXTENSION IF NOT EXISTS citext;
ALTER DATABASE greenlight OWNER TO greenlight;
```
