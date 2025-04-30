package main

import (
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) joinGame(clientId, gameCode string) (*Reply, error) {
	gameId, err := s.storage.getGameIdForGameCode(gameCode)
	if err != nil {
		return nil, err
	}
	game, err := s.storage.getGame(gameId)
	if err != nil {
		return nil, err
	}
	if len(game.Player1) > 0 && len(game.Player2) > 0 {
		log.Printf("client %s wants to join game %s but number of players are already reached\n", clientId, game.Id)
		return nil, status.Error(codes.InvalidArgument, "failed to join game because player count limit is reached")
	}
	if game.Player1 == clientId || game.Player2 == clientId {
		log.Printf("client %s is already part of the game %s", clientId, game.Id)
		return s.newCreateGameReply(), nil
	}
	if len(game.Player1) == 0 {
		game.Player1 = clientId
	} else if len(game.Player2) == 0 {
		game.Player2 = clientId
	}
	if len(game.Player1) > 0 && len(game.Player2) > 0 {
		game.Status = GameStatus_READY_TO_START
	} else {
		game.Status = GameStatus_WAITING_FOR_OTHER_PLAYER
	}
	err = s.storage.updateGame(game)
	if err != nil {
		return nil, err
	}
	err = s.storage.setGameForClient(game.Id, clientId)
	if err != nil {
		return nil, err
	}
	switch game.Status {
	case GameStatus_READY_TO_START:
		s.enqueueUpdatesAndSignal(game.Player1, s.newReadyToStartUpdate())
		s.enqueueUpdatesAndSignal(game.Player2, s.newReadyToStartUpdate())
	case GameStatus_WAITING_FOR_OTHER_PLAYER:
		if len(game.Player1) > 0 {
			s.enqueueUpdatesAndSignal(game.Player1, s.newWaitingForOtherPlayerUpdate())
		}
		if len(game.Player2) > 0 {
			s.enqueueUpdatesAndSignal(game.Player2, s.newWaitingForOtherPlayerUpdate())
		}
	}
	return s.newCreateGameReply(), nil
}
