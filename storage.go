package main

type storage interface {
	insertClient(publicKeyPEM string) (*Client, error)
	getClient(id string) (*Client, error)

	insertGame(player1 string) (*Game, error)
	insertGameCode(gameCode, gameId string) error
	getGame(gameId string) (*Game, error)
	getGameForClient(clientId string) (*Game, error)
	getGameIdForGameCode(gameCode string) (string, error)
	getGameCodeForGame(gameId string) (string, error)
	removeGameCode(gameCode, gameId string) error
	updateGame(game *Game) error
	setGameForClient(gameId, clientId string) error

	getAllUpdates(clientId string) ([]*Update, error)
	enqueueUpdate(clientId string, updateType isUpdate_Type) (*Update, error)
	dequeueUpdates(clientId string, updates []*Update) error
}
