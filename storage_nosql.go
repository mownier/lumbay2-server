package main

import "github.com/dgraph-io/badger/v4"

type storageNoSql struct {
	db *badger.DB
	storage
}

func newStorageNoSql(db *badger.DB) *storageNoSql {
	return &storageNoSql{}
}

func (s *storageNoSql) validateClientId(clientId string) error {
	return nil
}
