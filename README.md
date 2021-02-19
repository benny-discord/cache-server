# cache-server
A lightweight cross-process in-memory cache written in go.

## Running the Server
### With Docker
```shell
git clone https://github.com/benny-discord/cache-server
docker build - < Dockerfile
```
Or on Windows...
```shell
git clone https://github.com/benny-discord/cache-server -t cache-server
Get-Content Dockerfile | docker build -
```
You can then use the following command to run the server
```shell
docker run -it -p YOUR_DESIRED_PORT:7000 cache-server
```
### Without Docker
Requires Go v1.15+
```shell
git clone https://github.com/benny-discord/cache-server
go build
main.exe -p YOUR_DESIRED_PORT
```
