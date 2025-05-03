package main

import (
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
		client, err := s.storage.getClient(clientId)
		if err != nil {
			return err
		}
		err = verifyClient(client, publicKey)
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
	updates := []isUpdate_Type{}
	world, _ := s.storage.getWorldForClient(clientId)
	game, _ := s.storage.getGameForClient(clientId)
	if world != nil {
		switch world.Type.(type) {
		case *World_WorldOne:
			worldOne := world.GetWorldOne()
			updates = append(updates, s.newWorldOneRegionUpdate(worldOne.Region.Id))
			if game != nil {
				if clientId == game.Player1 {
					in := &ProcessWorldOneObjectRequest{
						RegionId:     worldOne.Region.Id,
						ObjectId:     WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_ONE,
						ObjectStatus: WorldOneObjectStatus_WORLD_ONE_OBJECT_STATUS_ASSIGNED,
						ObjectData:   nil,
					}
					updates = append(updates, s.newWorldOneObjectUpdate(in))
				} else if clientId == game.Player2 {
					in := &ProcessWorldOneObjectRequest{
						RegionId:     worldOne.Region.Id,
						ObjectId:     WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_ONE,
						ObjectStatus: WorldOneObjectStatus_WORLD_ONE_OBJECT_STATUS_ASSIGNED,
						ObjectData:   nil,
					}
					updates = append(updates, s.newWorldOneObjectUpdate(in))
				}
			}
			for _, object := range worldOne.Region.Objects {
				in := &ProcessWorldOneObjectRequest{
					RegionId:     worldOne.Region.Id,
					ObjectId:     object.Id,
					ObjectStatus: object.Status,
					ObjectData:   object.Data,
				}
				updates = append(updates, s.newWorldOneObjectUpdate(in))
			}
		}
	}
	if game != nil {
		switch game.Status {
		case GameStatus_READY_TO_START:
			updates = append(updates, s.newReadyToStartUpdate())
		case GameStatus_STARTED:
			updates = append(updates, s.newYouAreInGameUpdate())
		case GameStatus_WAITING_FOR_OTHER_PLAYER:
			updates = append(updates, s.newWaitingForOtherPlayerUpdate())
		}
		gameCode, _ := s.storage.getGameCodeForGame(game.Id)
		if len(gameCode) > 0 {
			updates = append(updates, s.newGameCodeGeneratedUpdate(gameCode))
		}
	}
	for _, update := range updates {
		err := stream.Send(&Update{Type: update})
		if err != nil {
			return err
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

func (s *server) enqueueUpdatesAndSignal(clientId string, updateTypes ...isUpdate_Type) {
	for _, u := range updateTypes {
		s.storage.enqueueUpdate(clientId, u)
	}
	if signal, exists := s.clientSignal.get(clientId); exists {
		select {
		case signal <- struct{}{}:
			// Signal sent
		default:
			// Non-blocking send if the channel is full
		}
	}
}
