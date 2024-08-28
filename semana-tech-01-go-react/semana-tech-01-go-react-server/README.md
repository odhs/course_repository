# Server Project Setup

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

# Development Steps

In the application directory of the application:

```sh
go mod init odhs/semana-tech-01-go-react-server-main
```

Install a tool to migrations

```sh
go install github.com/jackc/tern/v2@latest
```

Creating the dir internal

The internal directory contains all the software that is internal to your own package and that will not be imported as a dependency.

```sh
mkdir internal/store/pgstore
```

Then run the tool `tern` to create the files to migrations

```sh
tern init ./internal/store/pgstore/migrations

tern new --migrations ./internal/store/pgstore/migrations create_rooms_table

tern new --migrations ./internal/store/pgstore/migrations create_messages_table
```

The content of the 001_create_rooms_table.sql is:

```sh
001 
CREATE TABLE IF NOT EXISTS rooms (
 "id" uuid  PRIMARY KEY  NOT NULL DEFAULT gen_random_uuid(),
 "theme" VARCHAR(255)   NOT NULL
)

---- create above / drop below ----

DROP TABLE IF EXISTS rooms;
```
