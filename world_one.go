package main

import "errors"

func newWorldOne() *World {
	return &World{
		Id: WorldId_WORLD_ID_ONE,
		Type: &World_WorldOne{
			WorldOne: newWorldTypeWorldOne(),
		},
	}
}

func newWorldTypeWorldOne() *WorldOne {
	region := newWorldOneRegionOne()
	return &WorldOne{
		Status:  WorldOneStatus_WORLD_ONE_STATUS_NONE,
		Region:  region,
		Regions: []*WorldOneRegion{region},
	}
}

func newWorldOneRegionOne() *WorldOneRegion {
	return &WorldOneRegion{
		Id:      WorldOneRegionId_WORLD_ONE_REGION_ID_ONE,
		Objects: []*WorldOneObject{},
	}
}

func (w *WorldOne) needToMove() bool {
	switch w.Status {
	case WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_MOVED,
		WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_MOVED:
		return true
	default:
		return false
	}
}

func (w *WorldOne) needToDetermineFirstMover() bool {
	return w.Status == WorldOneStatus_WORLD_ONE_STATUS_NONE
}

func (r *WorldOneRegion) playerOneStoneCount() int {
	count := 0
	for _, object := range r.Objects {
		if object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_ONE ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_TWO ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_THREE {
			count += 1
		}
	}
	return count
}

func (r *WorldOneRegion) playerTwoStoneCount() int {
	count := 0
	for _, object := range r.Objects {
		if object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_ONE ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_TWO ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_THREE {
			count += 1
		}
	}
	return count
}

func playerOneDoesOwnThisWorldOneObjectId(id WorldOneObjectId) bool {
	if id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_ONE ||
		id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_TWO ||
		id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_THREE {
		return true
	}
	return false
}

func playerTwoDoesOwnThisWorldOneObjectId(id WorldOneObjectId) bool {
	if id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_ONE ||
		id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_TWO ||
		id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_THREE {
		return true
	}
	return false
}

func (r *WorldOneRegion) locationIsOccupied(location *WorldLocation) (bool, error) {
	if location == nil {
		return false, errors.New("cannot deterine if location is occupied")
	}
	occupied := false
	for _, object := range r.Objects {
		loc := object.Data.GetLocation()
		if loc == nil {
			continue
		}
		if location.X == loc.X && location.Y == loc.Y {
			occupied = true
			break
		}
	}
	return occupied, nil
}

func (r *WorldOneRegion) getObject(id WorldOneObjectId) (*WorldOneObject, error) {
	for _, object := range r.Objects {
		if object.Id == id {
			return object, nil
		}
	}
	return nil, errors.New("world one object not found")
}
