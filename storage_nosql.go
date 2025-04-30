package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	clientPrefix           = "client:"
	gamePrefix             = "game:"
	gameClientPrefix       = "game:client:"
	clientLastSeqNumPrefix = "client:last_seq_num:"
	clientUpdatePrefix     = "client:update:"
)

type storageNoSql struct {
	db *badger.DB
}

func newStorageNoSql(db *badger.DB) *storageNoSql {
	return &storageNoSql{db: db}
}

func (s *storageNoSql) insertClient(publicKeyPEM string) (*Client, error) {
	client := &Client{Id: generateClientId(publicKeyPEM)}
	clientData, err := proto.Marshal(client)
	if err != nil {
		log.Printf("unable to marshal to be inserted client: %v\n", err)
		return nil, status.Error(codes.Internal, "failed to marshal to be inserted client")
	}
	key := fmt.Sprintf("%s%s", clientPrefix, client.Id)
	err = s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), clientData)
		if err != nil {
			log.Printf("unable to insert client: %v\n", err)
			return status.Error(codes.Internal, "failed to insert client")
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
			log.Printf("client with id %s not found: %v\n", id, err)
			return status.Error(codes.NotFound, "client not found")
		}
		return item.Value(func(val []byte) error {
			c := &Client{}
			err := proto.Unmarshal(val, c)
			if err != nil {
				log.Printf("unable to unmarshal found client with id %s: %v\n", id, err)
				return status.Error(codes.Internal, "failed to unmarshal found client")
			}
			client = c
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *storageNoSql) insertGame(player1 string) (*Game, error) {
	game := &Game{
		Id:       uuid.New().String(),
		Player1:  player1,
		Player2:  "",
		Status:   GameStatus_WAITING_FOR_OTHER_PLAYER,
		GameCode: "",
	}
	gameData, err := proto.Marshal(game)
	if err != nil {
		log.Printf("unable to marshal to be inserted game: %v\n", err)
		return nil, status.Error(codes.Internal, "failed to marshal to be inserted game")
	}
	gameKey := fmt.Sprintf("%s%s", gamePrefix, game.Id)
	gameClientKey := fmt.Sprintf("%s%s", gameClientPrefix, game.Player1)
	err = s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(gameKey), gameData)
		if err != nil {
			log.Printf("unable to set game data %s: %v\n", player1, err)
			return status.Error(codes.Internal, "failed to set game data")
		}
		err = txn.Set([]byte(gameClientKey), []byte(game.Id))
		if err != nil {
			log.Printf("unable to set game-client data: %v\n", err)
			return status.Error(codes.Internal, "failed to set game-client data")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *storageNoSql) getGame(id string) (*Game, error) {
	key := fmt.Sprintf("%s%s", gamePrefix, id)
	var game *Game
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			log.Printf("unable to get game with id %s: %v\n", id, err)
			return status.Error(codes.Internal, "failed to get game")
		}
		return item.Value(func(val []byte) error {
			g := &Game{}
			err := proto.Unmarshal(val, g)
			if err != nil {
				log.Printf("unable to unmarshal while getting game with id %s: %v\n", id, err)
				return status.Error(codes.Internal, "failed to unmarshal while getting game")
			}
			game = g
			return nil
		})
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
			log.Printf("unable to get game's id for client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to get game id")
		}
		gameIdData, err := item.ValueCopy(nil)
		if err != nil {
			log.Printf("unable to get game id data for client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to get game id data")
		}
		gameId := string(gameIdData)
		gameKey := fmt.Sprintf("%s%s", gamePrefix, gameId)
		item, err = txn.Get([]byte(gameKey))
		if err != nil {
			log.Printf("unable to get game for client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to get game")
		}
		gameData, err := item.ValueCopy(nil)
		if err != nil {
			log.Printf("unable to get game data for client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to get game data")
		}
		g := &Game{}
		err = proto.Unmarshal(gameData, g)
		if err != nil {
			log.Printf("unable to marshal game for client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to marshal game")
		}
		game = g
		return nil
	})
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *storageNoSql) updateGame(game *Game) error {
	gameData, err := proto.Marshal(game)
	if err != nil {
		log.Printf("unable to marshal game %s when updating: %v\n", game.Id, err)
		return status.Error(codes.Internal, "failed to update game")
	}
	key := fmt.Sprintf("%s%s", gamePrefix, game.Id)
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), gameData)
	})
}

func (s *storageNoSql) enqueueUpdate(clientId string, updateType isUpdate_Type) (*Update, error) {
	lastSeqNumKey := fmt.Sprintf("%s%s", clientLastSeqNumPrefix, clientId)
	var update *Update
	err := s.db.Update(func(txn *badger.Txn) error {
		var lastSeqNum int64 = 0
		item, err := txn.Get([]byte(lastSeqNumKey))
		if err != nil && err != badger.ErrKeyNotFound {
			log.Printf("unable to get last sequence number of client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to get last sequence number")
		}
		if item != nil {
			bytes, err := item.ValueCopy(nil)
			if err != nil {
				log.Printf("unable to copy last sequence number of client %s as bytes: %v\n", clientId, err)
				return status.Error(codes.Internal, "failed to copy last sequence number")
			}
			lastSeqNum, err = strconv.ParseInt(string(bytes), 10, 64)
			if err != nil {
				log.Printf("unable to parse sequence number bytes of client %s as int64: %v\n", clientId, err)
				return status.Error(codes.Internal, "failed to parse last sequence number")
			}
		}
		nextSeqNum := lastSeqNum + 1
		u := &Update{SequenceNumber: nextSeqNum, Type: updateType}
		updateData, err := proto.Marshal(u)
		if err != nil {
			log.Printf("unable to marshal update to be enqueued for client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to marshal update to be enqueued")
		}
		updatesKey := fmt.Sprintf("%s%s:%d", clientUpdatePrefix, clientId, nextSeqNum)
		err = txn.Set([]byte(updatesKey), updateData)
		if err != nil {
			log.Printf("unable to enqueue update for client %s: %v\n", clientId, err)
			return status.Error(codes.Internal, "failed to enqueue update")
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
			err := item.Value(func(value []byte) error {
				update := &Update{}
				if err := proto.Unmarshal(value, update); err != nil {
					log.Printf("unable to unmarshal update for client %s: %v\n", clientId, err)
					return status.Error(codes.Internal, "failed to unmarshal update")
				}
				list = append(list, update)
				return nil
			})
			if err != nil {
				log.Printf("unable to get update item value for client %s: %v\n", clientId, err)
				return status.Error(codes.Internal, "failed to get update item value")
			}
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
				log.Printf("unable to dequeue %d update(s) for client %s: %v\n", len(updates), clientId, err)
				return status.Error(codes.Internal, "failed to dequeue updates")
			}
		}
		return nil
	})
}
