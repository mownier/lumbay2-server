package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) acquirePublicKey(in *AcquirePublicKeyRequest) (*Reply, error) {
	publicKey := ""
	s.consumers.forEach(func(k string, v *consumer) bool {
		if v.Name == in.Name {
			publicKey = v.PublicKey
			return true
		}
		return false
	})
	if len(publicKey) == 0 {
		return nil, status.Error(codes.NotFound, "consumer not found")
	}
	return s.createAcquirePublickKeyReply(publicKey), nil
}
