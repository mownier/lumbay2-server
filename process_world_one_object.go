package main

import (
	"google.golang.org/grpc/codes"
)

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
		return nil, sverror(codes.Internal, "1 failed to process world one object", nil)
	}
	if worldOne.GetRegion().Id != in.RegionId {
		return nil, sverror(codes.InvalidArgument, "2 failed to process world one object", nil)
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
		return nil, sverror(codes.InvalidArgument, "3 failed to process world one object", nil)
	}
	if !clientIsPlayer1 && !clientIsPlayer2 {
		return nil, sverror(codes.InvalidArgument, "4 failed to process world one object", nil)
	}
	if clientIsPlayer1 && !playerOneDoesOwnThisWorldOneObjectId(in.ObjectId) {
		return nil, sverror(codes.InvalidArgument, "5 failed to process world one object", nil)
	}
	if clientIsPlayer2 && !playerTwoDoesOwnThisWorldOneObjectId(in.ObjectId) {
		return nil, sverror(codes.InvalidArgument, "6 failed to process world one object", nil)
	}
	var worldOneWinner WorldOneObjectId
	switch in.ObjectStatus {
	case WorldOneObjectStatus_WORLD_ONE_OBJECT_STATUS_SPAWNED:
		var stoneCount int
		if clientIsPlayer1 {
			stoneCount = worldOne.Region.playerOneStoneCount()
		} else if clientIsPlayer2 {
			stoneCount = worldOne.Region.playerTwoStoneCount()
		} else {
			return nil, sverror(codes.Internal, "7 failed to process world one object", nil)
		}
		const stoneCountLimit = 3
		if stoneCount >= stoneCountLimit {
			return nil, sverror(codes.InvalidArgument, "8 failed to process world one object", nil)
		}
		occupied, err := worldOne.Region.locationIsOccupied(in.ObjectData.GetLocation())
		if err != nil {
			return nil, sverror(codes.InvalidArgument, "9 failed to process world one object", nil)
		}
		if occupied {
			return nil, sverror(codes.InvalidArgument, "9.1 failed to process world one object", nil)
		}
		worldObject := &WorldOneObject{
			Id:     in.ObjectId,
			Status: in.ObjectStatus,
			Data:   in.ObjectData,
		}
		if clientIsPlayer1 {
			worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_MOVED
		} else if clientIsPlayer2 {
			worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_MOVED
		}
		worldOne.Region.Objects = append(worldOne.Region.Objects, worldObject)
		worldOneStatusForWinner, winner := worldOne.whoIsTheWinner()
		if winner == WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_ONE ||
			winner == WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_TWO {
			worldOne.Status = worldOneStatusForWinner
		}
		worldOneWinner = winner
		err = s.storage.updateWorld(world, clientId)
		if err != nil {
			return nil, err
		}
	case WorldOneObjectStatus_WORLD_ONE_OBJECT_STATUS_MOVED:
		occupied, err := worldOne.Region.locationIsOccupied(in.ObjectData.GetLocation())
		if err != nil {
			return nil, sverror(codes.InvalidArgument, "10 failed to process world one object", err)
		}
		if occupied {
			return nil, sverror(codes.InvalidArgument, "11 failed to process world one object", nil)
		}
		worldObject, err := worldOne.Region.getObject(in.ObjectId)
		if err != nil {
			return nil, sverror(codes.InvalidArgument, "12 failed to process world one object", err)
		}
		if worldObject.Data.GetLocation() == nil {
			return nil, sverror(codes.InvalidArgument, "13 failed to process world one object", nil)
		}

		cells := map[int][2]int64{
			0: {0, 0}, 1: {1, 0}, 2: {2, 0},
			3: {0, 1}, 4: {1, 1}, 5: {2, 1},
			6: {0, 2}, 7: {1, 2}, 8: {2, 2},
		}

		var objectCurrentIndex int = -1
		var objectTargetIndex int = -1
		for cellIndex, cell := range cells {
			if cell[0] == worldObject.Data.GetLocation().X &&
				cell[1] == worldObject.Data.GetLocation().Y {
				objectCurrentIndex = cellIndex
			}
			if cell[0] == in.ObjectData.GetLocation().X &&
				cell[1] == in.ObjectData.GetLocation().Y {
				objectTargetIndex = cellIndex
			}
			if objectCurrentIndex != -1 && objectTargetIndex != -1 {
				break
			}
		}

		if objectCurrentIndex == -1 || objectTargetIndex == -1 {
			return nil, sverror(codes.InvalidArgument, "14 failed to process world one object", nil)
		}

		movementPaths := map[int][]int64{
			0: {1, 3, 4},
			1: {0, 4, 2},
			2: {1, 5, 4},
			3: {0, 4, 6},
			4: {0, 1, 2, 3, 5, 6, 7, 8},
			5: {2, 4, 8},
			6: {3, 4, 7},
			7: {4, 6, 8},
			8: {4, 5, 7},
		}

		for index, movementPath := range movementPaths {
			if index == objectCurrentIndex {
				validMove := false
				for _, targetIndex := range movementPath {
					if targetIndex == int64(objectTargetIndex) {
						validMove = true
						break
					}
				}
				if !validMove {
					return nil, sverror(codes.InvalidArgument, "15 failed to process world one object", nil)
				}
				break
			}
		}

		worldObject.Data.GetLocation().X = in.ObjectData.GetLocation().X
		worldObject.Data.GetLocation().Y = in.ObjectData.GetLocation().Y
		if clientIsPlayer1 {
			worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_MOVED
		} else if clientIsPlayer2 {
			worldOne.Status = WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_MOVED
		}
		worldOneStatusForWinner, winner := worldOne.whoIsTheWinner()
		if winner == WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_ONE ||
			winner == WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_TWO {
			worldOne.Status = worldOneStatusForWinner
		}
		worldOneWinner = winner
		err = s.storage.updateWorld(world, clientId)
		if err != nil {
			return nil, err
		}
	}
	player1Updates := []isUpdate_Type{s.newWorldOneObjectUpdate(in)}
	player2Updates := []isUpdate_Type{s.newWorldOneObjectUpdate(in)}
	if worldOne.needToMove() {
		player1Updates = append(player1Updates, s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status))
		player2Updates = append(player2Updates, s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status))
	}
	if worldOneWinner == WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_ONE ||
		worldOneWinner == WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_TWO {
		player1Updates = append(player1Updates, s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status))
		player2Updates = append(player2Updates, s.newWorldOneStatusUpdate(worldOne.Region.Id, worldOne.Status))
	}
	s.enqueueUpdatesAndSignal(game.Player1, player1Updates...)
	s.enqueueUpdatesAndSignal(game.Player2, player2Updates...)
	return s.newProcessWorldOneObjectReply(), nil
}
