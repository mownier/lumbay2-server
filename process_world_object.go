package main

import "google.golang.org/grpc/codes"

func (s *server) processWorldObject(clientId string, in *ProcessWorldObjectRequest) (*Reply, error) {
	_, err := s.storage.getClient(clientId)
	if err != nil {
		return nil, err
	}
	switch in.WorldId {
	case WorldId_WORLD_ID_WORLD_ONE:
		s.processWorldOneObject(clientId, in.WorldRegionId, in.WorldObject)
	default:
		return nil, sverror(codes.InvalidArgument, "failed to process world object", nil)
	}
	return s.newProcessWorldObjectReply(), nil
}
