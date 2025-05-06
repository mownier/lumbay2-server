package main

func (s *server) createAcquireClientIdReply(clientId string) *Reply {
	return &Reply{
		Type: &Reply_AcquireClientIdReply{
			AcquireClientIdReply: &AcquireClientIdReply{ClientId: clientId},
		},
	}
}

func (s *server) createAcquirePublickKeyReply(publicKey string) *Reply {
	return &Reply{
		Type: &Reply_AcquirePublicKeyReply{
			AcquirePublicKeyReply: &AcquirePublicKeyReply{PublicKey: publicKey},
		},
	}
}

func (s *server) newCreateGameReply() *Reply {
	return &Reply{
		Type: &Reply_CreateGameReply{
			CreateGameReply: &CreateGameReply{},
		},
	}
}

func (s *server) newGenerateGameCodeReply() *Reply {
	return &Reply{
		Type: &Reply_GenerateGameCodeReply{
			GenerateGameCodeReply: &GenerateGameCodeReply{},
		},
	}
}

func (s *server) newJoinGameReply() *Reply {
	return &Reply{
		Type: &Reply_JoinGameReply{
			JoinGameReply: &JoinGameReply{},
		},
	}
}

func (s *server) newQuitGameReply() *Reply {
	return &Reply{
		Type: &Reply_QuitGameReply{
			QuitGameReply: &QuitGameReply{},
		},
	}
}

func (s *server) newStartGameReply() *Reply {
	return &Reply{
		Type: &Reply_StartGameReply{
			StartGameReply: &StartGameReply{},
		},
	}
}

func (s *server) newProcessWorldOneObjectReply() *Reply {
	return &Reply{
		Type: &Reply_ProcessWorldOneObjectReply{
			ProcessWorldOneObjectReply: &ProcessWorldOneObjectReply{},
		},
	}
}

func (s *server) newExitWorldReply() *Reply {
	return &Reply{
		Type: &Reply_ExitWorldReply{
			ExitWorldReply: &ExitWorldReply{},
		},
	}
}

func (s *server) newRestartWorldReply() *Reply {
	return &Reply{
		Type: &Reply_RestartWorldReply{
			RestartWorldReply: &RestartWorldReply{},
		},
	}
}

func (s *server) newGameCodeGeneratedUpdate(gameCode string) isUpdate_Type {
	return &Update_GameCodeGenerated{
		GameCodeGenerated: &GameCodeGeneratedUpdate{
			GameCode: gameCode,
		},
	}
}

func (s *server) newGameStatusUpdate(gameStatus GameStatus) isUpdate_Type {
	return &Update_GameStatusUpdate{
		GameStatusUpdate: &GameStatusUpdate{Status: gameStatus},
	}
}

func (s *server) newWorldOneRegionUpdate(regionId WorldOneRegionId) isUpdate_Type {
	return &Update_WorldOneRegionUpdate{
		WorldOneRegionUpdate: &WorldOneRegionUpdate{RegionId: regionId},
	}
}

func (s *server) newWorldOneStatusUpdate(regionId WorldOneRegionId, status WorldOneStatus) isUpdate_Type {
	return &Update_WorldOneStatusUpdate{
		WorldOneStatusUpdate: &WorldOneStatusUpdate{RegionId: regionId, Status: status},
	}
}

func (s *server) newWorldOneObjectUpdate(in *ProcessWorldOneObjectRequest) isUpdate_Type {
	return &Update_WorldOneObjectUpdate{
		WorldOneObjectUpdate: &WorldOneObjectUpdate{
			RegionId:     in.RegionId,
			ObjectId:     in.ObjectId,
			ObjectStatus: in.ObjectStatus,
			ObjectData:   in.ObjectData,
		},
	}
}
