package main

type storage interface {
	insertClient(publicKeyPEM string) (*Client, error)
	getClient(id string) (*Client, error)

	insertGame(player1 string) (*Game, error)
	insertGameCode(clientId string) (string, error)
	getGame(gameId string) (*Game, error)
	getGameForClient(clientId string) (*Game, error)
	getGameIdForGameCode(gameCode string) (string, error)
	getGameCodeForGame(gameId string) (string, error)
	joinGame(clientId, gameCode string) (*Game, error)
	quitGame(clientId string) (*Game, error)
	startGame(clientId string) (*Game, bool, error)

	getAllUpdates(clientId string) ([]*Update, error)
	enqueueUpdate(clientId string, updateType isUpdate_Type) (*Update, error)
	dequeueUpdates(clientId string, updates []*Update) error
}
