package main

type storage interface {
	saveClientId(clientId string) error
	containsClientId(clientId string) (bool, error)
	getAllClientIds() ([]string, error)
	insertGame(player1 string) (*Game, error)
	hasGame(gameId string) (bool, error)
}
