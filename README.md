# cache-server

A lightweight cross-process in-memory cache written in go.

## Running the Server

### With Docker

Premade Docker builds: https://hub.docker.com/repository/docker/alex31r/cache/builds

```shell
$ docker run --name YOUR_CONTAINER_NAME -p YOUR_DESIRED_PORT:7000 alex31r/cache:latest 
```

### Without Docker

Requires Go v1.15+

From the CLI

```shell
git clone https://github.com/benny-discord/cache-server
go build
main.exe -p YOUR_DESIRED_PORT
```

From the ZIP:

- Download and extract https://github.com/benny-discord/cache-server/archive/master.zip
- Enter the folder
- Run `go build`
- Run `main.exe`
