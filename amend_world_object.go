package main

func (s *server) amendWorldObject(clientId string, in *AmendWorldObjectRequest) (*Reply, error) {
	_, err := s.storage.getClient(clientId)
	if err != nil {
		return nil, err
	}
	return s.newAmendWorldObjectReply(), nil
}
