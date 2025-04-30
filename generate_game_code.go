package main

func (s *server) generateGameCode(clientId string) (*Reply, error) {
	game, err := s.storage.getGameForClient(clientId)
	if err != nil {
		return nil, err
	}
	game.GameCode = generateGameCode()
	err = s.storage.updateGame(game)
	if err != nil {
		return nil, err
	}
	s.enqueueUpdatesAndSignal(clientId, s.newGameCodeGeneratedUpdate(game.GameCode))
	return s.newGenerateGameCodeReply(), nil
}
