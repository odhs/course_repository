# Server Project Setup

## The Environment

The environment is defined in `compose.yml` to run with **Docker**, the are a **PostgresSQL** database and a **PG Admin**, to manage the database in address `localhost:8081`.

To start the environment use:

```sh
docker compose start
```

Or if you want to see the log, but you will not be able to use the terminal at sametime to develop, use:

```sh
docker compose up
```

## Install Golang in Ubuntu

Download Golang

```sh
curl -OL <https://go.dev/dl/go1.23.0.linux-amd64.tar.gz>
sudo tar -C /usr/local -xvf go1.23.0.linux-amd64.tar.gz
```

## Define Environment Variables

Put on the end of the file

```sh
export PATH=$PATH:~/go/bin
export PATH=$PATH:/usr/local/go/bin
```

Run for the terminal to recognize the changes:

```sh
source ~/.profile
```

Check if it is working

```sh
go version
```

# Foundation - Development Steps

In the application directory of the application:

```sh
go mod init odhs/semana-tech-01-go-react-server-main
```

## Migrations

### Tern tool 

Install the tern package for use with migrations

```sh
go install github.com/jackc/tern/v2@latest
```

Create the internal directory

The internal directory contains all the software that is internal to your own package and that will not be imported as a dependency.

```sh
mkdir internal/store/pgstore
```

Then run the tool `tern` to create the files to migrations

```sh
tern init ./internal/store/pgstore/migrations
```

Delete the file .sql in the `/internal/store/pgstore/migrations` and create new ones.

```sh
tern new --migrations ./internal/store/pgstore/migrations create_rooms_table

tern new --migrations ./internal/store/pgstore/migrations create_messages_table
```

The content of the `001_create_rooms_table.sql` is:

```sql
001 
CREATE TABLE IF NOT EXISTS rooms (
 "id" uuid  PRIMARY KEY  NOT NULL DEFAULT gen_random_uuid(),
 "theme" VARCHAR(255)   NOT NULL
)

---- create above / drop below ----

DROP TABLE IF EXISTS rooms;
```

```sql
The content of the `001_create_messsages_table.sql` is:

CREATE TABLE IF NOT EXISTS messages(
	"id" 				uuid	PRIMARY KEY 	NOT NULL	DEFAULT gen_random_uuid(),
	"room_id"			uuid					NOT NULL,
	"message"			VARCHAR(255)			NOT NULL,
	"reaction_count"	BIGINT					NOT NULL 	DEFAULT 0,
	"answered"			BOOLEAN					NOT NULL 	DEFAULT false,
	FOREIGN KEY (room_id) REFERENCES rooms(id)
) 

---- create above / drop below ----

DROP TABLE IF EXISTS messages;
```

### Install the package SQLC

1. SQLC to generates the library to access the database

```sh
 go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

2. Creates the config file `sqlc.yaml` in `pgstore` dir and
the file `./internal/pgstore/queries/queries.sql` with the content:

```sql
-- name: GetRoom :one
SELECT "id", "theme" FROM rooms WHERE id = $id;
```

Other SQL was added laterin the queries.sql

3. Run SQLC

```sh
sqlc generate -f ./internal/store/pgstore/sqlc.yaml
```

## Automate Tern and SQLC

To automatically run the **Tern** and **SQLC** with the **GO Generate** tool the file `gen.go` was created at the root of the project with the content:

```go
//go:generate go run ./cmd/tools/terndotenv/main.go
//go:generate sqlc generate -f ./internal/store/pgstore/sqlc.yaml
```

So when executing `GO GENERATE./...` everything will be executed.

# API - Development Steps

The API handle was created in the internal/api directory

## Add Chi

## Add Chi CORS 

go get github.com/go-chi/cors