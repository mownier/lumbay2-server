package main

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
