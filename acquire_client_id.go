package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) acquireClientId(publicKey string) (*Reply, error) {
	if _, ok := s.consumers.get(publicKey); !ok {
		return nil, status.Error(codes.InvalidArgument, "unknown public key")
	}
	client, err := s.storage.insertClient(publicKey)
	if err != nil {
		return nil, err
	}
	return s.createAcquireClientIdReply(client.Id), nil
}
