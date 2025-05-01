package main

func (s *server) joinGame(clientId, gameCode string) (*Reply, error) {
	game, err := s.storage.joinGame(clientId, gameCode)
	if err != nil {
		return nil, err
	}
	switch game.Status {
	case GameStatus_READY_TO_START:
		s.enqueueUpdatesAndSignal(game.Player1, s.newReadyToStartUpdate())
		s.enqueueUpdatesAndSignal(game.Player2, s.newReadyToStartUpdate())
	case GameStatus_WAITING_FOR_OTHER_PLAYER:
		if len(game.Player1) > 0 {
			s.enqueueUpdatesAndSignal(game.Player1, s.newWaitingForOtherPlayerUpdate())
		}
		if len(game.Player2) > 0 {
			s.enqueueUpdatesAndSignal(game.Player2, s.newWaitingForOtherPlayerUpdate())
		}
	}
	return s.newJoinGameReply(), nil
}
