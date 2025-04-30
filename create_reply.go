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

func (s *server) newWaitingForOtherPlayerUpdate() isUpdate_Type {
	return &Update_WaitingForOtherPlayerUpdate{
		WaitingForOtherPlayerUpdate: &WaitingForOtherPlayerUpdate{},
	}
}

func (s *server) newReadyToStartUpdate() isUpdate_Type {
	return &Update_ReadyToStartUpdate{
		ReadyToStartUpdate: &ReadyToStartUpdate{},
	}
}

func (s *server) newYouAreInGameUpdate() isUpdate_Type {
	return &Update_YouAreInGameUpdate{
		YouAreInGameUpdate: &YouAreInGameUpdate{},
	}
}

func (s *server) newGameCodeGeneratedUpdate(gameCode string) isUpdate_Type {
	return &Update_GameCodeGenerated{
		GameCodeGenerated: &GameCodeGeneratedUpdate{
			GameCode: gameCode,
		},
	}
}
