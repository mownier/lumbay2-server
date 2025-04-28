package main

func (s *server) createGame(clientId string) (*Reply, error) {
	_, err := s.storage.insertGame(clientId)
	if err != nil {
		return nil, err
	}
	return s.createGameReply(), nil
}
