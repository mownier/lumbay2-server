package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("unable to load env file, defaults will be used")
	}

	consumersPath := os.Getenv("LUMBAY2_SERVER_CONSUMERS")
	if len(consumersPath) == 0 {
		consumersPath = "./consumers.json"
	}

	port := os.Getenv("LUMBAY2_SERVER_PORT")
	if len(port) == 0 {
		port = "50052"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v\n", port, err)
	}

	server := newServer()

	var consumersData []byte
	if _, err := os.Stat(consumersPath); err == nil {
		consumersData, err = os.ReadFile(consumersPath)
		if err != nil {
			log.Fatalf("unable to read existing consumers file %s: %v\n", consumersPath, err)
		}
	} else if os.IsNotExist(err) {
		names := []string{"iOS", "android", "macOS", "linux", "windows", "browser"}
		list := []consumer{}
		for _, name := range names {
			privateKey, publicKey, err := server.generateKeyPair(2048)
			if err != nil {
				log.Fatalf("unable to generate key pair for %s: %v", name, err)
			}
			list = append(list, consumer{Name: name, PublicKey: publicKey, PrivateKey: privateKey})
		}
		consumersData, err = json.Marshal(list)
		if err != nil {
			log.Fatalf("unable to marshal genererated consumers file %s: %v\n", consumersPath, err)
		}
		err = os.WriteFile(consumersPath, consumersData, 0644)
		if err != nil {
			log.Fatalf("unable to write genererated consumers file %s: %v\n", consumersPath, err)
		}
	} else {
		log.Fatalf("unable to determine existence of consumers file %s: %v\n", consumersPath, err)
	}

	var consumers []*consumer
	err = json.Unmarshal(consumersData, &consumers)
	if err != nil {
		log.Fatalf("unable to unmarshal consumers from %s: %v\n", consumersPath, err)
	}
	if len(consumers) == 0 {
		log.Fatalf("no consumers found in %s\n", consumersPath)
	}

	for _, consumer := range consumers {
		server.consumers.set(consumer.PublicKey, consumer)
	}

	grpcServer := grpc.NewServer()
	RegisterLumbayLumbayServer(grpcServer, server)
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
