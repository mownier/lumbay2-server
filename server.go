package main

type server struct {
	UnimplementedLumbayLumbayServer
}

func newServer() *server {
	return &server{}
}
