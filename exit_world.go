package main

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
		if worldOne.Status == WorldOneStatus_WORLD_ONE_STATUS_ABANDONED {
			break
		}
		if clientId == game.Player1 {
			oldStatus := worldOne.Status
			if worldOne.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_EXITED {
				worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_ABANDONED
			} else if worldOne.Status != WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_EXITED {
				worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_EXITED
			}
			if oldStatus != worldOne.Status {
				err := s.storage.detachWorldFromClient(world, clientId)
				if err != nil {
					return nil, err
				}
			}
			s.enqueueUpdatesAndSignal(game.Player1, s.newYouExitWorldUpdate())
			s.enqueueUpdatesAndSignal(game.Player2, s.newOtherExitsWorldUpdate())
		} else if clientId == game.Player2 {
			oldStatus := worldOne.Status
			if worldOne.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_EXITED {
				worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_ABANDONED
			} else if worldOne.Status != WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_EXITED {
				worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_EXITED
			}
			if oldStatus != worldOne.Status {
				err := s.storage.detachWorldFromClient(world, clientId)
				if err != nil {
					return nil, err
				}
			}
			s.enqueueUpdatesAndSignal(game.Player1, s.newOtherExitsWorldUpdate())
			s.enqueueUpdatesAndSignal(game.Player2, s.newYouExitWorldUpdate())
		}
	}
	return s.newExitWorldReply(), nil
}
