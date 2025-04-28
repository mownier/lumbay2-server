package main

import (
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	clientIdPrefix   = "client_id:"
	gamePrefix       = "game:"
	gamePlayerPrefix = "game:player:"
)

type storageNoSql struct {
	db *badger.DB
}

func newStorageNoSql(db *badger.DB) *storageNoSql {
	return &storageNoSql{db: db}
}

func (s *storageNoSql) saveClientId(clientId string) error {
	key := clientIdPrefix + clientId
	err := s.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), []byte(clientId))
		return txn.SetEntry(entry)
	})
	if err != nil {
		log.Printf("unable to save client id %s: %v\n", clientId, err)
		return status.Error(codes.Internal, "client id not saved")
	}
	return nil
}

func (s *storageNoSql) containsClientId(clientId string) (bool, error) {
	key := clientIdPrefix + clientId
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("unable to determine if db contains client id %s: %v\n", clientId, err)
		if err == badger.ErrKeyNotFound {
			return false, nil
		}
		return false, status.Error(codes.Internal, "client id existence cannot be determined")
	}
	return true, nil
}

func (s *storageNoSql) getAllClientIds() ([]string, error) {
	list := []string{}
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(clientIdPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				list = append(list, string(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("unable to get all client ids: %v\n", err)
		return nil, status.Error(codes.Internal, "failed to get all client ids")
	}
	return list, nil
}

func (s *storageNoSql) hasGame(gameId string) (bool, error) {
	key := gamePrefix + gameId
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("unable to determine if db contains game %s: %v\n", gameId, err)
		if err == badger.ErrKeyNotFound {
			return false, nil
		}
		return false, status.Error(codes.Internal, "game existence cannot be determined")
	}
	return true, nil
}

func (s *storageNoSql) insertGame(player1 string) (*Game, error) {
	gameId := uuid.New().String()
	gameExists, err := s.hasGame(gameId)
	if err != nil {
		return nil, err
	}
	if gameExists {
		return nil, status.Error(codes.AlreadyExists, "generated id for the game to be inserted already exists")
	}
	game := &Game{Id: gameId, Player1: player1, Player2: ""}
	gameData, err := proto.Marshal(game)
	if err != nil {
		log.Printf("unable to marshal to be inserted game to data: %v\n", err)
		return nil, status.Error(codes.Internal, "failed to marshal the to-be-inserted game")
	}
	gameKey := gamePrefix + game.Id
	gamePlayerKey := gamePlayerPrefix + game.Player1
	err = s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(gameKey), gameData)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(gamePlayerKey), []byte(game.Id))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("unable to insert game for player1 %s: %v\n", player1, err)
		return nil, status.Error(codes.Internal, "failed to insert game")
	}
	return game, nil
}
