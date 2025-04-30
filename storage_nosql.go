package main

import (
	"fmt"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

const (
	clientPrefix           = "client:"
	gamePrefix             = "game:"
	gameClientPrefix       = "game:client:"
	clientLastSeqNumPrefix = "client:last_seq_num:"
	clientUpdatePrefix     = "client:update:"
	gameCodePrefix         = "game_code:"
	gameGameCodePrefix     = "game:game_code:"
)

type storageNoSql struct {
	db *badger.DB
}

func newStorageNoSql(db *badger.DB) *storageNoSql {
	return &storageNoSql{db: db}
}

func (s *storageNoSql) insertClient(publicKeyPEM string) (*Client, error) {
	salt := generateClientSalt()
	client := &Client{
		Id:   generateClientId(publicKeyPEM, salt),
		Salt: salt,
	}
	clientData, err := proto.Marshal(client)
	if err != nil {
		return nil, sverror(codes.Internal, "failed to insert client", err)
	}
	key := fmt.Sprintf("%s%s", clientPrefix, client.Id)
	err = s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), clientData)
		if err != nil {
			return sverror(codes.Internal, "failed to insert client", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *storageNoSql) getClient(id string) (*Client, error) {
	var client *Client
	key := fmt.Sprintf("%s%s", clientPrefix, id)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return sverror(codes.NotFound, "failed to get client", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to get client", err)
		}
		c := &Client{}
		err = proto.Unmarshal(bytes, c)
		if err != nil {
			return sverror(codes.Internal, "failed to get client", err)
		}
		client = c
		return nil
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *storageNoSql) insertGame(player1 string) (*Game, error) {
	game := &Game{
		Id:      uuid.New().String(),
		Player1: player1,
		Player2: "",
		Status:  GameStatus_WAITING_FOR_OTHER_PLAYER,
	}
	gameData, err := proto.Marshal(game)
	if err != nil {
		return nil, sverror(codes.Internal, "failed to insert game", err)
	}
	gameKey := fmt.Sprintf("%s%s", gamePrefix, game.Id)
	gameClientKey := fmt.Sprintf("%s%s", gameClientPrefix, game.Player1)
	err = s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(gameKey), gameData)
		if err != nil {
			return sverror(codes.Internal, "failed to insert game", err)
		}
		err = txn.Set([]byte(gameClientKey), []byte(game.Id))
		if err != nil {
			return sverror(codes.Internal, "failed to insert game", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *storageNoSql) getGame(id string) (*Game, error) {
	var game *Game
	key := fmt.Sprintf("%s%s", gamePrefix, id)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return sverror(codes.NotFound, "failed to get game", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to get game", err)
		}
		g := &Game{}
		err = proto.Unmarshal(bytes, g)
		if err != nil {
			return sverror(codes.Internal, "failed to get game", err)
		}
		game = g
		return nil
	})
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *storageNoSql) getGameForClient(clientId string) (*Game, error) {
	var game *Game
	key := fmt.Sprintf("%s%s", gameClientPrefix, clientId)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return sverror(codes.NotFound, "failed to get game for client", err)
		}
		gameIdData, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to get game for client", err)
		}
		gameId := string(gameIdData)
		gameKey := fmt.Sprintf("%s%s", gamePrefix, gameId)
		item, err = txn.Get([]byte(gameKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to get game for client", err)
		}
		gameData, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to get game for client", err)
		}
		g := &Game{}
		err = proto.Unmarshal(gameData, g)
		if err != nil {
			return sverror(codes.Internal, "failed to get game for client", err)
		}
		game = g
		return nil
	})
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *storageNoSql) getGameIdForGameCode(gameCode string) (string, error) {
	gameId := ""
	key := fmt.Sprintf("%s%s", gameCodePrefix, gameCode)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return sverror(codes.Internal, "failed to get game id for game code", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to get game id for game code", err)
		}
		gameId = string(bytes)
		return nil
	})
	if err != nil {
		return gameId, err
	}
	return gameId, nil
}

func (s *storageNoSql) getGameCodeForGame(gameId string) (string, error) {
	gameCode := ""
	key := fmt.Sprintf("%s%s", gameGameCodePrefix, gameId)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return sverror(codes.Internal, "failed to get game code for game", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to get game code for game", err)
		}
		gameCode = string(bytes)
		return nil
	})
	if err != nil {
		return gameCode, err
	}
	return gameCode, nil
}

func (s *storageNoSql) insertGameCode(clientId string) (string, error) {
	gameCode := ""
	err := s.db.Update(func(txn *badger.Txn) error {
		gameClientKey := fmt.Sprintf("%s%s", gameClientPrefix, clientId)
		item, err := txn.Get([]byte(gameClientKey))
		if err != nil {
			return sverror(codes.Internal, "failed to insert game code", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to insert game code", err)
		}
		gameId := string(bytes)
		gameGameCodeKey := fmt.Sprintf("%s%s", gameGameCodePrefix, gameId)
		item, err = txn.Get([]byte(gameGameCodeKey))
		if err != nil {
			return sverror(codes.Internal, "failed to insert game code", err)
		}
		bytes, err = item.ValueCopy(nil)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return sverror(codes.Internal, "failed to insert game code", err)
			}
		} else {
			currentGameCode := string(bytes)
			gameCodeKey := fmt.Sprintf("%s%s", gameCodePrefix, currentGameCode)
			// ignore error on delete
			txn.Delete([]byte(gameGameCodeKey))
			txn.Delete([]byte(gameCodeKey))
		}
		newGameCode := generateGameCode()
		gameCodeKey := fmt.Sprintf("%s%s", gameCodePrefix, newGameCode)
		err = txn.Set([]byte(gameCodeKey), []byte(gameId))
		if err != nil {
			return sverror(codes.Internal, "failed to insert game code", err)
		}
		gameCode = newGameCode
		return nil
	})
	return gameCode, err
}

func (s *storageNoSql) joinGame(clientId, gameCode string) (*Game, error) {
	var updatedGame *Game
	gameCodeKey := fmt.Sprintf("%s%s", gameCodePrefix, gameCode)
	err := s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(gameCodeKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to join game", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to join game", err)
		}
		gameId := string(bytes)
		gameKey := fmt.Sprintf("%s%s", gamePrefix, gameId)
		item, err = txn.Get([]byte(gameKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to join game", err)
		}
		bytes, err = item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to join game", err)
		}
		game := &Game{}
		err = proto.Unmarshal(bytes, game)
		if err != nil {
			return sverror(codes.Internal, "failed to join game", err)
		}
		if len(game.Player1) > 0 && len(game.Player2) > 0 {
			return sverror(codes.Internal, "failed to join game because player count limit is reached", nil)
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
		updatedGameData, err := proto.Marshal(game)
		if err != nil {
			return sverror(codes.Internal, "failed to join game", err)
		}
		gameClientKey := fmt.Sprintf("%s%s", gameClientPrefix, clientId)
		err = txn.Set([]byte(gameKey), updatedGameData)
		if err != nil {
			return sverror(codes.Internal, "failed to join game", err)
		}
		err = txn.Set([]byte(gameClientKey), []byte(game.Id))
		if err != nil {
			return sverror(codes.Internal, "failed to join game", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedGame, nil
}

func (s *storageNoSql) enqueueUpdate(clientId string, updateType isUpdate_Type) (*Update, error) {
	lastSeqNumKey := fmt.Sprintf("%s%s", clientLastSeqNumPrefix, clientId)
	var update *Update
	err := s.db.Update(func(txn *badger.Txn) error {
		var lastSeqNum int64 = 0
		item, err := txn.Get([]byte(lastSeqNumKey))
		if err != nil && err != badger.ErrKeyNotFound {
			return sverror(codes.NotFound, "failed to enqueue update", err)
		}
		if item != nil {
			bytes, err := item.ValueCopy(nil)
			if err != nil {
				return sverror(codes.Internal, "failed to enqueue update", err)
			}
			lastSeqNum, err = strconv.ParseInt(string(bytes), 10, 64)
			if err != nil {
				return sverror(codes.Internal, "failed to enqueue update", err)
			}
		}
		nextSeqNum := lastSeqNum + 1
		u := &Update{SequenceNumber: nextSeqNum, Type: updateType}
		updateData, err := proto.Marshal(u)
		if err != nil {
			return sverror(codes.Internal, "failed to enqueue update", err)
		}
		updatesKey := fmt.Sprintf("%s%s:%d", clientUpdatePrefix, clientId, nextSeqNum)
		err = txn.Set([]byte(updatesKey), updateData)
		if err != nil {
			return sverror(codes.Internal, "failed to enqueue update", err)
		}
		update = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return update, nil
}

func (s *storageNoSql) getAllUpdates(clientId string) ([]*Update, error) {
	list := []*Update{}
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(clientUpdatePrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			bytes, err := item.ValueCopy(nil)
			if err != nil {
				return sverror(codes.Internal, "failed to get all updates", err)
			}
			update := &Update{}
			if err := proto.Unmarshal(bytes, update); err != nil {
				return sverror(codes.Internal, "failed to get all updates", err)
			}
			list = append(list, update)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *storageNoSql) dequeueUpdates(clientId string, updates []*Update) error {
	if len(updates) == 0 {
		return nil
	}
	return s.db.Update(func(txn *badger.Txn) error {
		for _, update := range updates {
			key := fmt.Sprintf("%s%s:%d", clientUpdatePrefix, clientId, update.SequenceNumber)
			if err := txn.Delete([]byte(key)); err != nil {
				return sverror(codes.Internal, "failed to dequeue updates", err)
			}
		}
		return nil
	})
}
