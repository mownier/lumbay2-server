package main

import "google.golang.org/grpc/codes"

func (s *server) processWorldOneObject(clientId string, worldRegionId WorldRegionId, worldObject *WorldObject) (*Game, error) {
	if !regionBelongsToWorldOne(worldRegionId) {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	if !objectBelongsToWorldOne(worldObject.Id) {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	world, err := s.storage.getWorldForClient(clientId)
	if err != nil {
		return nil, err
	}
	var worldRegion *WorldRegion
	for _, r := range world.Regions {
		if r.Id == worldRegionId {
			worldRegion = r
			break
		}
	}
	if worldRegion == nil {
		return nil, sverror(codes.Internal, "failed to process world one object", nil)
	}
	game, err := s.storage.getGameForClient(clientId)
	if err != nil {
		return nil, err
	}
	clientIsPlayer1 := false
	clientIsPlayer2 := false
	if game.Player1 == clientId {
		clientIsPlayer1 = true
	}
	if game.Player2 == clientId {
		clientIsPlayer2 = true
	}
	if clientIsPlayer1 && clientIsPlayer2 {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	if !clientIsPlayer1 && !clientIsPlayer2 {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	if clientIsPlayer1 && !objectBelongsToPlayer1InWorldOne(worldObject.Id) {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	if clientIsPlayer2 && !objectBelongsToPlayer2InWorldOne(worldObject.Id) {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	switch worldObject.Status {
	case WorldObjectStatus_WORLD_ONE_OBJECT_STATUS_SPAWNED:
		countStones := func(oId WorldObjectId) int {
			var count = 0
			for _, o := range worldRegion.Objects {
				if o.Id == oId {
					count = count + 1
				}
			}
			return count
		}
		var stoneCount int
		if clientIsPlayer1 {
			stoneCount = countStones(WorldObjectId_WORLD_ONE_OBJECT_STONE_1)
		} else if clientIsPlayer2 {
			stoneCount = countStones(WorldObjectId_WORLD_ONE_OBJECT_STONE_2)
		} else {
			return nil, sverror(codes.Internal, "failed to process world one object", nil)
		}
		const stoneCountLimit = 3
		if stoneCount >= stoneCountLimit {
			return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
		}
		verifyLocation := func(location *WorldLocation) bool {
			if location.X < 0 || location.X >= 3 || location.Y < 0 || location.Y >= 3 {
				return false
			}
			okay := true
			for _, object := range worldRegion.Objects {
				if object.Location.X == location.X && object.Location.Y == location.Y {
					okay = false
					break
				}
			}
			return okay
		}
		if !verifyLocation(worldObject.Location) {
			return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
		}
		worldRegion.Objects = append(worldRegion.Objects, worldObject)
		err = s.storage.updateWorld(world, clientId)
		if err != nil {
			return nil, err
		}
	case WorldObjectStatus_WORLD_ONE_OBJECT_STATUS_MOVED:
		verifyLocation := func(location *WorldLocation) bool {
			if location.X < 0 || location.X >= 3 || location.Y < 0 || location.Y >= 3 {
				return false
			}
			okay := true
			for _, object := range worldRegion.Objects {
				if object.Location.X == location.X && object.Location.Y == location.Y {
					okay = false
					break
				}
			}
			return okay
		}
		if !verifyLocation(worldObject.Location) {
			return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
		}
		worldRegion.Objects = append(worldRegion.Objects, worldObject)
		err = s.storage.updateWorld(world, clientId)
		if err != nil {
			return nil, err
		}
	}
	return game, nil
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

func objectBelongsToWorldOne(objectId WorldObjectId) bool {
	objectIds := []WorldObjectId{
		WorldObjectId_WORLD_ONE_OBJECT_STONE_1,
		WorldObjectId_WORLD_ONE_OBJECT_STONE_2,
	}
	okay := false
	for _, oId := range objectIds {
		if oId == objectId {
			okay = true
			break
		}
	}
	return okay
}

func objectBelongsToPlayer1InWorldOne(objectId WorldObjectId) bool {
	objectIds := []WorldObjectId{
		WorldObjectId_WORLD_ONE_OBJECT_STONE_1,
	}
	okay := false
	for _, oId := range objectIds {
		if oId == objectId {
			okay = true
			break
		}
	}
	return okay
}

func objectBelongsToPlayer2InWorldOne(objectId WorldObjectId) bool {
	objectIds := []WorldObjectId{
		WorldObjectId_WORLD_ONE_OBJECT_STONE_2,
	}
	okay := false
	for _, oId := range objectIds {
		if oId == objectId {
			okay = true
			break
		}
	}
	return okay
}
