package main

func newWorldOne() *World {
	return &World{
		Id:              "world-1",
		Name:            "World 1",
		CurrentRegionId: "world-1-region-1",
		Status: &WorldStatus{
			Type: &WorldStatus_WorldOneStatus{
				WorldOneStatus: WorldOneStatus_WORLD_ONE_STATUS_NOTHING,
			},
		},
		Regions: []*WorldRegion{
			{
				Id:   "world-1-region-1",
				Name: "World 1, Region 1",
				Status: &WorldRegionStatus{
					Type: &WorldRegionStatus_WorldOneRegionStatus{
						WorldOneRegionStatus: WorldOneRegionStatus_WORLD_ONE_REGION_STATUS_NOTHING,
					},
				},
			},
		},
	}
}
