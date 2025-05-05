package main

import "google.golang.org/grpc/codes"

func (s *server) exitWorld(clientId string) (*Reply, error) {
	game, err := s.storage.getGameForClient(clientId)
	if err != nil {
		return nil, err
	}
	world, err := s.storage.getWorldForClient(clientId)
	if err != nil {
		return nil, err
	}
	switch world.Type.(type) {
	case *World_WorldOne:
		worldOne := world.GetWorldOne()
		if !worldOne.viableForLife() {
			return nil, sverror(codes.InvalidArgument, "failed to exit world", nil)
		}
		err := s.storage.removeWorldForClient(clientId)
		if err != nil {
			return nil, err
		}
		if clientId == game.Player1 {
			s.enqueueUpdatesAndSignal(game.Player1, s.newYouExitWorldUpdate())
			s.enqueueUpdatesAndSignal(game.Player2, s.newOtherExitsWorldUpdate())
		} else if clientId == game.Player2 {
			s.enqueueUpdatesAndSignal(game.Player1, s.newOtherExitsWorldUpdate())
			s.enqueueUpdatesAndSignal(game.Player2, s.newYouExitWorldUpdate())
		}
	}
	return s.newExitWorldReply(), nil
}
