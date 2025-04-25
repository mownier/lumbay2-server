package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/hashicorp/mdns"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load env file")
	}

	port := os.Getenv("LUMBAY2_SERVER_PORT")
	if len(port) == 0 {
		port = "50052"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		portInt = 50052
	}

	serviceName := os.Getenv("LUMBAY2_SERVER_SERVICE_NAME")
	if len(serviceName) == 0 {
		serviceName = "_lumbay._tcp"
	}

	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("unable to get hostname: %v\n", err)
	}

	mdnsService, err := mdns.NewMDNSService(host, serviceName, "", "", portInt, nil, nil)
	if err != nil {
		log.Fatalf("unable to create mDNS service: %v\n", err)
	}

	mdnsServer, err := mdns.NewServer(&mdns.Config{Zone: mdnsService})
	if err != nil {
		log.Fatalf("unable to create mDNS server: %v\n", err)
	}
	defer mdnsServer.Shutdown()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v\n", port, err)
	}

	grpcServer := grpc.NewServer()
	RegisterLumbayLumbayServer(grpcServer, newServer())
	log.Printf("server listening on port %s\n", port)

	shutdownSigChan := make(chan os.Signal, 1)
	signal.Notify(shutdownSigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownSigChan
		log.Println("server will stop listening")
		grpcServer.GracefulStop()
		log.Println("server stopped listening")
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}

	log.Println("server shutdown")
}
