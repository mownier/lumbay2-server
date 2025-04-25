package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) acquireClientId(publicKey string) (*Reply, error) {
	if _, ok := s.consumers.get(publicKey); !ok {
		return nil, status.Error(codes.InvalidArgument, "unknown public key")
	}
	clientId := s.generateClientId(publicKey)
	// TODO: Save public key, client id
	return s.createAcquireClientIdReply(clientId), nil
}
