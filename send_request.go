package main

import (
	context "context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *server) SendRequest(ctx context.Context, in *Request) (*Reply, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "send request cancelled")

	default:
		publicKey := ""
		clientId := ""
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get("public_key-bin")
			if len(values) > 0 {
				publicKey = values[0]
			}
			values = md.Get("client_id-bin")
			if len(values) > 0 {
				clientId = values[0]
			}
		}
		r, e := s.sendRequestInternal(publicKey, clientId, in)
		return r, e
	}
}

func (s *server) sendRequestInternal(publicKey, clientId string, in *Request) (*Reply, error) {
	switch in.Type.(type) {
	case *Request_AcquireClientIdRequest:
		return s.acquireClientId(publicKey)
	case *Request_AcquirePublicKeyRequest:
		return s.acquirePublicKey(in.GetAcquirePublicKeyRequest())
	case *Request_CreateGameRequest:
		return s.createGame(clientId)
	case *Request_GenerateGameCodeRequest:
		return s.generateGameCode(clientId)
	case *Request_JoinGameRequest:
		return s.joinGame(clientId, in.GetJoinGameRequest().GetGameCode())
	case *Request_QuitGameRequest:
		return s.quitGame(clientId)
	case *Request_StartGameRequest:
		return s.startGame(clientId)
	}

	return nil, status.Error(codes.InvalidArgument, "request not known")
}
