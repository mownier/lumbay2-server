package main

import (
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *server) Subscribe(emp *Empty, stream LumbayLumbay_SubscribeServer) error {
	select {
	case <-stream.Context().Done():
		return status.Error(codes.Canceled, "subscribe was cancelled")

	default:
		clientId := ""
		publicKey := ""
		md, ok := metadata.FromIncomingContext(stream.Context())
		if ok {
			values := md.Get("public_key-bin")
			if len(values) > 0 {
				publicKey = values[0]
			}
			values = md.Get("client_id-bin")
			if len(values) > 0 {
				clientId = values[0]
			}
		}
		if len(clientId) == 0 {
			return status.Error(codes.InvalidArgument, "failed to subscribe because client id is unknown")
		}
		if len(publicKey) == 0 {
			return status.Error(codes.InvalidArgument, "failed to subscribe because public key is unknown")
		}
		err := verifyClientId(clientId, publicKey)
		if err != nil {
			return err
		}
		if _, exists := s.clientSignal.get(clientId); !exists {
			s.clientSignal.set(clientId, make(chan struct{}, 1))
		}
		defer s.cleanUpResources(clientId)
		if err := s.sendInitialUpdates(clientId, stream); err != nil {
			return err
		}
		signal, _ := s.clientSignal.get(clientId)

		for {
			select {
			case <-stream.Context().Done():
				return status.Error(codes.Internal, "subscribe was done")

			case <-signal:
				if err := s.sendUpdates(clientId, stream); err != nil {
					return err
				}
			}
		}
	}
}

func (s *server) cleanUpResources(clientId string) {
	s.clientSignal.delete(clientId)
}

func (s *server) sendInitialUpdates(clientId string, stream LumbayLumbay_SubscribeServer) error {
	list := []*Update{
		{
			Type: &Update_YouAreInGameUpdate{
				YouAreInGameUpdate: &YouAreInGameUpdate{},
			},
		},
		{
			Type: &Update_Ping{
				Ping: &Ping{},
			},
		},
	}
	for _, update := range list {
		if err := stream.Send(update); err != nil {
			log.Printf("unable to send initial updates for client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to send intial updates")
		}
	}
	return nil
}

func (s *server) sendUpdates(clientId string, stream LumbayLumbay_SubscribeServer) error {
	updates, err := s.storage.getAllUpdates(clientId)
	if err != nil {
		return err
	}
	if len(updates) == 0 {
		return nil
	}
	updatesToDequeue := []*Update{}
	for _, update := range updates {
		err = stream.Send(update)
		if err != nil {
			break
		}
		updatesToDequeue = append(updatesToDequeue, update)
	}
	dequeueErr := s.storage.dequeueUpdates(clientId, updatesToDequeue)
	if err != nil {
		return err
	}
	return dequeueErr
}
