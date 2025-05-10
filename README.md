# Lumbay-Lumbay Game Server

## How to generate go files from the proto file

```
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative lumbay2.proto
```

## ngrok

You must register first in ngrok website. Follow their instructions on how to install the ngrok cli. Then, do this.

```
ngrok http 50052 --app-protocol=http2

ngrok http --url=<free_static_url> 50052 --upstream-protocol=http2
```

## env

You must define first an env file. If environment variable is set to empty, the program will use the default value.

```
# .env
LUMBAY2_SERVER_CONSUMERS=
LUMBAY2_SERVER_PORT=
LUMBAY2_SERVER_DB=
```