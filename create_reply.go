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
