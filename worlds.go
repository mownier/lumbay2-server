package main

import (
	"fmt"

	"github.com/google/uuid"
)

func newWorldOne(player1, player2 string) *World {
	return &World{
		Id:   uuid.New().String(),
		Name: "World 1",
		Status: &WorldStatus{
			Type: &WorldStatus_WorldOneStatus{
				WorldOneStatus: WorldOneStatus_WORLD_ONE_STATUS_NOTHING,
			},
		},
		Map: &WorldMap{
			Id:     uuid.New().String(),
			Height: 3,
			Width:  3,
		},
		Region: []*WorldRegion{
			{
				Id:   uuid.New().String(),
				Name: "World 1, Region 1",
				Status: &WorldStatus{
					Type: &WorldStatus_WorldOneStatus{
						WorldOneStatus: WorldOneStatus_WORLD_ONE_STATUS_NOTHING,
					},
				},
				Map: &WorldMap{
					Id:     uuid.New().String(),
					Height: 3,
					Width:  3,
				},
				Objects: []*WorldObject{
					{
						Type: &WorldObject_WorldOneObject{
							WorldOneObject: &WorldOneObject{
								Id: fmt.Sprintf("%s:1", player1),
							},
						},
						Location: &WorldMapLocation{X: -1, Y: -1},
					},
					{
						Type: &WorldObject_WorldOneObject{
							WorldOneObject: &WorldOneObject{
								Id: fmt.Sprintf("%s:2", player1),
							},
						},
						Location: &WorldMapLocation{X: -1, Y: -1},
					},
					{
						Type: &WorldObject_WorldOneObject{
							WorldOneObject: &WorldOneObject{
								Id: fmt.Sprintf("%s:3", player1),
							},
						},
						Location: &WorldMapLocation{X: -1, Y: -1},
					},
					{
						Type: &WorldObject_WorldOneObject{
							WorldOneObject: &WorldOneObject{
								Id: fmt.Sprintf("%s:1", player2),
							},
						},
						Location: &WorldMapLocation{X: -1, Y: -1},
					},
					{
						Type: &WorldObject_WorldOneObject{
							WorldOneObject: &WorldOneObject{
								Id: fmt.Sprintf("%s:2", player2),
							},
						},
						Location: &WorldMapLocation{X: -1, Y: -1},
					},
					{
						Type: &WorldObject_WorldOneObject{
							WorldOneObject: &WorldOneObject{
								Id: fmt.Sprintf("%s:3", player2),
							},
						},
						Location: &WorldMapLocation{X: -1, Y: -1},
					},
				},
			},
		},
	}
}
