# cache-server

A lightweight cross-process in-memory cache written in go.

## Running the Server

### With Docker

Premade Docker builds: https://github.com/benny-discord/cache-server/packages/631531

```shell
$ docker login docker.pkg.github.com #login with your github credentials here
$ docker run --name YOUR_CONTAINER_NAME -p YOUR_DESIRED_PORT:7000 docker.pkg.github.com/benny-discord/cache-server/cache:latest
```

### Without Docker

Requires Go v1.12+

From the CLI

```shell
$ git clone https://github.com/benny-discord/cache-server
$ go build
$ main.exe -p YOUR_DESIRED_PORT
```

From the ZIP:

- Download and extract https://github.com/benny-discord/cache-server/archive/master.zip
- Enter the folder
- Run `go build`
- Run `main.exe`

## Usage

You can connect to the server over a Websocket connection. Authentication is not handled by the server. It is
recommended that you use a reverse proxy such as NginX to handle the authentication for you.

### Messages

| op       | Description                                                                                                                          | Direction | Example                                                                                                                 |
|----------|--------------------------------------------------------------------------------------------------------------------------------------|-----------|-------------------------------------------------------------------------------------------------------------------------|
| GET      | Used for requesting a cached value                                                                                                   | Sent      | `{"op":"GET", "key": "foo"}`                                                                                            |
| SET      | Used to set a value. Has an optional `expires` field that accepts a Unix timestamp                                                   | Sent      | `{"op": "SET", "key": "foo", "value": "bar", "expires": 1613821456922}`                                                 |
| DELETE   | Used to delete a cached value                                                                                                        | Sent      | `{"op": "DELETE", "key": "foo"}`                                                                                        |
| CLEAR    | Clears the cache of the entire server                                                                                                | Sent      | `{"op": "CLEAR"}`                                                                                                       |
| RESPONSE | Responds to a GET request that the client sent                                                                                       | Received  | `{"op": "RESPONSE", "key": "foo", "value": "bar"}`                                                                      |
| WARN     | Receives a warning from the server                                                                                                   | Received  | `{"op": "WARN", "message": "Warning: valid op property must be present in payload"}`                                    |
| ERROR    | Receives an error from the server, regarding an invalid payload the client sent. The client will be disconnected after this is sent. | Received  | `{"op": "ERROR", "message": "json: cannot unmarshal number into Go struct field wsRequestPayload.value of type string"` |

When you receive a warning, you should review your code to find where it came from. A list of all possible warnings is
available below.

- {"op": "WARN", "message": "Warning: valid op property must be present in payload"}
- {"op": "WARN", "message": "Warning: expires property should not be present when op is `CLEAR|DELETE|GET`"}
- {"op": "WARN", "message": "Warning: value property should not be present when op is `CLEAR|DELETE|GET`"}
- {"op": "WARN", "message": "Warning: key property must be present when op is `DELETE|GET|SET`"}
- {"op": "WARN", "message": "Warning: expires property should be greater than current time"}

When receiving an error, the connection to the server will be terminated, the client must fix the payload and reconnect.
Errors are only received when sending malformed JSON payloads. Here is the structure expected.
```json5
{
  "key": "string", //Required on DELETE, GET, SET. Must be a string
  "value": "string", //Required on SET. Must be a string.
  "op": "string", //Must be one of CLEAR|DELETE|GET|SET. Must be a string.
  "expires": "integer" //Unix timestamp, example: 1613821456922. Must be an integer
}
```
