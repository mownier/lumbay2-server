package main

func (s *server) startGame(clientId string) (*Reply, error) {
	game, startAlreadyInitiated, err := s.storage.startGame(clientId)
	if err != nil {
		return nil, err
	}
	if startAlreadyInitiated {
		return s.newStartGameReply(), nil
	}
	world := newWorldOne()
	err = s.storage.insertWorld(world, game.Player1, game.Player2)
	if err != nil {
		return nil, err
	}
	regionId := world.GetWorldOne().Region.Id
	in1 := &ProcessWorldOneObjectRequest{
		RegionId:     regionId,
		ObjectId:     WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_ONE,
		ObjectStatus: WorldOneObjectStatus_WORLD_ONE_OBJECT_STATUS_ASSIGNED,
		ObjectData:   nil,
	}
	in2 := &ProcessWorldOneObjectRequest{
		RegionId:     regionId,
		ObjectId:     WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_TWO,
		ObjectStatus: WorldOneObjectStatus_WORLD_ONE_OBJECT_STATUS_ASSIGNED,
		ObjectData:   nil,
	}
	s.enqueueUpdatesAndSignal(game.Player1,
		s.newWorldOneRegionUpdate(regionId),
		s.newWorldOneObjectUpdate(in1),
		s.newGameStartedUpdate(),
	)
	s.enqueueUpdatesAndSignal(game.Player2,
		s.newWorldOneRegionUpdate(regionId),
		s.newWorldOneObjectUpdate(in2),
		s.newGameStartedUpdate(),
	)
	return s.newStartGameReply(), nil
}
