package main

import (
	"log"

	"github.com/dgraph-io/badger/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type storageNoSql struct {
	db *badger.DB
	storage
}

func newStorageNoSql(db *badger.DB) *storageNoSql {
	return &storageNoSql{db: db}
}

func (s *storageNoSql) saveClientId(clientId string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(clientId), []byte{1})
		return txn.SetEntry(entry)
	})
	if err != nil {
		log.Printf("unable to save client id %s: %v", clientId, err)
		return status.Error(codes.Internal, "client id not saved")
	}
	return nil
}

func (s *storageNoSql) containsClientId(clientId string) (bool, error) {
	var value []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(clientId))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil {
		log.Printf("unable to determine if db contains client id %s: %v", clientId, err)
		return false, status.Error(codes.Internal, "client id existence cannot be determined")
	}
	if len(value) == 0 {
		return false, nil
	}
	return true, nil
}
