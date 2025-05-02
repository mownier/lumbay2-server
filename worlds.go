package main

func newWorldOne() *World {
	return &World{
		Id:              WorldId_WORLD_ID_WORLD_ONE,
		CurrentRegionId: WorldRegionId_WORLD_ONE_REGION_ONE,
		Status:          WorldStatus_WORLD_STATUS_NOTHING,
		Regions: []*WorldRegion{
			{
				Id:     WorldRegionId_WORLD_ONE_REGION_ONE,
				Status: WorldRegionStatus_WORLD_REGION_STATUS_NOTHING,
			},
		},
	}
}
