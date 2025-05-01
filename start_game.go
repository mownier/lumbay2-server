package main

func (s *server) startGame(clientId string) (*Reply, error) {
	game, startAlreadyInitiated, err := s.storage.startGame(clientId)
	if err != nil {
		return nil, err
	}
	if startAlreadyInitiated {
		return s.newStartGameReply(), nil
	}
	world := newWorldOne(game.Player1, game.Player2)
	err = s.storage.insertWorld(world, game.Player1, game.Player2)
	if err != nil {
		return nil, err
	}
	s.enqueueUpdatesAndSignal(game.Player1, s.newGameStartedUpdate(), s.newWorldUpdate(world))
	s.enqueueUpdatesAndSignal(game.Player2, s.newGameStartedUpdate(), s.newWorldUpdate(world))
	return s.newStartGameReply(), nil
}
