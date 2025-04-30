package main

func (s *server) generateGameCode(clientId string) (*Reply, error) {
	game, err := s.storage.getGameForClient(clientId)
	if err != nil {
		return nil, err
	}
	currentGameCode, err := s.storage.getGameCodeForGame(game.Id)
	if err != nil {
		return nil, err
	}
	if len(currentGameCode) > 0 {
		err := s.storage.removeGameCode(currentGameCode, game.Id)
		if err != nil {
			return nil, err
		}
	}
	newGameCode := generateGameCode()
	err = s.storage.insertGameCode(newGameCode, game.Id)
	if err != nil {
		return nil, err
	}
	s.enqueueUpdatesAndSignal(clientId, s.newGameCodeGeneratedUpdate(newGameCode))
	return s.newGenerateGameCodeReply(), nil
}
