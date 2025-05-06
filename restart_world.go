package main

import "google.golang.org/grpc/codes"

func (s *server) restartWorld(clientId string) (*Reply, error) {
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
		if !worldOne.gameIsOver() {
			return nil, sverror(codes.InvalidArgument, "failed to restart world", nil)
		}
		if clientId == game.Player1 {
			if worldOne.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_CONFIRMS_RESTART {
				worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_FIRST_MOVE
				worldOne.Region.Objects = []*WorldOneObject{}
				err := s.storage.updateWorld(world, clientId)
				if err != nil {
					return nil, err
				}
				s.enqueueUpdatesAndSignal(game.Player1, s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status))
				s.enqueueUpdatesAndSignal(game.Player2, s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status))
			} else if worldOne.Status != WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_CONFIRMS_RESTART {
				worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_CONFIRMS_RESTART
				err := s.storage.updateWorld(world, clientId)
				if err != nil {
					return nil, err
				}
				s.enqueueUpdatesAndSignal(game.Player1, s.newWorldOneStatusUpdate(worldOne.Region.Id, WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_CONFIRMS_RESTART))
				s.enqueueUpdatesAndSignal(game.Player2, s.newWorldOneStatusUpdate(worldOne.Region.Id, WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_CONFIRMS_RESTART))
			}
		} else if clientId == game.Player2 {
			if worldOne.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_CONFIRMS_RESTART {
				worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_FIRST_MOVE
				worldOne.Region.Objects = []*WorldOneObject{}
				err := s.storage.updateWorld(world, clientId)
				if err != nil {
					return nil, err
				}
				s.enqueueUpdatesAndSignal(game.Player1, s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status))
				s.enqueueUpdatesAndSignal(game.Player2, s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status))
			} else if worldOne.Status != WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_CONFIRMS_RESTART {
				worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_CONFIRMS_RESTART
				err := s.storage.updateWorld(world, clientId)
				if err != nil {
					return nil, err
				}
				s.enqueueUpdatesAndSignal(game.Player1, s.newWorldOneStatusUpdate(worldOne.Region.Id, WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_CONFIRMS_RESTART))
				s.enqueueUpdatesAndSignal(game.Player2, s.newWorldOneStatusUpdate(worldOne.Region.Id, WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_CONFIRMS_RESTART))
			}
		}
	}
	return s.newRestartWorldReply(), nil
}
