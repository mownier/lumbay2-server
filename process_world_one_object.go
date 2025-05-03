package main

import "google.golang.org/grpc/codes"

func (s *server) processWorldOneObject(clientId string, in *ProcessWorldOneObjectRequest) (*Reply, error) {
	_, err := s.storage.getClient(clientId)
	if err != nil {
		return nil, err
	}
	world, err := s.storage.getWorldForClient(clientId)
	if err != nil {
		return nil, err
	}
	var worldOne *WorldOne
	switch world.Type.(type) {
	case *World_WorldOne:
		worldOne = world.GetWorldOne()
	default:
		return nil, sverror(codes.Internal, "failed to process world one object", nil)
	}
	if worldOne.GetRegion().Id != in.RegionId {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
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
	if clientIsPlayer1 && in.ObjectId != WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_ONE {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	if clientIsPlayer2 && in.ObjectId != WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_TWO {
		return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
	}
	switch in.ObjectStatus {
	case WorldOneObjectStatus_WORLD_ONE_OBJECT_STATUS_SPAWNED:
		countStones := func(oId WorldOneObjectId) int {
			var count = 0
			for _, o := range worldOne.Region.Objects {
				if o.Id == oId {
					count = count + 1
				}
			}
			return count
		}
		var stoneCount int
		if clientIsPlayer1 {
			stoneCount = countStones(WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_ONE)
		} else if clientIsPlayer2 {
			stoneCount = countStones(WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_TWO)
		} else {
			return nil, sverror(codes.Internal, "failed to process world one object", nil)
		}
		const stoneCountLimit = 3
		if stoneCount >= stoneCountLimit {
			return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
		}
		verifyLocation := func(location *WorldLocation) bool {
			if location == nil || location.X < 0 || location.X >= 3 || location.Y < 0 || location.Y >= 3 {
				return false
			}
			okay := true
			for _, object := range worldOne.Region.Objects {
				var loc *WorldLocation
				switch object.Data.Type.(type) {
				case *WorldOneObjectData_Location:
					loc = object.Data.GetLocation()
				default:
					loc = nil
				}
				if loc == nil {
					continue
				}
				if loc.X == location.X && loc.Y == location.Y {
					okay = false
					break
				}
			}
			return okay
		}
		if !verifyLocation(in.ObjectData.GetLocation()) {
			return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
		}
		worldObject := &WorldOneObject{
			Id:     in.ObjectId,
			Status: in.ObjectStatus,
			Data:   in.ObjectData,
		}
		worldOne.Region.Objects = append(worldOne.Region.Objects, worldObject)
		err = s.storage.updateWorld(world, clientId)
		if err != nil {
			return nil, err
		}
	case WorldOneObjectStatus_WORLD_ONE_OBJECT_STATUS_MOVED:
		verifyLocation := func(location *WorldLocation) bool {
			if location == nil || location.X < 0 || location.X >= 3 || location.Y < 0 || location.Y >= 3 {
				return false
			}
			okay := true
			for _, object := range worldOne.Region.Objects {
				var loc *WorldLocation
				switch object.Data.Type.(type) {
				case *WorldOneObjectData_Location:
					loc = object.Data.GetLocation()
				default:
					loc = nil
				}
				if loc == nil {
					continue
				}
				if loc.X == location.X && loc.Y == location.Y {
					okay = false
					break
				}
			}
			return okay
		}
		if !verifyLocation(in.ObjectData.GetLocation()) {
			return nil, sverror(codes.InvalidArgument, "failed to process world one object", nil)
		}
		worldObject := &WorldOneObject{
			Id:     in.ObjectId,
			Status: in.ObjectStatus,
			Data:   in.ObjectData,
		}
		worldOne.Region.Objects = append(worldOne.Region.Objects, worldObject)
		err = s.storage.updateWorld(world, clientId)
		if err != nil {
			return nil, err
		}
	}
	s.enqueueUpdatesAndSignal(game.Player1, s.newWorldOneObjectUpdate(in))
	s.enqueueUpdatesAndSignal(game.Player2, s.newWorldOneObjectUpdate(in))
	return s.newProcessWorldOneObjectReply(), nil
}
