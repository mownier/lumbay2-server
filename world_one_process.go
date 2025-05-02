package main

import "google.golang.org/grpc/codes"

func (s *server) processWorldOneObject(clientId string, worldRegionId WorldRegionId, worldObject *WorldObject) error {
	if !regionBelongsToWorldOne(worldRegionId) {
		return sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	return nil
}

func regionBelongsToWorldOne(regionId WorldRegionId) bool {
	regionIds := []WorldRegionId{
		WorldRegionId_WORLD_ONE_REGION_ONE,
	}
	okay := false
	for _, rId := range regionIds {
		if rId == regionId {
			okay = true
			break
		}
	}
	return okay
}
