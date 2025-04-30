package main

func (s *server) generateGameCode(clientId string) (*Reply, error) {
	gameCode, err := s.storage.insertGameCode(clientId)
	if err != nil {
		return nil, err
	}
	s.enqueueUpdatesAndSignal(clientId, s.newGameCodeGeneratedUpdate(gameCode))
	return s.newGenerateGameCodeReply(), nil
}
