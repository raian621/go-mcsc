[![Go Coverage](https://github.com/raian621/minecraft-server-controller/wiki/coverage.svg)](https://raw.githack.com/wiki/raian621/minecraft-server-controller/coverage.html)

# Minecraft Server Controller

This is an HTTP and WebSocket server that essentially wraps around a Minecraft server process and controls it by writing and reading pipes to the Minecraft server's standard input, standard output, and standard error streams.

Currently, this server controller can download any version of Minecraft Java Edition's server jar and run it. This controller can also start and stop the server as well as pass commands to the Minecraft server console.

The goal of this project is for the server controller to be remotely controllable by a control plane of sorts via a REST API and stream the Minecraft server console input and output using WebSockets. Eventually, this server controller should also be usable as a general server controller that can be controlled via a graphical web dashboard.

## Starting the Server Controller

The server can be started by navigating to the project directory and executing this in your terminal:

```sh
go run .
```

## Controlling the Minecraft Server

There are currently only three HTTP endpoints:

- `POST /start`: Downloads the Minecraft server jar file if necessary and starts the Minecraft server process.
- `POST /stop`: Stops the Minecraft server by writing the `/stop` command to the standard input of the Minecraft server process
- `POST /command`: Sends a command to the Minecraft server's standard input
  - Body: `{ "command": <command here> }` 

## Code Generation

To generate code from the `openapi.yml` OpenAPI 3.0 spec, run this in your terminal

```sh
go get # installs the correct oapi-codegen version
oapi-codegen -config oapi-codegen.yml openapi.yml
```

## TODO

- [ ] Write OpenAPI spec for the server controller's REST API
- [ ] Implement REST API endpoints
- [ ] Implement WebSocket streaming of the standard input, output, and error streams for the Minecraft server console
- [ ] Implement API key authentication
- [ ] Add support for Bedrock Minecraft Servers
- [ ] Add support for mods
- [ ] Add support for modded versions of Java Minecraft (Neoforge, Spigot, etc.)
