package main

func (s *server) quitGame(clientId string) (*Reply, error) {
	game, err := s.storage.quitGame(clientId)
	if err != nil {
		return nil, err
	}
	s.enqueueUpdatesAndSignal(game.Player1, s.newGameStatusUpdate(game.Status))
	s.enqueueUpdatesAndSignal(game.Player2, s.newGameStatusUpdate(game.Status))
	s.enqueueUpdatesAndSignal(clientId, s.newGameStatusUpdate(GameStatus_GAME_STATUS_NONE))
	return s.newQuitGameReply(), nil
}
