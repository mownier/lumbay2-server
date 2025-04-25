package main

type storageNoSql struct {
	storage
}

func newStorageNoSql() *storageNoSql {
	return &storageNoSql{}
}

func (s *storageNoSql) validateClientId(clientId string) error {
	return nil
}
