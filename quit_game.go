package main

import "log"

func (s *server) quitGame(clientId string) (*Reply, error) {
	game, err := s.storage.quitGame(clientId)
	if err != nil {
		return nil, err
	}
	log.Printf("client %s, quitGame: %v\n", clientId, game)
	switch game.Status {
	case GameStatus_WAITING_FOR_OTHER_PLAYER:
		if len(game.Player1) > 0 {
			log.Printf("newWaitingForOtherPlayerUpdate, player1 %s\n", game.Player1)
			s.enqueueUpdatesAndSignal(game.Player1, s.newWaitingForOtherPlayerUpdate())
		}
		if len(game.Player2) > 0 {
			log.Printf("newWaitingForOtherPlayerUpdate, player2 %s\n", game.Player2)
			s.enqueueUpdatesAndSignal(game.Player2, s.newWaitingForOtherPlayerUpdate())
		}
	}
	s.enqueueUpdatesAndSignal(clientId, s.newYouQuitTheGameUpdate())
	return s.newQuitGameReply(), nil
}
