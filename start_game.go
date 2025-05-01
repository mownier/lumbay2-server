package main

func (s *server) startGame(clientId string) (*Reply, error) {
	game, startAlreadyInitiated, err := s.storage.startGame(clientId)
	if err != nil {
		return nil, err
	}
	if !startAlreadyInitiated {
		s.enqueueUpdatesAndSignal(game.Player1, s.newGameStartedUpdate())
		s.enqueueUpdatesAndSignal(game.Player2, s.newGameStartedUpdate())
	}
	return s.newStartGameReply(), nil
}
