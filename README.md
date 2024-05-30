[![Go Coverage](https://github.com/raian621/minecraft-server-controller/wiki/coverage.svg)](https://raw.githack.com/wiki/raian621/minecraft-server-controller/coverage.html)

# Go Minecraft Server Controller (GoMCSC)

GoMCSC is an HTTP and WebSocket server that essentially wraps around a Minecraft server process and controls it by writing and reading pipes to the Minecraft server's standard input, standard output, and standard error streams.

Currently, this server controller can download any version of Minecraft Java Edition's server jar and run it. This controller can also start and stop the server as well as pass commands to the Minecraft server console.

The goal of GoMCSC is to be remotely controllable by a control plane of sorts via a REST API and be able to stream the Minecraft server console input and output using WebSockets. Eventually, this server controller should also be usable as a general server controller that can be controlled via a graphical web dashboard.

## Code Generation

To generate code from the `openapi.yml` OpenAPI 3.0 spec, run this in your terminal

```sh
go get # installs the correct oapi-codegen version
oapi-codegen -config oapi-codegen.yml openapi.yml
```

## TODO

- [x] Write OpenAPI spec for the server controller's REST API
- [ ] Implement REST API endpoints
- [ ] Implement WebSocket streaming of the standard input, output, and error streams for the Minecraft server console
- [ ] Implement API key authentication
- [ ] Add support for Bedrock Minecraft Servers
- [ ] Add support for mods
- [ ] Add support for modded versions of Java Minecraft (Neoforge, Spigot, etc.)
