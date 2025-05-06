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
		oldStatus := worldOne.Status
		youArePlayer1 := clientId == game.Player1
		youArePlayer2 := clientId == game.Player2
		player1Exited := worldOne.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_EXITED
		player2Exited := worldOne.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_EXITED
		noOneYetExited := worldOne.Status != WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_EXITED && worldOne.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_EXITED
		if (youArePlayer1 && player2Exited) || (youArePlayer2 && player1Exited) {
			worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_ABANDONED
		} else if youArePlayer1 && noOneYetExited {
			worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_EXITED
		} else if youArePlayer2 && noOneYetExited {
			worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_EXITED
		}
		if oldStatus != worldOne.Status {
			updatedGame, err := s.storage.detachWorldFromClient(world, clientId)
			if err != nil {
				return nil, err
			}
			s.enqueueUpdatesAndSignal(game.Player1,
				s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status),
				s.newGameStatusUpdate(updatedGame.Status),
			)
			s.enqueueUpdatesAndSignal(game.Player2,
				s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status),
				s.newGameStatusUpdate(updatedGame.Status),
			)
		}
	}
	return s.newExitWorldReply(), nil
}
