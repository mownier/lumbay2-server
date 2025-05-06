package main

func (s *server) joinGame(clientId, gameCode string) (*Reply, error) {
	game, err := s.storage.joinGame(clientId, gameCode)
	if err != nil {
		return nil, err
	}
	s.enqueueUpdatesAndSignal(game.Player1, s.newGameStatusUpdate(game.Status))
	s.enqueueUpdatesAndSignal(game.Player2, s.newGameStatusUpdate(game.Status))
	return s.newJoinGameReply(), nil
}
