package main

import (
	"errors"

	"github.com/google/uuid"
)

func newWorldOne() *World {
	return &World{
		DbId: uuid.New().String(),
		Id:   WorldId_WORLD_ID_ONE,
		Type: &World_WorldOne{
			WorldOne: newWorldTypeWorldOne(),
		},
	}
}

func newWorldTypeWorldOne() *WorldOne {
	region := newWorldOneRegionOne()
	return &WorldOne{
		Status:    WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_FIRST_MOVE,
		Region:    region,
		RegionIds: []WorldOneRegionId{region.Id},
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
		WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_MOVED,
		WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_FIRST_MOVE,
		WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_FIRST_MOVE:
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

func (w *WorldOne) whoIsTheWinner() (WorldOneStatus, WorldOneObjectId) {
	switch w.Status {
	case WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_WINS,
		WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_WINS_BY_OUT_OF_MOVES:
		return w.Status, WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_ONE
	case WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_WINS,
		WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_WINS_BY_OUT_OF_MOVES:
		return w.Status, WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_TWO
	}

	objects := []*WorldOneObject{}
	for _, object := range w.Region.Objects {
		if object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_ONE ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_TWO ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_THREE ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_ONE ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_TWO ||
			object.Id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_THREE {
			objects = append(objects, object)
		}
	}
	if len(objects) > 6 {
		return WorldOneStatus_WORLD_ONE_STATUS_NONE, WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE
	}

	winStatuses := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	cells := map[int][2]int64{
		0: {0, 0}, 1: {1, 0}, 2: {2, 0},
		3: {0, 1}, 4: {1, 1}, 5: {2, 1},
		6: {0, 2}, 7: {1, 2}, 8: {2, 2},
	}

	board := [9]WorldOneObjectId{
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
		WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE,
	}

	for _, object := range objects {
		location := object.Data.GetLocation()
		if location == nil || location.X < 0 || location.Y > 2 {
			return WorldOneStatus_WORLD_ONE_STATUS_NONE, WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE
		}
		for boardIndex, cell := range cells {
			if cell[0] == location.X && cell[1] == location.Y {
				board[boardIndex] = object.Id
				break
			}
		}
	}

	for _, winStatus := range winStatuses {
		if worldOneStoneBelongsToPlayerOne(board[winStatus[0]]) &&
			worldOneStoneBelongsToPlayerOne(board[winStatus[1]]) &&
			worldOneStoneBelongsToPlayerOne(board[winStatus[2]]) {
			return WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_WINS, WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_ONE
		}
		if worldOneStoneBelongsToPlayerTwo(board[winStatus[0]]) &&
			worldOneStoneBelongsToPlayerTwo(board[winStatus[1]]) &&
			worldOneStoneBelongsToPlayerTwo(board[winStatus[2]]) {
			return WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_WINS, WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_TWO
		}
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

	if w.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_MOVED {
		playerTwoCanMove := true
		for _, object := range objects {
			if !object.stoneBelongsToPlayerTwo() {
				continue
			}
			playerTwoCanMove = false
			for movementPathIndex, cell := range cells {
				if cell[0] == object.GetData().GetLocation().X &&
					cell[1] == object.GetData().GetLocation().Y {
					movementPath := movementPaths[movementPathIndex]
					for _, boardIndex := range movementPath {
						if board[boardIndex] == WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE {
							playerTwoCanMove = true
							break
						}
					}
					break
				}
			}
			if playerTwoCanMove {
				break
			}
		}
		if !playerTwoCanMove {
			return WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_WINS_BY_OUT_OF_MOVES, WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_ONE
		}
	}

	if w.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_MOVED {
		playerOneCanMove := true
		for _, object := range objects {
			if !object.stoneBelongsToPlayerOne() {
				continue
			}
			playerOneCanMove = false
			for movementPathIndex, cell := range cells {
				if cell[0] == object.GetData().GetLocation().X &&
					cell[1] == object.GetData().GetLocation().Y {
					movementPath := movementPaths[movementPathIndex]
					for _, boardIndex := range movementPath {
						if board[boardIndex] == WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE {
							playerOneCanMove = true
							break
						}
					}
					break
				}
			}
			if playerOneCanMove {
				break
			}
		}
		if !playerOneCanMove {
			return WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_WINS_BY_OUT_OF_MOVES, WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_TWO
		}
	}

	return WorldOneStatus_WORLD_ONE_STATUS_NONE, WorldOneObjectId_WORLD_ONE_OBJECT_ID_NONE
}

func (o *WorldOneObject) stoneBelongsToPlayerOne() bool {
	return worldOneStoneBelongsToPlayerOne(o.Id)
}

func (o *WorldOneObject) stoneBelongsToPlayerTwo() bool {
	return worldOneStoneBelongsToPlayerTwo(o.Id)
}

func worldOneStoneBelongsToPlayerOne(id WorldOneObjectId) bool {
	return id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_ONE ||
		id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_TWO ||
		id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_ONE_STONE_THREE
}

func worldOneStoneBelongsToPlayerTwo(id WorldOneObjectId) bool {
	return id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_ONE ||
		id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_TWO ||
		id == WorldOneObjectId_WORLD_ONE_OBJECT_ID_PLAYER_TWO_STONE_THREE
}

func (w *WorldOne) gameIsOver() bool {
	_, winner := w.whoIsTheWinner()
	if winner == WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_ONE ||
		winner == WorldOneObjectId_WORLD_ONE_OBJECT_ID_STONE_PLAYER_TWO {
		return true
	}
	return false
}

func (w *WorldOne) viableForLife() bool {
	return w.Status != WorldOneStatus_WORLD_ONE_STATUS_ABANDONED &&
		w.Status != WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_EXITED &&
		w.Status != WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_EXITED
}

func (w *WorldOne) restartIsInitiated() bool {
	return w.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_ONE_CONFIRMS_RESTART ||
		w.Status == WorldOneStatus_WORLD_ONE_STATUS_PLAYER_TWO_CONFIRMS_RESTART
}
