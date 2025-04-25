package main

func (s *server) createAcquireClientIdReply(clientId string) *Reply {
	return &Reply{
		Type: &Reply_AcquireClientIdReply{
			AcquireClientIdReply: &AcquireClientIdReply{ClientId: clientId},
		},
	}
}
