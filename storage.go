package main

type storage interface {
	insertClient(publicKeyPEM string) (*Client, error)
	getClient(id string) (*Client, error)

	insertGame(player1 string) (*Game, error)
	getGame(gameId string) (*Game, error)
	getGameForClient(clientId string) (*Game, error)
	updateGame(game *Game) error

	getAllUpdates(clientId string) ([]*Update, error)
	enqueueUpdate(clientId string, updateType isUpdate_Type) (*Update, error)
	dequeueUpdates(clientId string, updates []*Update) error
}
