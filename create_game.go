package main

func (s *server) createGame(clientId string) (*Reply, error) {
	game, err := s.storage.insertGame(clientId)
	if err != nil {
		return nil, err
	}
	s.enqueueUpdatesAndSignal(clientId, s.newGameStatusUpdate(game.Status))
	return s.newCreateGameReply(), nil
}
