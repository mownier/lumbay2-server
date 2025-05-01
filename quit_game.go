package main

func (s *server) quitGame(clientId string) (*Reply, error) {
	game, err := s.storage.quitGame(clientId)
	if err != nil {
		return nil, err
	}
	switch game.Status {
	case GameStatus_WAITING_FOR_OTHER_PLAYER:
		if len(game.Player1) > 0 {
			s.enqueueUpdatesAndSignal(game.Player1, s.newWaitingForOtherPlayerUpdate())
		}
		if len(game.Player2) > 0 {
			s.enqueueUpdatesAndSignal(game.Player2, s.newWaitingForOtherPlayerUpdate())
		}
	}
	s.enqueueUpdatesAndSignal(clientId, s.newYouQuitTheGameUpdate())
	return s.newQuitGameReply(), nil
}
