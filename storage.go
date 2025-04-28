package main

type storage interface {
	saveClientId(clientId string) error
	containsClientId(clientId string) (bool, error)
	getAllClientIds() ([]string, error)
}
