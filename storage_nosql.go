package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

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
	worldPrefix            = "world:"
	wordlClientPrefix      = "world:client:"
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
	log.Printf("get game for client %s\n", clientId)
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
			if err != badger.ErrKeyNotFound {
				return sverror(codes.Internal, "failed to insert game code", err)
			}
		} else {
			bytes, err := item.ValueCopy(nil)
			if err != nil {
				return sverror(codes.Internal, "failed to insert game code", err)
			}
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
		err = txn.Set([]byte(gameGameCodeKey), []byte(newGameCode))
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
			return sverror(codes.NotFound, "1 failed to join game", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "2 failed to join game", err)
		}
		gameId := string(bytes)
		gameKey := fmt.Sprintf("%s%s", gamePrefix, gameId)
		item, err = txn.Get([]byte(gameKey))
		if err != nil {
			return sverror(codes.NotFound, "3 failed to join game", err)
		}
		bytes, err = item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "4 failed to join game", err)
		}
		game := &Game{}
		err = proto.Unmarshal(bytes, game)
		if err != nil {
			return sverror(codes.Internal, "5 failed to join game", err)
		}
		if len(game.Player1) > 0 && len(game.Player2) > 0 {
			return sverror(codes.Internal, "6 failed to join game because player count limit is reached", nil)
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
			return sverror(codes.Internal, "7 failed to join game", err)
		}
		gameClientKey := fmt.Sprintf("%s%s", gameClientPrefix, clientId)
		err = txn.Set([]byte(gameKey), updatedGameData)
		if err != nil {
			return sverror(codes.Internal, "8 failed to join game", err)
		}
		err = txn.Set([]byte(gameClientKey), []byte(game.Id))
		if err != nil {
			return sverror(codes.Internal, "9 failed to join game", err)
		}
		updatedGame = game
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedGame, nil
}

func (s *storageNoSql) quitGame(clientId string) (*Game, error) {
	var updatedGame *Game
	gameClientKey := fmt.Sprintf("%s%s", gameClientPrefix, clientId)
	err := s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(gameClientKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to quit game", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to quit game", err)
		}
		gameId := string(bytes)
		gameKey := fmt.Sprintf("%s%s", gamePrefix, gameId)
		item, err = txn.Get([]byte(gameKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to quit game", err)
		}
		bytes, err = item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to quit game", err)
		}
		game := &Game{}
		err = proto.Unmarshal(bytes, game)
		if err != nil {
			return sverror(codes.Internal, "failed to quit game", err)
		}
		if game.Player1 == clientId {
			game.Player1 = ""
		} else if game.Player2 == clientId {
			game.Player2 = ""
		} else {
			return sverror(codes.Internal, "failed to quit game because you are not part of the game", nil)
		}
		if len(game.Player1) == 0 && len(game.Player2) == 0 {
			game.Status = GameStatus_NONE
		} else {
			game.Status = GameStatus_WAITING_FOR_OTHER_PLAYER
		}
		switch game.Status {
		case GameStatus_NONE:
			err := txn.Delete([]byte(gameKey))
			if err != nil {
				return sverror(codes.Internal, "failed to quit game", err)
			}
		default:
			updatedGameData, err := proto.Marshal(game)
			if err != nil {
				return sverror(codes.Internal, "failed to quit game", err)
			}
			err = txn.Set([]byte(gameKey), updatedGameData)
			if err != nil {
				return sverror(codes.Internal, "failed to quit game", err)
			}
		}
		err = txn.Delete([]byte(gameClientKey))
		if err != nil {
			return sverror(codes.Internal, "failed to quit game", err)
		}
		updatedGame = game
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedGame, nil
}

func (s *storageNoSql) startGame(clientId string) (*Game, bool, error) {
	var game *Game
	var startAlreadyInitiated bool
	gameClientKey := fmt.Sprintf("%s%s", gameClientPrefix, clientId)
	err := s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(gameClientKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to start game", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to start game", err)
		}
		gameId := string(bytes)
		gameKey := fmt.Sprintf("%s%s", gamePrefix, gameId)
		item, err = txn.Get([]byte(gameKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to start game", err)
		}
		bytes, err = item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to start game", err)
		}
		g := &Game{}
		err = proto.Unmarshal(bytes, g)
		if err != nil {
			return sverror(codes.Internal, "failed to start game", err)
		}
		var initiatedAlready bool
		if g.Status == GameStatus_STARTED {
			initiatedAlready = true
		} else {
			g.Status = GameStatus_STARTED
			initiatedAlready = false
		}
		gData, err := proto.Marshal(g)
		if err != nil {
			return sverror(codes.Internal, "failed to start game", err)
		}
		err = txn.Set([]byte(gameKey), gData)
		if err != nil {
			return sverror(codes.Internal, "failed to start game", err)
		}
		game = g
		startAlreadyInitiated = initiatedAlready
		return nil
	})
	if err != nil {
		return game, startAlreadyInitiated, err
	}
	return game, startAlreadyInitiated, nil
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
			seqNum, err := strconv.ParseInt(string(bytes), 10, 64)
			if err != nil {
				return sverror(codes.Internal, "failed to enqueue update", err)
			}
			lastSeqNum = seqNum
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
	lastSeqNumKey := fmt.Sprintf("%s%s", clientLastSeqNumPrefix, clientId)
	list := []*Update{}
	err := s.db.View(func(txn *badger.Txn) error {
		var lastSeqNum int64 = 0
		item, err := txn.Get([]byte(lastSeqNumKey))
		if err == nil && item != nil {
			bytes, err := item.ValueCopy(nil)
			if err != nil {
				return sverror(codes.Internal, "failed to get all updates", err)
			}
			num, err := strconv.ParseInt(string(bytes), 10, 64)
			if err != nil {
				return sverror(codes.Internal, "failed to get all updates", err)
			}
			lastSeqNum = num
		}
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(clientUpdatePrefix + clientId)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			keyBytes := it.Item().KeyCopy(nil)
			keyString := string(keyBytes)
			suffix := strings.TrimPrefix(keyString, string(clientUpdatePrefix))
			suffixParts := strings.Split(suffix, ":")
			if len(suffixParts) != 2 {
				return sverror(codes.Internal, "failed to get all updates, malformed key", nil)
			}
			seqNum, err := strconv.ParseInt(suffixParts[1], 10, 64)
			if err != nil {
				return sverror(codes.Internal, "failed to get all updates", err)
			}
			if seqNum <= lastSeqNum {
				continue
			}
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
	sort.Sort(BySeqNum(list))
	return list, nil
}

func (s *storageNoSql) dequeueUpdates(clientId string, updates []*Update) error {
	if len(updates) == 0 {
		return nil
	}
	sort.Sort(BySeqNum(updates))
	return s.db.Update(func(txn *badger.Txn) error {
		for index, update := range updates {
			key := fmt.Sprintf("%s%s:%d", clientUpdatePrefix, clientId, update.SequenceNumber)
			if err := txn.Delete([]byte(key)); err != nil {
				return sverror(codes.Internal, "failed to dequeue updates", err)
			}
			if index == len(updates)-1 {
				lastSeqNumKey := fmt.Sprintf("%s%s", clientLastSeqNumPrefix, clientId)
				err := txn.Set([]byte(lastSeqNumKey), []byte(fmt.Sprintf("%d", update.SequenceNumber)))
				if err != nil {
					return sverror(codes.Internal, "failed to dequeue updates", err)
				}
			}
		}
		return nil
	})
}

func (s *storageNoSql) insertWorld(world *World, clientIds []string) error {
	worldData, err := proto.Marshal(world)
	if err != nil {
		return sverror(codes.Internal, "failed to insert world", err)
	}
	worldKey := fmt.Sprintf("%s%s", worldPrefix, world.Id)
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(worldKey), worldData)
		if err != nil {
			return sverror(codes.Internal, "failed to insert world", err)
		}
		for _, clientId := range clientIds {
			worldClientKey := fmt.Sprintf("%s%s", wordlClientPrefix, clientId)
			err := txn.Set([]byte(worldClientKey), []byte(world.Id))
			if err != nil {
				return sverror(codes.Internal, "failed to insert world", err)
			}
		}
		return nil
	})
}

func (s *storageNoSql) getWorld(worldId string, clientId string) (*World, error) {
	var world *World
	worldKey := fmt.Sprintf("%s%s", worldPrefix, worldId)
	worldClientKey := fmt.Sprintf("%s%s", wordlClientPrefix, clientId)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(worldClientKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to get world", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to get world", err)
		}
		if worldId != string(bytes) {
			return sverror(codes.InvalidArgument, "failed to get world because you do not belong to it", nil)
		}
		item, err = txn.Get([]byte(worldKey))
		if err != nil {
			return sverror(codes.Internal, "failed to get world", err)
		}
		bytes, err = item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to get world", err)
		}
		w := &World{}
		err = proto.Unmarshal(bytes, w)
		if err != nil {
			return sverror(codes.Internal, "failed to get world", err)
		}
		world = w
		return nil
	})
	if err != nil {
		return nil, err
	}
	return world, nil
}

func (s *storageNoSql) updateWorld(world *World, clientId string) error {
	worldData, err := proto.Marshal(world)
	if err != nil {
		return sverror(codes.Internal, "failed to update world", err)
	}
	worldKey := fmt.Sprintf("%s%s", worldPrefix, world.Id)
	worldClientKey := fmt.Sprintf("%s%s", wordlClientPrefix, clientId)
	return s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(worldClientKey))
		if err != nil {
			return sverror(codes.NotFound, "failed to update world", err)
		}
		bytes, err := item.ValueCopy(nil)
		if err != nil {
			return sverror(codes.Internal, "failed to update world", err)
		}
		if world.Id != string(bytes) {
			return sverror(codes.InvalidArgument, "failed to update world because you do not belong to it", nil)
		}
		err = txn.Set([]byte(worldKey), worldData)
		if err != nil {
			return sverror(codes.Internal, "failed to update world", err)
		}
		return nil
	})
}
