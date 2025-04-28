# Lumbay-Lumbay Game Server

## How to generate go files from the proto file

```
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative lumbay2.proto
```

## ngrok

```
ngrok http 50052 --app-protocol=http2
```